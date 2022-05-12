package model

type Peer struct {
	ID     string `json:"id"`
	Serial int32  `json:"serial"`
	UUID   string `json:"uuid"`
	PK     []byte `json:"pk"`
}
