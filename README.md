# goauthtoken

This CLI tool retrieves OAuth2 access token.

Operating environment:

* Go1.8 or later

Install:
```
go get github.com/nfukasawa/goauthtoken
```

Usage:
```
Usage:
  goauthtoken [OPTIONS]

Application Options:
  -c, --config=   configuration file path
  -s, --save      save OAuth2 token on configuration file / use with --config
  -t, --template  generate configuration template / exclusive with --config

Help Options:
  -h, --help      Show this help message
```

Generate configuration template:
```
goauthtoken -t > path/to/config.json
```

Then, edit configuration file:
```
  "oauth": {
    "client_id": "CLIENT_ID",
    "client_secret": "CLIENT_SECRET",
    "auth_url": "https://example.com/auth",
    "token_url": "https://example.com/token",
    "scopes": [
      "openid",
      "profile",
      "email"
    ],
    "response_type": "code"
  },
  "local_server": {
    "port": 8888
  }
}
```

Run command:
```
goauthtoken -c=path/to/config.json -s
```
