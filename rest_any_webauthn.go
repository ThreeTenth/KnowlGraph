package main

import (
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	ua "github.com/mileusna/useragent"
)

// WebAuthnID is user ID according to the Relying Party
func (u *Terminal) WebAuthnID() []byte {
	return ([]byte)(u.Code)
}

// WebAuthnName is user Name according to the Relying Party
func (u *Terminal) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName is display Name of the user
func (u *Terminal) WebAuthnDisplayName() string {
	return u.Name
}

// WebAuthnIcon is user's icon url
func (u *Terminal) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials is credentials owned by the user
func (u *Terminal) WebAuthnCredentials() []webauthn.Credential {
	return nil
}

var (
	web *webauthn.WebAuthn
)

func init() {
	var err error
	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Knowledge graph", // Display Name for your site
		RPID:          "localhost",       // Generally the FQDN for your site
		RPOrigin:      "http://localhost",
	})
	panicIfErrNotNil(err)
}

func verifyChallenge(c *Context) (string, *Terminal, error) {
	var form struct {
		State     string `form:"state" binding:"required"`
		Challenge string `form:"challenge" binding:"required"`
	}

	if err := c.ShouldBindQuery(&form); err != nil {
		return "", nil, &RestfulAPIError{Status: http.StatusBadRequest, Content: err.Error()}
	}

	var t Terminal
	if err := GetV4Redis(RChallenge(form.Challenge), &t); err != nil {
		return "", &t, &RestfulAPIError{Status: http.StatusBadRequest, Content: err.Error()}
	}
	if form.State != t.ClientState {
		return "", nil, &RestfulAPIError{Status: http.StatusBadRequest, Content: "Bad state"}
	}
	if 0 != t.UserID && !t.authorized() {
		// 如果已绑定认证账号，但认证账号未授权，则返回“没有权限”

		return "", nil, &RestfulAPIError{Status: http.StatusUnauthorized, Content: "Permission denied"}
	}

	return form.Challenge, &t, nil
}

func beginRegistration(c *Context) error {
	challenge, t, err := verifyChallenge(c)
	if err != nil {
		return c.StatusError(err)
	}

	_ua := ua.Parse(c.GetHeader("User-Agent"))
	userVerification := protocol.VerificationPreferred
	if _ua.IsIOS() {
		userVerification = protocol.VerificationDiscouraged
	}
	t.Code = New16BitID()
	options, sessionData, err := web.BeginRegistration(t, func(pkcco *protocol.PublicKeyCredentialCreationOptions) {
		pkcco.AuthenticatorSelection = protocol.AuthenticatorSelection{
			AuthenticatorAttachment: protocol.Platform,
			UserVerification:        userVerification,
		}
	})
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if err = SetWebAuthnSession(challenge, sessionData); err != nil {
		return c.InternalServerError(err.Error())
	}

	if err = SetV2Redis(RChallenge(challenge), &t, ExpireTimeChallengeConfirm); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(&options)
}

func finishRegistration(c *Context) error {
	challenge, t, err := verifyChallenge(c)
	if err != nil {
		return c.StatusError(err)
	}

	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	sessionData, err := GetWebAuthnSession(challenge)
	if err != nil {
		return c.BadRequest(err.Error())
	}
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	if err != nil {
		return c.Unauthorized(err.Error())
	}
	credential, err := web.CreateCredential(t, sessionData, parsedResponse)
	// Handle validation or input errors
	// If creation was successful, store the credential object

	if err != nil {
		return c.Unauthorized(err.Error())
	}

	fmt.Println(&credential)
	return c.Ok(&credential)
}
