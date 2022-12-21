package publicsite

import "tiim/go-comment-api/config"

type publicSitePlugin struct {
	DebugApertureId string `json:"debug_aperture_id"`
}

func init() {
	config.RegisterModule(&publicSitePlugin{})
}

func (p *publicSitePlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "public-site",
		New:  func() config.Module { return new(publicSitePlugin) },
		Docs: config.ConfigDocs{
			DocString: `Public site module. This module enables the public info site.`,
			Fields: map[string]string{
				"DebugApertureId": "Used for debugging only. Sets the rel=microsub link to the url of the aperture.p3k.io endpoint with the given id.",
			},
		},
	}
}

func (p *publicSitePlugin) Load(c config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	var plugin config.ApiPluginInstance
	plugin = newPublicModule(p.DebugApertureId)

	return plugin, nil
}
