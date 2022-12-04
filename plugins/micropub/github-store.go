package micropub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"willnorris.com/go/microformats"
)

type UrlConverter interface {
	UrlToFilePath(url string) string
	FilePathToUrl(path string) string
}

type micropubGithubStore struct {
	token        string
	user         string
	repo         string
	folder       string
	urlConverter UrlConverter
	client       *http.Client
	rand         *rand.Rand
}

func newMicropubGithubStore(token, user, repo, folder string, urlConverter UrlConverter, client *http.Client) *micropubGithubStore {
	return &micropubGithubStore{
		token:        token,
		user:         user,
		repo:         repo,
		folder:       folder,
		urlConverter: urlConverter,
		client:       client,
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *micropubGithubStore) Create(post MicropubPost) (string, error) {
	filePath := m.nextFile()
	url, err := url.Parse("https://api.github.com/repos/" + m.user + "/" + m.repo + "/contents/" + filePath + ".md")
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(map[string]interface{}{
		"message": "create post " + filePath,
		"content": base64.StdEncoding.EncodeToString([]byte(post.ToMarkdown())),
	})
	if err != nil {
		return "", err
	}

	req := http.Request{
		Method: http.MethodPut,
		URL:    url,
		Header: http.Header{
			"Authorization": {"Bearer " + m.token},
		},
		Body: io.NopCloser(bytes.NewBuffer(buf)),
	}
	res, err := m.client.Do(&req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
	}
	return m.urlConverter.FilePathToUrl(filePath), err
}

func (m *micropubGithubStore) Modify(u string, deleteProps interface{}, addProps, replaceProps map[string][]interface{}) error {
	filePath := m.urlConverter.UrlToFilePath(u)
	url, err := url.Parse("https://api.github.com/repos/" + m.user + "/" + m.repo + "/contents/" + filePath + ".md")
	if err != nil {
		return err
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: http.Header{
			"Authorization": {"Bearer " + m.token},
			"Accept":        {"application/vnd.github+json"},
		},
	}
	res, err := m.client.Do(&req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var responseData map[string]interface{}
	json.Unmarshal(result, &responseData)
	content, err := base64.StdEncoding.DecodeString(responseData["content"].(string))
	sha := responseData["sha"].(string)
	if err != nil {
		return err
	}

	post := PostFromMarkdown(string(content))
	ModifyEntry(&post, deleteProps, addProps, replaceProps)

	buf, err := json.Marshal(map[string]interface{}{
		"message": "update post " + filePath,
		"content": base64.StdEncoding.EncodeToString([]byte(post.ToMarkdown())),
		"sha":     sha,
	})
	if err != nil {
		return err
	}
	req = http.Request{
		Method: http.MethodPut,
		URL:    url,
		Header: http.Header{
			"Authorization": {"Bearer " + m.token},
		},
		Body: io.NopCloser(bytes.NewBuffer(buf)),
	}
	res, err = m.client.Do(&req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		result, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("error updating post: %s", result)
	}
	return nil
}

func (m *micropubGithubStore) Delete(u string) error {
	filePath := m.urlConverter.UrlToFilePath(u)
	url, err := url.Parse("https://api.github.com/repos/" + m.user + "/" + m.repo + "/contents/" + filePath + ".md")
	if err != nil {
		return err
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: http.Header{
			"Authorization": {"Bearer " + m.token},
			"Accept":        {"application/vnd.github+json"},
		},
	}
	res, err := m.client.Do(&req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var responseData map[string]interface{}
	json.Unmarshal(result, &responseData)
	sha := responseData["sha"].(string)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(map[string]interface{}{
		"message": "delete post " + filePath,
		"sha":     sha,
	})
	if err != nil {
		return err
	}

	req = http.Request{
		Method: http.MethodDelete,
		URL:    url,
		Header: http.Header{
			"Authorization": {"Bearer " + m.token},
		},
		Body: io.NopCloser(bytes.NewBuffer(buf)),
	}
	res, err = m.client.Do(&req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		result, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("error deleting post: %s", result)
	}
	return nil
}

func (m *micropubGithubStore) UnDelete(url string) error {
	return fmt.Errorf("undelete not implemented")
}

func (m *micropubGithubStore) Get(u string) (*microformats.Microformat, error) {
	filePath := m.urlConverter.UrlToFilePath(u)
	url, err := url.Parse("https://api.github.com/repos/" + m.user + "/" + m.repo + "/contents/" + filePath + ".md")
	if err != nil {
		return nil, err
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: http.Header{
			"Authorization": {"Bearer " + m.token},
			"Accept":        {"application/vnd.github+json"},
		},
	}
	res, err := m.client.Do(&req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseData map[string]interface{}
	json.Unmarshal(result, &responseData)
	if responseData["content"] == nil {
		return nil, fmt.Errorf("not found")
	}
	content, err := base64.StdEncoding.DecodeString(responseData["content"].(string))
	if err != nil {
		return nil, err
	}

	post := PostFromMarkdown(string(content))
	return post.Entry.ToMicroformat(), nil
}

func (m *micropubGithubStore) nextFile() string {
	// create randowm 6 char string
	name := strings.ToLower(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", m.rand.Int()))))[:6]
	now := time.Now().Format("2006/01")
	return fmt.Sprintf("%s/%s/%s", m.folder, now, name)
}
