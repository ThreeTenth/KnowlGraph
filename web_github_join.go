package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"knowlgraph.com/ent/user"
)

func joinGithub(c *Context) error {
	////////////// 1 GitHub OAuth2 get token ///////////////
	// 1.1 Request a user's GitHub identity
	// See https://docs.github.com/en/developers/apps/authorizing-oauth-apps#1-request-a-users-github-identity
	//
	// 1.2 Users are redirected back to your site by GitHub
	// See https://docs.github.com/en/developers/apps/authorizing-oauth-apps#2-users-are-redirected-back-to-your-site-by-github
	////////////////////////////////////////////////////////

	githubOAuthCode := c.Query(QueryCode)

	githubOAuthRequestBody := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		State        string `json:"state,omitempty"`
	}{
		ClientID:     config.Gci,
		ClientSecret: config.Gcs,
		Code:         githubOAuthCode,
	}

	githubOAuthState, _ := c.Cookie("github_oauth_state")
	if githubOAuthState != "" {
		githubOAuthRequestBody.State = githubOAuthState
	}

	githubOAuthRequestBodyJSON, err := json.Marshal(githubOAuthRequestBody)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	_httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(githubOAuthRequestBodyJSON))
	if err != nil {
		return c.InternalServerError(err.Error())
	}
	req.Header.Set(HeaderAccept, binding.MIMEJSON)
	req.Header.Set(HeaderContentType, binding.MIMEJSON)
	resp, err := _httpClient.Do(req)
	if err != nil {
		return c.Unauthorized(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.InternalServerError(err.Error())
	}
	var githubAccessToken struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	err = json.Unmarshal(body, &githubAccessToken)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	/////////////// 2 Get GitHub OAuth2 user ///////////////
	// Use the access token to access the API
	// see https://docs.github.com/en/rest/overview/resources-in-the-rest-api#authentication
	// see https://docs.github.com/en/rest/reference/users
	////////////////////////////////////////////////////////

	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	req.Header.Set(HeaderAuthorization, githubAccessToken.TokenType+" "+githubAccessToken.AccessToken)
	resp, err = _httpClient.Do(req)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	var githubUser struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}

	err = json.Unmarshal(body, &githubUser)
	if err != nil || 0 == githubUser.ID || "" == githubUser.Login {
		return c.Unauthorized(string(body))
	}

	_userID, err := client.User.Update().
		Where(user.GithubIDEQ(githubUser.ID)).
		SetName(githubUser.Login).
		SetEmail(githubUser.Email).
		SetAvatar(githubUser.AvatarURL).
		Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}
	if _userID == 0 {
		_user, err := client.User.Create().
			SetName(githubUser.Login).
			SetEmail(githubUser.Email).
			SetAvatar(githubUser.AvatarURL).
			SetGithubID(githubUser.ID).
			Save(ctx)

		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_userID = _user.ID
	}

	if setAuthorization(c.Context, _userID) != nil {
		return c.InternalServerError(err.Error())
	}

	return c.MovedPermanently("/")
}
