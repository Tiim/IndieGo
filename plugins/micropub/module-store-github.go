package micropub

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type githubStoreModule struct{}
type githubStoreModuleData struct {
	GithubToken  string `json:"github_token"`
	GithubUser   string `json:"github_user"`
	GithubRepo   string `json:"github_repo"`
	GithubFolder string `json:"github_folder"`
	UrlPrefix    string `json:"url_prefix"`
	UrlSuffix    string `json:"url_suffix"`
}

func init() {
	config.RegisterModule(&githubStoreModule{})
}

func (m *githubStoreModule) Name() string {
	return "micropub-store-github"
}

func (m *githubStoreModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	var d githubStoreModuleData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	mapper := &suffixPrefixUrlMapper{
		urlPrefix: d.UrlPrefix,
		urlSuffix: d.UrlSuffix,
		folder:    d.GithubFolder,
		extension: ".md",
	}

	return newMicropubGithubStore(
		d.GithubToken,
		d.GithubUser,
		d.GithubRepo,
		d.GithubFolder,
		mapper,
		config.HttpClient,
	), nil
}
