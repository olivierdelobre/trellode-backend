package models

type LogChange struct {
	Field     string `json:"field"`
	FromValue string `json:"fromValue"`
	ToValue   string `json:"toValue"`
}
