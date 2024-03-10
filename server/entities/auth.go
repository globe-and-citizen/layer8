package entities

import "net/url"

// AuthURL represents the URL to redirect the user to for authentication
type AuthURL struct {
	URL          string
	State        string
	Code         string
	CodeVerifier string
}

// String returns the URL to redirect the user to for authentication
func (u *AuthURL) String() string {
	urlc, _ := url.Parse(u.URL)
	q := urlc.Query()
	redirectURI := q.Get("redirect_uri")
	q.Del("client_id")
	q.Del("redirect_uri")
	urlc.RawQuery = q.Encode()
	return redirectURI + urlc.String()
}

// Query returns the query parameters of the AuthURL
func (u *AuthURL) Query() url.Values {
	m, _ := url.ParseQuery(u.URL)
	return m
}
