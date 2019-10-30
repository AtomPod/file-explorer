package models

//Token 令牌
type Token struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}
