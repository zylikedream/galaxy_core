package component

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoleInfo struct {
	RoleID  primitive.ObjectID `bson:"_id"`
	Account string             `bson:"account"`
	Name    string             `bson:"name"`
}

func (r *RoleInfo) GetName() string {
	return "role_info"
}
