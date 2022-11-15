package microformatsextract

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"willnorris.com/go/microformats"
)

func TestGetHApp(t *testing.T) {
	baseUrl, _ := url.Parse("https://webmention.rocks")
	type args struct {
		data *microformats.Data
	}
	tests := []struct {
		name string
		args args
		want MF2HApp
	}{
		{
			name: "webmention.rocks",
			args: args{
				data: microformats.Parse(strings.NewReader(`<div style="display: none;" class="h-app"><a href="/" class="u-url p-name">Webmention.rocks!</a><img src="/assets/webmention-rocks-icon.png" class="u-logo"><a href="https://indielogin.com/redirect/indieauth" class="u-redirect-uri"></a></div>`), baseUrl),
			},
			want: MF2HApp{
				Url:          "https://webmention.rocks/",
				Name:         "Webmention.rocks!",
				Logo:         "https://webmention.rocks/assets/webmention-rocks-icon.png",
				RedirectUris: []string{"https://indielogin.com/redirect/indieauth"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHApp(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEntryWithType(t *testing.T) {
	type args struct {
		data  *microformats.Data
		types []string
	}
	tests := []struct {
		name string
		args args
		want *microformats.Microformat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEntryWithType(tt.args.data, tt.args.types...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEntryWithType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringProp(t *testing.T) {
	type args struct {
		name string
		item *microformats.Microformat
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringProp(tt.args.name, tt.args.item); got != tt.want {
				t.Errorf("GetStringProp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringPropSlice(t *testing.T) {
	type args struct {
		name string
		item *microformats.Microformat
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringPropSlice(tt.args.name, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStringPropSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHCard(t *testing.T) {
	type args struct {
		name string
		item *microformats.Microformat
	}
	tests := []struct {
		name string
		args args
		want MF2HCard
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHCard(tt.args.name, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
