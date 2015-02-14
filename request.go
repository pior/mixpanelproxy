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
	req *http.Request

	mixpanelForward  bool
	payloadAvailable bool

	token      string
	distinctId string
}

func (m MixpanelRequest) String() string {
	return fmt.Sprintf("MixpanelRequest: distinctId=%")
}

func newMixpanelRequest(req *http.Request) (m MixpanelRequest, err error) {
	m.req = req
	m.mixpanelForward = true // forward by default
	m.payloadAvailable = false

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

		m.token = p.Properties.Token
		m.distinctId = p.Properties.DistinctId

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

		m.token = p.Token
		m.distinctId = p.DistinctId

	default:
		log.Debugf("unkown endpoint: %s", req.URL.Path)
		return
	}

	m.payloadAvailable = true
	log.Debugf("parsed: token=%s distinct_id=%s", m.token, m.distinctId)
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

	log.Debugf(`data: "%s"`, data)
	return
}
