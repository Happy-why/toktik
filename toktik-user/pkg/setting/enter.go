package setting

type group struct {
	config config
	dao    Dao
}

var Group = new(group)

func InitAll() {
	Group.config.Init()
	Group.dao.Init()
}
