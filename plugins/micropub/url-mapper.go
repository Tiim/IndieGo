package micropub

import "strings"

type suffixPrefixUrlMapper struct {
	urlSuffix string
	urlPrefix string
	folder    string
	extension string
}

func (m *suffixPrefixUrlMapper) UrlToFilePath(url string) string {
	filename := strings.TrimPrefix(url, m.urlPrefix)
	filename = strings.TrimSuffix(filename, m.urlSuffix)
	return m.folder + filename + m.extension
}

func (m *suffixPrefixUrlMapper) FilePathToUrl(path string) string {
	filename := strings.TrimPrefix(path, m.folder)
	filename = strings.TrimSuffix(filename, m.extension)
	return m.urlPrefix + filename + m.urlSuffix
}
