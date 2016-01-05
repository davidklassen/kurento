package kurento

import (
	"encoding/json"
	"github.com/DavidKlassen/rpc-codec/jsonrpc2"
	"io"
)

type Client struct {
	rpc *jsonrpc2.Client
}

type Response struct {
	SessionId string `json:"sessionId"`
	Value     string `json:"value"`
}

type Event struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Value struct {
			Data   json.RawMessage `json:"data"`
			Object string          `json:"object"`
			Type   string          `json:"type"`
		} `json:"value"`
	} `json:"params"`
}

func NewClient(ws io.ReadWriteCloser) *Client {
	rpc := jsonrpc2.NewClient(ws)

	return &Client{rpc}
}

func (c *Client) DescribeServerManager() (Response, error) {
	var res Response

	err := c.rpc.Call("describe", map[string]interface{}{
		"object": "manager_ServerManager",
	}, &res)

	return res, err
}

func (c *Client) CreateMediaPipeline(sid string) (*MediaPipeline, error) {
	var res Response

	err := c.rpc.Call("create", map[string]interface{}{
		"sessionId": sid,
		"type":      "MediaPipeline",
	}, &res)

	return &MediaPipeline{res.Value, res.SessionId, []Endpoint{}, c}, err
}

func (c *Client) CreateWebRtcEndpoint(sid string, mp *MediaPipeline) (*WebRtcEndpoint, error) {
	var res Response

	err := c.rpc.Call("create", map[string]interface{}{
		"sessionId": sid,
		"type":      "WebRtcEndpoint",
		"constructorParams": map[string]interface{}{
			"mediaPipeline": mp.ID,
		},
	}, &res)

	if err != nil {
		return nil, err
	}

	endpoint := &WebRtcEndpoint{res.Value, res.SessionId, c}
	mp.Endpoints = append(mp.Endpoints, endpoint)

	return endpoint, err
}

func (c *Client) CreateRecorderEndpoint(sid string, mp *MediaPipeline, uri string) (*RecorderEndpoint, error) {
	var res Response

	err := c.rpc.Call("create", map[string]interface{}{
		"sessionId": sid,
		"type":      "RecorderEndpoint",
		"constructorParams": map[string]interface{}{
			"mediaPipeline": mp.ID,
			"uri":           uri,
		},
	}, &res)

	endpoint := &RecorderEndpoint{res.Value, res.SessionId, c}
	mp.Endpoints = append(mp.Endpoints, endpoint)

	return endpoint, err
}
