package mixpanelproxy

import (
	"github.com/cenkalti/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	ReverseProxy *httputil.ReverseProxy

	Director func(*MixpanelRequest) (err error)
}

func NewProxy(u *string) *Proxy {

	target, err := url.Parse(*u)
	if err != nil {
		log.Fatalf("Invalid host: %v", err)
	}

	if target.Scheme == "" || target.Host == "" {
		log.Fatalf("Invalid host: %v", target)
	}

	log.Infof("Proxying to %s://%s\n", target.Scheme, target.Host)

	httpdirector := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	httpproxy := &httputil.ReverseProxy{Director: httpdirector}

	director := func(mr *MixpanelRequest) (err error) {
		log.Infof("Director: %+v", mr)
		return nil
	}

	return &Proxy{
		ReverseProxy: httpproxy,
		Director:     director,
	}
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m, err := newMixpanelRequest(req)
	if err != nil {
		log.Errorf("Can't parse request: %s %v", req)
		p.ReverseProxy.ServeHTTP(rw, req)
		return
	}

	if p.Director != nil {
		err = p.Director(&m)
		if err != nil {
			log.Errorf("director: %v", req)
			p.ReverseProxy.ServeHTTP(rw, req)
			return
		}
	}

	if m.mixpanelForward == false {
		serveDummyMixpanelResponse(rw, req)
		return
	}

	p.ReverseProxy.ServeHTTP(rw, req)
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
