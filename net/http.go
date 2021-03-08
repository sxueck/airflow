package net

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type HTTPOptions struct {
	URL         string
	body        *http.Response
	Cookie      []*http.Cookie
	ContentType string
	UserAgent   string
	Err         error
	ProxyURL    string
}

func New() *HTTPOptions {
	return &HTTPOptions{
		ContentType: "application/json",
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36",
		Err:         nil,
	}
}

func (h *HTTPOptions) ObtainCookie() {
	if h.body == nil {
		log.Fatal("please call the get method first")
	}
	h.Cookie = h.body.Cookies()
}

func (h *HTTPOptions) exBody() string {
	if h.Err != nil {
		return ""
	}

	var resContent []byte
	resContent, h.Err = ioutil.ReadAll(h.body.Body)

	return string(resContent)
}

func (h *HTTPOptions) GET() string {
	req, _ := http.NewRequest(http.MethodGet, h.URL, nil)
	if h.Cookie != nil {
		for _, c := range h.Cookie {
			req.AddCookie(c)
		}
	}
	return h.do(req)
}

func (h *HTTPOptions) POST(body string) string {
	log.Printf("POST %s to %s", body, h.URL)
	req, _ := http.NewRequest(http.MethodPost, h.URL, bytes.NewBufferString(body))
	if h.Cookie != nil {
		for _, c := range h.Cookie {
			req.AddCookie(c)
		}
	}
	return h.do(req)
}

func (h *HTTPOptions) do(r *http.Request) string {
	r.Header.Set("Content-Type", h.ContentType)
	r.Header.Set("User-Agent", h.UserAgent)

	client := &http.Client{}

	if h.ProxyURL != "" {
		proxy := func(r *http.Request) (*url.URL, error) {
			return url.Parse(r.URL.String())
		}
		client.Transport = &http.Transport{Proxy: proxy}
	}

	h.body, h.Err = client.Do(r)
	if h.Err == nil && h.body.StatusCode == 200 {
		return h.exBody()
	}
	return ""
}
