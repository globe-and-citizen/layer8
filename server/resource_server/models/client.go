package models

type Client struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
}

func (c *Client) TableName() string {
	return "clients"
}