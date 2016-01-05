package kurento

import (
	"encoding/json"
	"log"
)

type WebRtcEndpoint struct {
	ID        string
	SessionId string
	c         *Client
}

type IceCandidateEventData struct {
	Candidate json.RawMessage `json:"candidate"`
	Source    string          `json:"source"`
	Tags      []string        `json:"tags"`
	Timestamp string          `json:"timestamp"`
	Type      string          `json:"type"`
}

type IceCandidate struct {
	Module        string `json:"__module__"`
	Type          string `json:"__type__"`
	Candidate     string `json:"candidate"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	SdpMid        string `json:"sdpMid"`
}

func (we *WebRtcEndpoint) GetID() string {
	return we.ID
}

func (we *WebRtcEndpoint) Connect(e Endpoint) error {
	var res Response

	return we.c.rpc.Call("invoke", map[string]interface{}{
		"sessionId": we.SessionId,
		"object":    we.ID,
		"operation": "connect",
		"operationParams": map[string]interface{}{
			"sink": e.GetID(),
		},
	}, &res)
}

func (we *WebRtcEndpoint) ProcessOffer(offer string) (Response, error) {
	var res Response

	err := we.c.rpc.Call("invoke", map[string]interface{}{
		"sessionId": we.SessionId,
		"object":    we.ID,
		"operation": "processOffer",
		"operationParams": map[string]interface{}{
			"offer": offer,
		},
	}, &res)

	return res, err
}

func (we *WebRtcEndpoint) GatherCandidates() error {
	var res Response

	return we.c.rpc.Call("invoke", map[string]interface{}{
		"sessionId": we.SessionId,
		"object":    we.ID,
		"operation": "gatherCandidates",
	}, &res)
}

func (we *WebRtcEndpoint) SubscribeToOnIceCandidate() (<-chan IceCandidate, error) {
	var res Response

	ch := we.c.rpc.GetUnhandledChannel()
	out := make(chan IceCandidate)

	go (func() {
		for {
			raw := <-ch
			e := Event{}
			err := json.Unmarshal(raw.([]byte), &e)

			if err != nil {
				continue
			}

			versionOk := e.Version == "2.0"
			methodOk := e.Method == "onEvent"
			typeOk := e.Params.Value.Type == "OnIceCandidate"
			sourceOk := e.Params.Value.Object == we.ID

			log.Printf("preparing ice candidate: %+v\n", e)

			if !versionOk || !methodOk || !typeOk || !sourceOk {
				continue
			}

			icd := IceCandidateEventData{}
			err = json.Unmarshal(e.Params.Value.Data, &icd)

			if err != nil {
				continue
			}

			ic := IceCandidate{}
			err = json.Unmarshal(icd.Candidate, &ic)

			if err != nil {
				continue
			}

			log.Printf("sending ice candidate: %+v\n", ic)
			out <- ic
		}
	})()

	err := we.c.rpc.Call("subscribe", map[string]interface{}{
		"sessionId": we.SessionId,
		"object":    we.ID,
		"type":      "OnIceCandidate",
	}, &res)

	return out, err
}

func (we *WebRtcEndpoint) AddIceCandidate(c IceCandidate) error {
	var res Response

	c.Module = "kurento"
	c.Type = "IceCandidate"

	return we.c.rpc.Call("invoke", map[string]interface{}{
		"sessionId": we.SessionId,
		"object":    we.ID,
		"operation": "addIceCandidate",
		"operationParams": map[string]interface{}{
			"candidate": c,
		},
	}, &res)
}
