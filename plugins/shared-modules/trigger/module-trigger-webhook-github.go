package trigger

import "tiim/go-comment-api/config"

type GithubWebhookPlugin struct {
	Name          string `json:"name"`
	WebhookSecret string `json:"webhook_secret"`
}

func init() {
	config.RegisterModule(&GithubWebhookPlugin{})
}

func (p *GithubWebhookPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "trigger.webhook.github",
		New:  func() config.Module { return new(GithubWebhookPlugin) },
	}
}

func (p *GithubWebhookPlugin) Load(c config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	validator := newGithubValidator(p.WebhookSecret)
	webhook := newWebhookModule(p.Name)
	webhook.SetValidator(validator)

	var trigger Trigger = webhook
	var _ config.ApiPluginInstance = webhook
	return trigger, nil
}
