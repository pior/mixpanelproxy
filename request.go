package mixpanelproxy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/log"
	"net/http"
)

type EventPayload struct {
	Event      string                 `json:"event"`
	Properties EventPropertiesPayload `json:"properties"`
}

type EventPropertiesPayload struct {
	Token      string `json:"token"`
	DistinctId string `json:"distinct_id"`
}

type PeoplePayload struct {
	Token      string `json:"$token"`
	DistinctId string `json:"$distinct_id"`
}

type MixpanelRequest struct {
	req     *http.Request
	forward bool

	Token      string
	DistinctId string
}

func (m MixpanelRequest) String() string {
	return fmt.Sprintf("MixpanelRequest: token=%s distinctId=%s", m.Token, m.DistinctId)
}

func newMixpanelRequest(req *http.Request) (m MixpanelRequest, err error) {
	m.req = req
	m.forward = true // forward by default

	var data []byte

	switch req.URL.Path {
	case "/track", "/track/":
		data, err = decodeRequest(req)
		if err != nil {
			return
		}

		var p EventPayload
		err = json.Unmarshal(data, &p)
		if err != nil {
			return
		}

		m.Token = p.Properties.Token
		m.DistinctId = p.Properties.DistinctId

	case "/engage", "/engage/":
		data, err = decodeRequest(req)
		if err != nil {
			return
		}

		var p PeoplePayload
		err = json.Unmarshal(data, &p)
		if err != nil {
			return
		}

		m.Token = p.Token
		m.DistinctId = p.DistinctId

	default:
		log.Debugf("Not a Mixpanel endpoint: %s", req.URL.Path)
		return
	}

	log.Debugf("Parsed: %s", m)
	return
}

func decodeRequest(req *http.Request) (data []byte, err error) {
	form, err := dumpRequestForm(req)
	if err != nil {
		log.Error("dumpRequestForm: ", err)
		return
	}

	data, err = base64.StdEncoding.DecodeString(form.Get("data"))
	if err != nil {
		log.Error("base64: ", err)
	}

	return
}
