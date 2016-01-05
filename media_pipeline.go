package kurento

type MediaPipeline struct {
	ID        string
	SessionId string
	Endpoints []Endpoint
	c         *Client
}

type Endpoint interface {
	GetID() string
}

func (mp *MediaPipeline) Release() error {
	var res Response

	return mp.c.rpc.Call("release", map[string]interface{}{
		"sessionId": mp.SessionId,
		"object":    mp.ID,
	}, &res)
}
