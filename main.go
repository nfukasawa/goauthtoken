package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/nfukasawa/goauthtoken/oauth"
	"golang.org/x/oauth2"
)

type Options struct {
	ConfigFilePath string `short:"c" long:"config" description:"configuration file path"`
	SaveConfigFile bool   `short:"s" long:"save" description:"save OAuth2 token on configuration file / use with --config"`
	GenTemplate    bool   `short:"t" long:"template" description:"generate configuration template / exclusive with --config"`
}

func (o Options) isSpecified() bool {
	return (o.ConfigFilePath != "" || o.GenTemplate)
}

func main() {
	var opts Options
	p := flags.NewParser(&opts, flags.Default)
	_, err := p.Parse()
	exitWhenError(err)

	if !opts.isSpecified() {
		p.WriteHelp(os.Stdout)
		return
	}

	if opts.GenTemplate {
		t := oauth.NewConfigTemplate()
		json, err := marshalJSON(t, true)
		exitWhenError(err)
		fmt.Fprintf(os.Stdout, string(json))
		return
	}

	if opts.ConfigFilePath != "" {
		data, err := ioutil.ReadFile(opts.ConfigFilePath)
		exitWhenError(err)

		config := new(oauth.Config)
		err = json.Unmarshal(data, config)
		exitWhenError(err)

		at, token, err := runOAuthFlow(config)
		exitWhenError(err)

		if opts.SaveConfigFile {
			config.CachedToken = token
			json, err := marshalJSON(config, true)
			exitWhenError(err)
			err = ioutil.WriteFile(opts.ConfigFilePath, json, 0644)
			exitWhenError(err)
		}

		fmt.Fprintf(os.Stdout, at)
	}
}

func runOAuthFlow(config *oauth.Config) (accsessToken string, token *oauth2.Token, err error) {
	ctx := context.Background()

	if config.CachedToken != nil {
		token, err := oauth.Refresh(ctx, config, config.CachedToken)
		if err == nil {
			return token.AccessToken, token, nil
		}
	}

	auth, err := oauth.Authorize(ctx, config)
	if err != nil {
		return "", nil, err
	}
	if auth.Code == "" {
		return auth.AccessToken, nil, nil
	}

	token, err = oauth.Exchange(ctx, config, auth.Code)
	if err != nil {
		if auth.AccessToken != "" {
			return auth.AccessToken, nil, nil
		}
		return "", nil, err
	}
	return token.AccessToken, token, nil
}

func marshalJSON(obj interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(obj, "", "  ")
	}
	return json.Marshal(obj)
}

func exitWhenError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}
