package component

type RoleInfo struct {
	AccountName string
	RoleID      uint64
}

func Load() error {
	return nil
}

func Name() string {
	return "role_info"
}

func Unload() error {
	return nil
}
