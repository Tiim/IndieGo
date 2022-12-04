package comments

import "fmt"

type CommentPageToUrlMapper interface {
	Map(page string, id string) string
}

type formatPageMapper struct {
	format string
}

func (f *formatPageMapper) Map(page string, id string) string {
	return fmt.Sprintf(f.format, page, id)
}
