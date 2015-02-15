package main

import (
	"github.com/cenkalti/log"
	"github.com/pior/mixpanelproxy"
	"net/http"
	"net/url"
)

func main() {
	director := func(m *mixpanelproxy.MixpanelRequest) (err error) {
		log.Infof("Proxying: %s", m)
		return nil
	}

	url, _ := url.Parse("http://api.mixpanel.com")
	http.Handle("/", mixpanelproxy.NewProxy(url, director))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
