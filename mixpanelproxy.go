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

	if m.forward {
		p.reverseProxy.ServeHTTP(rw, req)
	} else {
		serveDummy(rw, req)
	}
}
