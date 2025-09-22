package worker

import (
	"context"
	"e-klinik/infra/types"
	"e-klinik/internal/domain/dto"
	"e-klinik/pkg"
	"e-klinik/pkg/constant"
	"e-klinik/pkg/logging"
	"e-klinik/utils"
	"fmt"

	"github.com/streadway/amqp"
)

type ConsumerService struct {
	Logger logging.Logger
	RMQ    *pkg.RabbitMQ
	TsRepo types.IndexRepository
	// PostRepo *pgr.PostRepositoryImpl
	Ch   *amqp.Channel
	Done chan struct{}
}

func (s *ConsumerService) StartRabbitConsumer(ctx context.Context) error {
	fmt.Println("starting consumer")
	msgs, err := s.Ch.Consume(
		constant.QueueName,       // queue
		constant.RMQConsumerName, // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			s.Logger.Info(logging.Rabbit, logging.Received, "No more messages to consume. Exiting.", nil)
			s.Done <- struct{}{}
		}()

		for {
			select {
			case <-ctx.Done():
				s.Logger.Info(logging.Rabbit, logging.Received, "Shutting down RMQ consumer...", nil)
				return

			case msg, ok := <-msgs:
				if !ok {
					fmt.Println("Message channel closed! Exiting.")
					s.Logger.Info(logging.Rabbit, logging.Received, "Message channel closed! Exiting.", nil)
					return
				}
				fmt.Println(msg.RoutingKey)
				s.Logger.Info(logging.Rabbit, logging.Received, fmt.Sprintf("Received message: %s", msg.RoutingKey), nil)

				var nack bool
				switch msg.RoutingKey {
				//Created Post
				case fmt.Sprintf("%s%s", constant.RoutingKey, constant.PostCreated):
					post, err := utils.ByteToAny[dto.PostCreatedResponse](msg.Body)
					if err != nil {
						nack = true
						continue
					}

					if err := s.TsRepo.CreateIndex(ctx, "pvsave", post); err != nil {
						nack = true
					}

				// case fmt.Sprintf("%s_%s", constant.RoutingKey, constant.PostUpdated):
				// 	postId, err := utils.ByteToAny[params.UpdatePost](msg.Body)
				// 	if err != nil {
				// 		nack = true
				// 		continue
				// 	}
				// data, err := s.PostRepo.GetPostForIndexByID(ctx, postId.ID)
				// if err != nil {
				// 	nack = true
				// 	continue
				// }
				// cleanedPost, err := dto.ToPostIndexCreatedResp(data)
				// if err := s.TsRepo.Upsert(ctx, "pvsave", cleanedPost); err != nil {
				// 	nack = true
				// }
				case fmt.Sprintf("%s_%s", constant.RoutingKey, constant.PostDeleted):
					post, err := utils.ByteToAny[dto.PostedDelete](msg.Body)
					if err != nil {
						nack = true
						break
					}
					if err := s.TsRepo.DeleteIndex(ctx, "pvsave", post.PostID); err != nil {
						nack = true
					}

				default:
					nack = true
				}

				if nack {
					s.Logger.Info(logging.Rabbit, logging.Received, "NAcking :(", nil)
					_ = msg.Nack(false, false)
				} else {
					s.Logger.Info(logging.Rabbit, logging.Received, "Acking :)", nil)
					_ = msg.Ack(false)
				}
			}
		}
	}()

	return nil
}
