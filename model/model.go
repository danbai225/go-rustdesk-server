package model

type Peer struct {
	Uid string `json:"_id"`
	ID  string `json:"id"`

	UUID string `json:"uuid"`
	PK   []byte `json:"pk"`
}
