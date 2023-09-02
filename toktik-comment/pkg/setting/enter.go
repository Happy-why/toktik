package setting

var Settings []InitSettings

type InitSettings interface {
	InitSetting()
}

func InitAllSetting() {
	InitBootStrap()
	for _, setting := range Settings {
		setting.InitSetting()
	}
}
