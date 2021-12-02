package component

type RoleInfo struct {
	Account string `bson:"account"`
	RoleID  uint64 `bson:"role_id"`
	Name    string `bson:"name"`
}

func (r *RoleInfo) GetName() string {
	return "role_info"
}
