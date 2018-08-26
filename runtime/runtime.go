package runtime

import (
	"net/url"
)

type URL string

func (raw *URL) Get() *url.URL {
	if raw == nil {
		return nil
	}
	if *raw == "" {
		return nil
	}
	u, err := url.Parse(string(*raw))
	if err != nil {
		return nil
	}
	return u
}

func (raw *URL) Set(u *url.URL) {
	if u == nil {
		return
	}
	*raw = URL(u.String())
}

type File struct {
	Name     string `firestore:"name"`
	URL      URL    `firestore:"url"`
	MIMEType string `firestore:"mimeType"`
}
