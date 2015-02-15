package mixpanelproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestServeDummy_loop(t *testing.T) {

	tests := []struct {
		Desc    string
		Status  int
		Params  url.Values
		Headers http.Header
		Body    string
	}{{
		Desc:    "simple",
		Status:  http.StatusOK,
		Body:    `1`,
		Headers: http.Header{},
	}, {
		Desc:   "verbose",
		Status: http.StatusOK,
		Params: url.Values{
			"verbose": {"1"},
		},
		Body:    `{"error": "", "status": 1}`,
		Headers: http.Header{},
	}, {
		Desc:   "redirect",
		Status: http.StatusFound,
		Params: url.Values{
			"redirect": {"http://test.example/foo"},
		},
		Body: ``,
		Headers: http.Header{
			"Location":      {"http://test.example/foo"},
			"Cache-Control": {"no-cache, no-store"},
		},
	}}

	for _, test := range tests {
		record := httptest.NewRecorder()
		req := &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/"},
			Form:   test.Params,
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
