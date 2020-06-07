package render

import "github.com/appconf/engine/util"

//Config 应用模板配置
type Config struct {
	Src       string   `toml:"src"`
	Dst       string   `toml:"dst"`
	Keys      []string `toml:"keys"`
	ReloadCmd string   `toml:"reload_cmd"`
}

//GetKeys 获取配置文件中模板里指定的所有key
//这些key用来去Storage中获取数据
func GetKeys(cfgs []*Config) []string {
	keys := make([]string, 0)
	for _, cfg := range cfgs {
		keys = append(keys, cfg.Keys...)
	}
	return util.RemoveRepeatElement(keys)
}
