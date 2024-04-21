package models

type Client struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Salt        string `json:"salt"`
}

func CreateClient(id, secret, name, redirect_uri string) Client {
	return Client{
		ID:          id,
		Secret:      secret,
		Name:        name,
		RedirectURI: redirect_uri,
	}
}

func (c *Client) TableName() string {
	return "clients"
}
