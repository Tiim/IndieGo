package admin

import (
	"tiim/go-comment-api/config"
)

type adminModule struct {
	Password string `json:"password"`
}

func init() {
	config.RegisterModule(&adminModule{})
}

func (p *adminModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "admin",
		New:  func() config.Module { return new(adminModule) },
	}
}

func (p *adminModule) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	return newAdminModule(p.Password), nil
}
