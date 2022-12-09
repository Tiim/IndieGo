package micropub

import (
	"tiim/go-comment-api/config"
)

type githubStoreModule struct {
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

func (m *githubStoreModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "micropub.store.github",
		New:  func() config.Module { return new(githubStoreModule) },
	}
}

func (m *githubStoreModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	mapper := &suffixPrefixUrlMapper{
		urlPrefix: m.UrlPrefix,
		urlSuffix: m.UrlSuffix,
		folder:    m.GithubFolder,
		extension: ".md",
	}

	return newMicropubGithubStore(
		m.GithubToken,
		m.GithubUser,
		m.GithubRepo,
		m.GithubFolder,
		mapper,
		config.HttpClient,
	), nil
}
