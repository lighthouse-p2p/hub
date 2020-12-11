package models

// RegisterBody is used to deserialize the HTTP register endpoint's body
type RegisterBody struct {
	PubKey   string `json:"pub_key"`
	NickName string `json:"nickname"`
}
