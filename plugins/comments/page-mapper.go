package comments

import (
	"strings"
)

type CommentPageToUrlMapper interface {
	Map(page string, id string) string
}

type formatPageMapper struct {
	format string
}

func (f *formatPageMapper) Map(page string, id string) string {
	url := strings.ReplaceAll(f.format, "{page}", page)
	url = strings.ReplaceAll(url, "{id}", id)
	return url
}
