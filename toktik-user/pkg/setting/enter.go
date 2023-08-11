package setting

var Settings []InitSettings

type InitSettings interface {
	InitSetting()
}

func InitAllSetting() {
	for _, setting := range Settings {
		setting.InitSetting()
	}
}
