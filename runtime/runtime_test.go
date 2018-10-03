package runtime

import (
	"net/url"
	"reflect"
	"testing"
)

func TestURL_Get(t *testing.T) {
	tests := []struct {
		name string
		raw  URL
		want *url.URL
	}{
		{
			name: "blank",
			raw:  URL(""),
			want: nil,
		},
		{
			name: "basic",
			raw:  URL("https://example.com/"),
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/",
			},
		},
		{
			name: "other scheme",
			raw:  URL("gs://example.com/path/to/file"),
			want: &url.URL{
				Scheme: "gs",
				Host:   "example.com",
				Path:   "/path/to/file",
			},
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			got := tt.raw.Get()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URL.Get() = %v, want %v", got, tt.want)
			}
			var other URL
			other.Set(got)
			if !reflect.DeepEqual(other, tt.raw) {
				t.Errorf("URL failed to round trip => %v, want %v", other, tt.raw)
			}
		})
	}
}

func TestURL_Set(t *testing.T) {
	tests := []struct {
		name   string
		raw    *URL
		args   *url.URL
		exp    *URL
		panics bool
	}{
		{
			name: "from empty",
			raw:  new(URL),
			args: &url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			exp: func() *URL { u := URL("https://example.com"); return &u }(),
		},
		{
			name:   "from nil panics",
			raw:    nil,
			args:   &url.URL{},
			panics: true,
		},
		{
			name: "from set",
			raw:  func() *URL { u := URL("https://example.com"); return &u }(),
			args: &url.URL{
				Scheme: "http",
				Host:   "yahoo.com",
			},
			exp: func() *URL { u := URL("http://yahoo.com"); return &u }(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("did not panic")
					}
				}()

				tt.raw.Set(tt.args)
			} else {
				tt.raw.Set(tt.args)
				got := tt.raw
				if !reflect.DeepEqual(*got, *tt.exp) {
					t.Errorf("URL.Set() = %v, want %v", *got, *tt.exp)
				}
			}
		})
	}
}
