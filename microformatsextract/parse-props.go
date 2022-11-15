package microformatsextract

import (
	"willnorris.com/go/microformats"
)

type MF2HCard struct {
	Name string
}

type MF2HApp struct {
	Url          string
	Name         string
	Logo         string
	Summary      string
	Author       MF2HCard
	RedirectUris []string
}

func GetHApp(data *microformats.Data) MF2HApp {
	item := getEntryWithType(data, "h-app", "h-x-app")
	if item == nil {
		return MF2HApp{}
	}
	return MF2HApp{
		Url:          GetStringProp("url", item),
		Name:         GetStringProp("name", item),
		Logo:         GetStringProp("logo", item),
		Summary:      GetStringProp("summary", item),
		Author:       GetHCard("author", item),
		RedirectUris: GetStringPropSlice("redirect-uri", item),
	}
}

func getEntryWithType(data *microformats.Data, types ...string) *microformats.Microformat {
	for _, item := range data.Items {
		for _, itemType := range item.Type {
			for _, searchType := range types {
				if itemType == searchType {
					return item
				}
			}
		}
	}
	return nil
}

func GetStringProp(name string, item *microformats.Microformat) string {
	propValue, ok := item.Properties[name]
	if !ok {
		return ""
	}
	if len(propValue) == 0 {
		return ""
	}
	value, _ := propValue[0].(string)
	return value
}

func GetStringPropSlice(name string, item *microformats.Microformat) []string {
	propValue, ok := item.Properties[name]
	if !ok {
		return []string{}
	}
	if len(propValue) == 0 {
		return []string{}
	}
	slice := make([]string, 0, len(propValue))
	for _, value := range propValue {
		if value, ok := value.(string); ok {
			slice = append(slice, value)
		}
	}
	return slice
}

func GetHCard(name string, item *microformats.Microformat) MF2HCard {
	author := item.Properties["name"]
	if len(author) > 0 {
		authorMf, ok := author[0].(*microformats.Microformat)
		if ok {
			return MF2HCard{Name: authorMf.Value}
		}
	}
	return MF2HCard{}
}
