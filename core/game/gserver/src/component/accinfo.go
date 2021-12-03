package component

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoleAccount struct {
	RoleID  primitive.ObjectID `bson:"_id,omitempty"`
	Account string             `bson:"account"`
}
