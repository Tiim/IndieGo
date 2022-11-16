package microformatsextract

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"willnorris.com/go/microformats"
)

func getTestFiles(mfType string) []string {
	files, _ := ioutil.ReadDir("testdata")
	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if strings.HasPrefix(file.Name(), mfType) && strings.HasSuffix(file.Name(), ".html") {
			fileNames = append(fileNames, strings.TrimSuffix(file.Name(), ".html"))
		}
	}
	return fileNames
}

func getTestFile(name string) (io.Reader, io.Reader, error) {
	got, err := os.Open("testdata/" + name + ".html")
	if err != nil {
		return nil, nil, err
	}
	want, err := os.Open("testdata/" + name + ".json")
	if err != nil {
		return nil, nil, err
	}
	return got, want, nil
}

func deepEqual(a, b interface{}, t *testing.T, debugMf *microformats.Data, fileName string) bool {
	if !reflect.DeepEqual(a, b) {
		// convert got and want to JSON for easier comparison
		aJSON, _ := json.MarshalIndent(a, "", "  ")
		bJSON, _ := json.MarshalIndent(b, "", "  ")
		debugJSON, _ := json.MarshalIndent(debugMf, "", "  ")
		t.Errorf("File: %s\n got %v\nwant %v\ngot:%s\nwant:%s\ndebug: %s", fileName, a, b, aJSON, bJSON, debugJSON)
		return false
	}
	return true
}

func TestGetHApp(t *testing.T) {
	baseUrl, _ := url.Parse("https://webmention.rocks")

	testFiles := getTestFiles("h-app")

	for _, testFile := range testFiles {
		gotR, wantR, err := getTestFile(testFile)
		if err != nil {
			t.Error(err)
			return
		}
		mf := microformats.Parse(gotR, baseUrl)
		got := GetHApp(mf)
		want := &MF2HApp{}
		json.NewDecoder(wantR).Decode(&want)
		if !deepEqual(got, want, t, mf, testFile) {
			return
		}
	}
}

func TestGetHEntry(t *testing.T) {
	baseUrl, _ := url.Parse("https://webmention.rocks")

	testFiles := getTestFiles("h-entry")

	for _, testFile := range testFiles {
		gotR, wantR, err := getTestFile(testFile)
		if err != nil {
			t.Error(err)
			return
		}
		mf := microformats.Parse(gotR, baseUrl)
		got := GetHEntry(mf)
		want := &MF2HEntry{}
		json.NewDecoder(wantR).Decode(want)
		if !deepEqual(got, want, t, mf, testFile) {
			return
		}
	}
}
