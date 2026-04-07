package protocol

type Payload struct {
	MsgType Type   `json:"msg_type"`
	Remark  string `json:"remark,omitempty"`
	Data   []byte `json:"data,omitempty"`
	Meta   *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
