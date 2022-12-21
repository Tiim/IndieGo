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
		Docs: config.ConfigDocs{
			DocString: `Admin module. This module enables the admin dashboard.`,
			Fields: map[string]string{
				"Password": "Password for the admin dashboard.",
			},
		},
	}
}

func (p *adminModule) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	return newAdminModule(p.Password), nil
}
