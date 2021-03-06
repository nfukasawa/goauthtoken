package oauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"path"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

type Authorization struct {
	AccessToken string
	Code        string
}

func Authorize(ctx context.Context, config *Config) (*Authorization, error) {
	authURL, state := config.AuthCodeURL()
	open.Start(authURL)

	queryCh := make(chan url.Values)
	errorCh := make(chan error)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Server.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			s := q.Get("state")
			b := path.Base(r.URL.Path)
			switch {
			case s == state:
				w.Write([]byte(`<script>window.open('about:blank', '_self').close();</script>`))
				queryCh <- q
			case s == "" && (b == "/" || b == "."):
				w.Write([]byte(`<script>location.href = '/auth_result?' + location.hash.substring(1);</script>`))
			case s != state:
				w.WriteHeader(400)
				errorCh <- fmt.Errorf("invalid callback")
			default:
				w.WriteHeader(404)
			}
		}),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errorCh <- fmt.Errorf("server error")
		}
	}()
	defer srv.Shutdown(ctx)

	select {
	case query := <-queryCh:
		at := query.Get("access_token")
		code := query.Get("code")
		if at == "" && code == "" {
			return nil, fmt.Errorf("accesstoken and authorization code are empty")
		}
		return &Authorization{AccessToken: at, Code: code}, nil
	case err := <-errorCh:
		return nil, err
	}
}

func Exchange(ctx context.Context, config *Config, code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, fmt.Errorf("authorization code is empty")
	}
	return config.OAuth2Config().Exchange(ctx, code)
}

func Refresh(ctx context.Context, config *Config, token *oauth2.Token) (*oauth2.Token, error) {
	if token == nil {
		return nil, fmt.Errorf("token is empty")
	}
	return config.OAuth2Config().TokenSource(ctx, token).Token()
}
