package mixpanelproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

func dumpRequestForm(req *http.Request) (form url.Values, err error) {
	savedBody := req.Body

	if req.Body != nil { // read the body and make 2 copies
		var buf bytes.Buffer
		if _, err = buf.ReadFrom(req.Body); err != nil {
			return
		}
		if err = req.Body.Close(); err != nil {
			return
		}
		savedBody = ioutil.NopCloser(&buf)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
	}

	// Read the form
	err = req.ParseForm()
	if err != nil {
		return
	}
	form = req.Form

	if req.Body != nil {
		req.Body = savedBody
	}
	return
}
