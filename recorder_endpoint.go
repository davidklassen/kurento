package kurento

type RecorderEndpoint struct {
	ID        string
	SessionId string
	c         *Client
}

func (re *RecorderEndpoint) GetID() string {
	return re.ID
}

func (re *RecorderEndpoint) Record() error {
	var res Response

	return re.c.rpc.Call("invoke", map[string]interface{}{
		"sessionId": re.SessionId,
		"object":    re.ID,
		"operation": "record",
	}, &res)
}
