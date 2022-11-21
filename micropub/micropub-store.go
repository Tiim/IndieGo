package micropub

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/mfobjects"

	"willnorris.com/go/microformats"
)

type MicropubStore interface {
	Create(post MicropubPost) (string, error)
	Modify(url string, deleteProps interface{}, addProps, replaceProps map[string][]interface{}) error
	Delete(url string) error
	UnDelete(url string) error
	Get(url string) (*microformats.Microformat, error)
}

type micropubMemoryStore struct {
	index int
	store map[string]MicropubPost
}

func NewMicropubPrintStore() *micropubMemoryStore {
	return &micropubMemoryStore{
		index: 1,
		store: make(map[string]MicropubPost),
	}
}

func (m *micropubMemoryStore) Create(post MicropubPost) (string, error) {
	url := m.nextUrl()
	m.store[url] = post
	m.printStore()
	return url, nil
}

func (m *micropubMemoryStore) Modify(url string, deleteProps interface{}, addProps, replaceProps map[string][]interface{}) error {
	post, ok := m.store[url]
	if !ok {
		return fmt.Errorf("post not found")
	}
	mf := post.Entry.ToMicroformat()
	for key, values := range replaceProps {
		mf.Properties[key] = values
	}
	for key, values := range addProps {
		_, ok := mf.Properties[key]
		if !ok {
			mf.Properties[key] = values
		} else {
			mf.Properties[key] = append(mf.Properties[key], values...)
		}
	}
	if deleteProps != nil {
		switch del := deleteProps.(type) {
		case []interface{}:
			for _, key := range del {
				delete(mf.Properties, key.(string))
			}
		case map[string]interface{}:
			for key, values := range del {
				for _, value := range values.([]interface{}) {
					for i, v := range mf.Properties[key] {
						if v == value {
							mf.Properties[key] = append(mf.Properties[key][:i], mf.Properties[key][i+1:]...)
						}
					}
				}
			}
		default:
			return fmt.Errorf("unknown delete type %T", del)
		}
	}
	m.store[url] = MicropubPost{Entry: mfobjects.GetHEntry(&microformats.Data{Items: []*microformats.Microformat{mf}})}
	m.printStore()
	return nil
}

func (m *micropubMemoryStore) Delete(url string) error {
	delete(m.store, url)
	m.printStore()
	return nil
}

func (m *micropubMemoryStore) UnDelete(url string) error {
	return nil
}

func (m *micropubMemoryStore) Get(url string) (*microformats.Microformat, error) {
	post, ok := m.store[url]
	if !ok {
		return nil, fmt.Errorf("post not found")
	}
	return post.Entry.ToMicroformat(), nil
}

func (m *micropubMemoryStore) nextUrl() string {
	m.index++
	return fmt.Sprintf("https://tiim.ch/mf2/test/%d", m.index)
}

func (m *micropubMemoryStore) printStore() {
	prettyJson, err := json.MarshalIndent(m.store, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(prettyJson))
}
