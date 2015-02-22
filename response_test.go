package mixpanelproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var serveDummy_tests = []struct {
	Desc       string
	Method     string
	Status     int
	Params     url.Values
	PostParams url.Values
	Headers    http.Header
	Body       string
}{{
	Desc:       "GET quiet",
	Method:     "GET",
	Status:     http.StatusOK,
	Params:     url.Values{},
	PostParams: url.Values{},
	Body:       `1`,
	Headers:    http.Header{},
}, {
	Desc:   "GET verbose",
	Method: "GET",
	Status: http.StatusOK,
	Params: url.Values{
		"verbose": {"1"},
	},
	PostParams: url.Values{},
	Body:       `{"error": "", "status": 1}`,
	Headers:    http.Header{},
}, {
	Desc:   "GET redirect",
	Method: "GET",
	Status: http.StatusFound,
	Params: url.Values{
		"redirect": {"http://test.example/foo"},
	},
	PostParams: url.Values{},
	Body:       ``,
	Headers: http.Header{
		"Location":      {"http://test.example/foo"},
		"Cache-Control": {"no-cache, no-store"},
	},
}, {
	Desc:       "POST quiet",
	Method:     "POST",
	Status:     http.StatusOK,
	Params:     url.Values{},
	PostParams: url.Values{},
	Body:       `1`,
	Headers:    http.Header{},
}, {
	Desc:   "POST verbose",
	Method: "POST",
	Status: http.StatusOK,
	Params: url.Values{},
	PostParams: url.Values{
		"verbose": {"1"},
	},
	Body:    `{"error": "", "status": 1}`,
	Headers: http.Header{},
}, {
	Desc:   "POST query verbose",
	Method: "POST",
	Status: http.StatusOK,
	Params: url.Values{
		"verbose": {"1"},
	},
	PostParams: url.Values{},
	Body:       `{"error": "", "status": 1}`,
	Headers:    http.Header{},
}, {
	Desc:   "POST redirect",
	Method: "POST",
	Status: http.StatusFound,
	Params: url.Values{},
	PostParams: url.Values{
		"redirect": {"http://test.example/foo"},
	},
	Body: ``,
	Headers: http.Header{
		"Location":      {"http://test.example/foo"},
		"Cache-Control": {"no-cache, no-store"},
	},
}}

func TestServeDummy(t *testing.T) {
	for _, test := range serveDummy_tests {
		record := httptest.NewRecorder()
		req := &http.Request{
			Method:   test.Method,
			URL:      &url.URL{Path: "/", RawQuery: test.Params.Encode()},
			PostForm: test.PostParams,
		}
		serveDummy(record, req)
		if got, want := record.Code, test.Status; got != want {
			t.Errorf("%s: response code = %d, want %d", test.Desc, got, want)
		}
		if got, want := string(record.Body.Bytes()), test.Body; got != want {
			t.Errorf("%s: body = `%s`, want `%s`", test.Desc, got, want)
		}
		for header, want := range test.Headers {
			if got := record.HeaderMap.Get(header); got != want[0] {
				t.Errorf("%s: header %s = `%s`, want `%s`", test.Desc, header, got, want)
			}
		}
	}
}
