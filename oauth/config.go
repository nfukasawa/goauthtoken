package oauth

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/oauth2"
)

const (
	DefaultServerPort = 8888
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	OAuth       *OAuthConfig  `json:"oauth"`
	Server      *ServerConfig `json:"local_server,omitempty"`
	CachedToken *oauth2.Token `json:"cached_token,omitempty"`
}

type OAuthConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret,omitempty"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url,omitempty"`
	Scopes       []string `json:"scopes"`
	ResponseType string   `json:"response_type"`
}

type ServerConfig struct {
	Port int `json:"port,omitempty"`
}

func NewConfigTemplate() *Config {
	return &Config{
		OAuth: &OAuthConfig{
			ClientID:     "CLIENT_ID",
			ClientSecret: "CLIENT_SECRET",
			AuthURL:      "https://example.com/auth",
			TokenURL:     "https://example.com/token",
			Scopes:       []string{"openid", "profile", "email"},
			ResponseType: "code",
		},
		Server: &ServerConfig{
			Port: DefaultServerPort,
		},
	}
}

func (c *Config) AuthCodeURL() (url, state string) {
	buf := make([]byte, 16)
	rand.Read(buf)
	state = fmt.Sprintf("%x", buf)
	url = c.oauth2Config().AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", c.OAuth.ResponseType))
	return
}

func (c *Config) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.OAuth.ClientID,
		ClientSecret: c.OAuth.ClientSecret,
		RedirectURL:  c.redirectURL(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.OAuth.AuthURL,
			TokenURL: c.OAuth.TokenURL,
		},
		Scopes: c.OAuth.Scopes,
	}
}

func (c *Config) redirectURL() string {
	port := DefaultServerPort
	if c.Server != nil && c.Server.Port > 0 {
		port = c.Server.Port
	}
	return fmt.Sprintf("http://localhost:%d", port)
}
