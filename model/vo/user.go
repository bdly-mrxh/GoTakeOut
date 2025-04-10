package vo

type UserLoginVO struct {
	ID     int    `json:"id"`
	Token  string `json:"token"`
	OpenID string `json:"openid"`
}
