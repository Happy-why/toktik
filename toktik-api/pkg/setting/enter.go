package setting

type group struct {
	config config
}

var Group = new(group)

func InitAll() {
	Group.config.Init()
}
