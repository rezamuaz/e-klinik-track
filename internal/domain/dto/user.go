package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserResponse struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Role      []string           `bson:"role"`
	Place     primitive.ObjectID `bson:"place"`
	PlaceName string             `bson:"place_name"`
	Room      primitive.ObjectID `bson:"room"`
	RoomName  string             `bson:"room_name"`
}
