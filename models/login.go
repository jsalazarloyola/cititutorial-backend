package models

// Estos son los parámetros para la API
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type ResponseLogin struct {
	Data   map[string]any `json:"data"`
	Expire string         `json:"expire"`
	Token  string         `json:"token"`
}
