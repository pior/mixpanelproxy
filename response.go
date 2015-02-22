package mixpanelproxy

import (
	"fmt"
	"net/http"
)

// Simulate a success response (handle verbose and redirect)
func serveDummy(rw http.ResponseWriter, req *http.Request) {
	verbose := req.URL.Query().Get("verbose")

	if err := req.ParseForm(); err != nil {
		dummyResponseFail(rw, verbose, "Failed to parse request")
		return
	}

	// Verbose in QS has precedence
	if verbose == "" {
		verbose = req.Form.Get("verbose")
	}

	if redirect := req.Form.Get("redirect"); redirect != "" {
		dummyRedirect(rw, redirect)
	} else {
		dummyResponseSuccess(rw, verbose)
	}
}

func dummyResponseSuccess(rw http.ResponseWriter, verbose string) {
	if verbose == "1" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error": "", "status": 1}`))
	} else {
		rw.Write([]byte("1"))
	}
}

func dummyResponseFail(rw http.ResponseWriter, verbose string, errmsg string) {
	if verbose == "1" {
		rw.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"error": "%s", "status": 0}`, errmsg)
		rw.Write([]byte(msg))
	} else {
		rw.Write([]byte("0"))
	}
}

func dummyRedirect(rw http.ResponseWriter, redirect string) {
	rw.WriteHeader(http.StatusFound)
	rw.Header().Set("Location", redirect)
	rw.Header().Set("Cache-Control", "no-cache, no-store")
}
