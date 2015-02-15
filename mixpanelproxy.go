package mixpanelproxy

import (
	"github.com/cenkalti/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Director func(*MixpanelRequest) (err error)

type Proxy struct {
	reverseProxy *httputil.ReverseProxy
	director     Director
}

func NewProxy(target *url.URL, director Director) *Proxy {
	httpDirector := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	httpproxy := &httputil.ReverseProxy{Director: httpDirector}

	return &Proxy{
		reverseProxy: httpproxy,
		director:     director,
	}
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m, err := newMixpanelRequest(req)
	if err != nil {
		log.Errorf("Decoding failed: %+v", req)
		p.reverseProxy.ServeHTTP(rw, req)
		return
	}

	if p.director != nil {
		err = p.director(&m)
		if err != nil {
			log.Errorf("director: %v", req)
			p.reverseProxy.ServeHTTP(rw, req)
			return
		}
	}

	if m.mixpanelForward == false {
		serveDummyMixpanelResponse(rw, req)
		return
	}

	p.reverseProxy.ServeHTTP(rw, req)
}

// Simulate a success response (handle verbose and redirect)
func serveDummyMixpanelResponse(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Errorf("%v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	redirect := req.Form.Get("redirect")
	if redirect != "" {
		rw.WriteHeader(302)
		rw.Header().Set("Location", redirect)
		rw.Header().Set("Cache-Control", "no-cache, no-store")
		return
	}

	verbose := req.Form.Get("verbose")
	if verbose == "1" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error": "", "status": 1}`))
	} else {
		rw.Write([]byte("1"))
	}
}
