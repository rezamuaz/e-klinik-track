package pkg

import (
	"context"
	"e-klinik/config"
	"e-klinik/pkg/logging"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewS3Storage(cfg *config.Config, log logging.Logger) (*minio.Client, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	minioClient, err := minio.New(cfg.Minio.EndPoint, &minio.Options{
		Region: cfg.Minio.Region,
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.SSL,
	})
	if err != nil {
		log.Fatal(logging.S3, logging.Startup, err.Error(), nil)
		return nil, err
	}

	// err = CreateOrSkipBucket(minioClient, ctx, cfg.Minio.Bucket1, cfg.Minio.Region, log)
	// if err != nil {
	// 	return nil, err
	// }
	// err = CreateOrSkipBucket(minioClient, ctx, cfg.Minio.Bucket2, cfg.Minio.Region, log)
	// if err != nil {
	// 	return nil, err
	// }
	return minioClient, nil
}

func CreateOrSkipBucket(client *minio.Client, ctx context.Context, bucketName string, region string, log logging.Logger) error {

	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Info(logging.S3, logging.IsExist, fmt.Sprintf("We already own %s\n", bucketName), nil)
		} else {
			log.Fatalf(err.Error())
			return err
		}
	} else {
		log.Info(logging.S3, logging.IsExist, fmt.Sprintf("Successfully created %s\n", bucketName), nil)
	}
	return nil
}
