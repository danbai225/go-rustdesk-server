package model

import "time"

type Peer struct {
	Uid         string     `json:"_id"`
	ID          string     `json:"id"`
	IP          string     `json:"ip"`
	UUID        string     `json:"uuid"`
	PK          []byte     `json:"pk"`
	User        string     `json:"user"`
	Disabled    bool       `json:"disabled"`
	LastRegTime *time.Time `json:"last_reg_time"`
}
type Relay struct {
	Uid         string     `json:"_id"`
	Name        string     `json:"name"`
	Port        uint16     `json:"port"`
	IP          string     `json:"ip"`
	Online      bool       `json:"online"`
	LastRegTime *time.Time `json:"last_reg_time"`
	Download    uint       `json:"download"`
	Upload      uint       `json:"upload"`
	Ping        uint       `json:"ping"`
	Cpu         uint       `json:"cpu"`
	NetFlow     float64    `json:"net_flow"`
}
type User struct {
	Uid      string `json:"_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
