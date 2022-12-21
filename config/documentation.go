package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type docsHierarchy struct {
	Namespace string
	Docs      string
	Children  []*docsHierarchy
}

func (d *docsHierarchy) AddModule(module Module) {
	nsString, docs := generateModuleDocs(module)
	namespace := strings.Split(nsString, ".")
	d.addDocs(namespace, docs)
}

func (d *docsHierarchy) addDocs(namespace []string, docs string) {
	if len(namespace) == 1 {
		d.Children = append(d.Children, &docsHierarchy{
			Namespace: namespace[0],
			Docs:      docs,
		})
		sort.Slice(d.Children, d.Less)
		return
	} else {
		for _, child := range d.Children {
			if child.Namespace == namespace[0] {
				child.addDocs(namespace[1:], docs)
				return
			}
		}
		newChild := &docsHierarchy{
			Namespace: namespace[0],
		}
		newChild.addDocs(namespace[1:], docs)
		d.Children = append(d.Children, newChild)
		sort.Slice(d.Children, d.Less)
	}
}

func (d *docsHierarchy) Less(i, j int) bool {
	return d.Children[i].Namespace < d.Children[j].Namespace
}

func (d *docsHierarchy) String(topNamespaces []string) string {

	var namespaceSegments []string
	if d.Namespace != "" {
		namespaceSegments = append(topNamespaces, d.Namespace)
	} else {
		namespaceSegments = topNamespaces
	}
	heading := "h3"
	namespaceString := strings.Join(namespaceSegments, ".")
	docs := ""
	if len(namespaceSegments) > 0 {
		docs += fmt.Sprintf(`<%s id="%s">Namespace: %s </%s>`+"\n", heading, namespaceString, namespaceString, heading)
	} else {
		docs += "<h2>Namespaces</h2>\n"
	}
	docs += d.Docs
	for _, child := range d.Children {
		docs += child.String(namespaceSegments)
	}
	return docs
}

func GenerateDocs() string {
	docs := &docsHierarchy{}
	for _, module := range modules {
		docs.AddModule(module)
	}
	return docs.String([]string{})
}

func generateModuleDocs(module Module) (string, string) {

	info := module.IndieGoModule()
	namespace := info.Name

	docs := "<p>"
	docs += info.Docs.DocString
	docs += "</p>\n"

	exampleJson := "<pre><code>"
	exampleJson += "{\n\t\"name\": \"" + namespace + "\",\n\t\"args\": {\n"

	typ := reflect.TypeOf(module)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fieldDesc := "<ul>\n"
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Tag.Get("json")
		typ := field.Type
		prefix := ""
		prefixJson := ""
		suffixJson := ""
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
			prefix = "[]"
			prefixJson = "["
			suffixJson = "]"
			typ = typ.Elem()
		}

		fieldDesc += "<li>\n"
		if typ == reflect.TypeOf(ModuleRaw{}) {
			fc := field.Tag.Get("config")
			fieldDesc += fmt.Sprintf("Field <code>%s</code>, type: <a href=\"#%s\"><code>%s%s</code></a>\n", name, fc, prefix, fc)
			exampleJson += fmt.Sprintf("\t\t\"%s\": %s{ \"name\": \"%s\", \"args\": { <a href=\"#%s\">...</a> } }%s,\n", name, prefixJson, fc, fc, suffixJson)
		} else {
			fieldDesc += fmt.Sprintf("Field <code>%s</code>, type: <code>%s%s</code>\n", name, prefix, typ.Name())
			jsond, err := json.Marshal(reflect.New(typ).Interface())
			if err != nil {
				panic(err)
			}
			exampleJson += fmt.Sprintf("\t\t\"%s\": %s%s%s,\n", name, prefixJson, jsond, suffixJson)
		}
		fieldDesc += fmt.Sprintf("<p>%s</p>", info.Docs.Fields[field.Name])
		fieldDesc += "</li>\n"
	}
	fieldDesc += "</ul>\n"
	exampleJson += "\t}\n}\n"
	exampleJson += "</code></pre>"
	return namespace, docs + fieldDesc + exampleJson
}
