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
		Docs: config.ConfigDocs{
			DocString: `Github webhook trigger. This trigger is used to receive webhooks from Github.`,
			Fields: map[string]string{
				"Name":          "The name of the trigger. This is used to for the endpoint url. Resulting endpoint: webhooks/{name}",
				"WebhookSecret": "The secret used to verify the webhook, as specified on github. Webhook that don't have a valid signature will be ignored.",
			},
		},
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
