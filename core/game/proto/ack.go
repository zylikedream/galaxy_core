package proto

type AckMessage struct {
	Code         int    `json:"code"`
	Reason       string `json:"reason,omitempty"`
	ResponseData []byte `json:"response_data,omitempty"`
}
