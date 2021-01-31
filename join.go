package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func joinGithub(c *Context) error {
	githubOAuthCode := c.Query("code")

	params := url.Values{"client_id": {config.Gci}, "client_secret": {config.Gcs}, "code": {githubOAuthCode}}

	githubOAuthState, _ := c.Cookie("github_oauth_state")
	if githubOAuthState != "" {
		params.Add("state", githubOAuthState)
	}

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", params)
	if err != nil {
		return c.Unauthorized(err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return c.Ok(&body)
}
