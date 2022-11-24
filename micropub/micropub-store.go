package micropub

import (
	"willnorris.com/go/microformats"
)

type MicropubStore interface {
	Create(post MicropubPost) (string, error)
	Modify(url string, deleteProps interface{}, addProps, replaceProps map[string][]interface{}) error
	Delete(url string) error
	UnDelete(url string) error
	Get(url string) (*microformats.Microformat, error)
}
