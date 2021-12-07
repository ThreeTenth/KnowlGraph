package main

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	ua "github.com/mileusna/useragent"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/terminal"
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
	credID, err := base64.StdEncoding.DecodeString(u.Cred.CredID)
	if err != nil {
		return nil
	}
	publicKey, err := base64.StdEncoding.DecodeString(u.Cred.PublicKey)
	if err != nil {
		return nil
	}

	cred := webauthn.Credential{
		ID:              credID,
		PublicKey:       publicKey,
		AttestationType: u.Cred.AttestationType,
	}
	return []webauthn.Credential{cred}
}

var (
	web *webauthn.WebAuthn
)

func initWebAuthn() {
	var err error
	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Knowledge graph", // Display Name for your site
		RPID:          config.Rpid,       // Generally the FQDN for your site
		RPOrigin:      config.Rpog,
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
	t.Code = New16bitID()
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
	if t.OnlyOnce {
		return c.MethodNotAllowed("Unauthorized access to the interface.")
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

	id := 0
	token := New64BitID()
	d := ExpireTimeToken

	terminalMap := make(map[int]string)

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		if 0 == t.UserID {
			user, err1 := tx.User.Create().Save(ctx)
			if err1 != nil {
				return err1
			}
			t.UserID = user.ID
		} else {
			GetV4Redis(RUser(t.UserID), &terminalMap)
		}

		cidBase64 := base64.StdEncoding.EncodeToString(credential.ID)
		publicKeyBase64 := base64.StdEncoding.EncodeToString(credential.PublicKey)
		_credential, err1 := tx.Credential.
			Create().
			SetCredID(cidBase64).
			SetPublicKey(publicKeyBase64).
			SetAttestationType(credential.AttestationType).
			Save(ctx)

		if err1 != nil {
			return err1
		}

		terminal, err1 := tx.Terminal.
			Create().
			SetCode(t.Code).SetName(t.Name).SetUa(t.UA).SetOnlyOnce(false).SetUserID(t.UserID).SetCredential(_credential).
			Save(ctx)

		if err1 != nil {
			return err1
		}

		id = terminal.ID

		terminalMap[id] = token

		_, err1 = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if err1 = SetV2RedisPipe(pipe, RUser(t.UserID), &terminalMap, d); err1 != nil {
				return err1
			}
			pipe.Set(ctx, RToken(token), t.UserID, d)
			pipe.Del(ctx, RChallenge(challenge))
			return nil
		})

		return err1
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(gin.H{"id": id, "name": t.Name, "token": token, "onlyOnce": t.OnlyOnce})
}

// 仅支持用户添加终端时使用
func finishRegOnlyOnce(c *Context) error {
	challenge, t, err := verifyChallenge(c)
	if err != nil {
		return c.StatusError(err)
	}
	if t.State == TokenStateIdle {
		return c.MethodNotAllowed("Unauthorized access to the interface.")
	}

	id := 0
	token := New64BitID()
	d := ExpireTimeTokenOnce
	terminalMap := make(map[int]string)
	if err = GetV4Redis(RUser(t.UserID), &terminalMap); err != nil {
		return c.InternalServerError(err.Error())
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		terminal, err1 := tx.Terminal.
			Create().
			SetCode(t.Code).SetName(t.Name).SetUa(t.UA).SetOnlyOnce(true).SetUserID(t.UserID).
			Save(ctx)

		if err1 != nil {
			return err1
		}

		id = terminal.ID

		terminalMap[id] = token

		_, err1 = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if err1 = SetV2RedisPipe(pipe, RUser(t.UserID), &terminalMap, d); err1 != nil {
				return err1
			}
			pipe.Set(ctx, RToken(token), t.UserID, d)
			pipe.Del(ctx, RChallenge(challenge))
			return nil
		})
		return err1
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(gin.H{"id": id, "name": t.Name, "token": token, "onlyOnce": t.OnlyOnce})
}

func beginWebAuthnLogin(_terminal *ent.Terminal) (*protocol.CredentialAssertion, error) {
	t := Terminal{
		Code: _terminal.Code,
		Name: _terminal.Name,
		Cred: _terminal.Edges.Credential,
	}
	options, sessionData, err := web.BeginLogin(&t)
	if err != nil {
		return nil, &RestfulAPIError{Status: http.StatusInternalServerError, Content: err.Error()}
	}

	if err = SetWebAuthnSession(_terminal.Code, sessionData); err != nil {
		return nil, &RestfulAPIError{Status: http.StatusInternalServerError, Content: err.Error()}
	}
	return options, nil
}

func beginValidate(c *Context) error {
	terminalID := c.QueryInt("id")
	_terminal, err := client.Terminal.Query().Where(terminal.ID(terminalID)).WithCredential().WithUser().First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}
	if _terminal.Edges.User.ID != c.GetInt(GinKeyUserID) {
		return c.Unauthorized("The terminal does not belong to you")
	}

	options, err := beginWebAuthnLogin(_terminal)
	if err != nil {
		return c.StatusError(err)
	}

	return c.Ok(&options)
}

func beginLogin(c *Context) error {
	terminalID := c.QueryInt("id")
	_terminal, err := client.Terminal.Query().Where(terminal.ID(terminalID)).WithCredential().First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	options, err := beginWebAuthnLogin(_terminal)
	if err != nil {
		return c.StatusError(err)
	}

	return c.Ok(&options)
}

func finishWebAuthnLogin(credential io.Reader, _terminal *ent.Terminal) error {
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	sessionData, err := GetWebAuthnSession(_terminal.Code)
	if err != nil {
		return &RestfulAPIError{Status: http.StatusUnauthorized, Content: err.Error()}
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(credential)
	if err != nil {
		return &RestfulAPIError{Status: http.StatusUnauthorized, Content: err.Error()}
	}

	t := Terminal{
		Code: _terminal.Code,
		Name: _terminal.Name,
		Cred: _terminal.Edges.Credential,
	}

	if _, err = web.ValidateLogin(&t, sessionData, parsedResponse); err != nil {
		return &RestfulAPIError{Status: http.StatusUnauthorized, Content: err.Error()}
	}

	return nil
}

func finishValidate(c *Context) error {
	terminalID := c.QueryInt("id")
	_terminal, err := client.Terminal.Query().Where(terminal.ID(terminalID)).WithCredential().WithUser().First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}
	if _terminal.Edges.User.ID != c.GetInt(GinKeyUserID) {
		return c.Unauthorized("The terminal does not belong to you")
	}

	if err = finishWebAuthnLogin(c.Request.Body, _terminal); err != nil {
		return c.StatusError(err)
	}

	return c.Ok(true)
}

func finishLogin(c *Context) error {
	terminalID := c.QueryInt("id")
	_terminal, err := client.Terminal.Query().Where(terminal.ID(terminalID)).WithCredential().WithUser().First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}
	if err = finishWebAuthnLogin(c.Request.Body, _terminal); err != nil {
		return c.StatusError(err)
	}

	token := New64BitID()
	userID := _terminal.Edges.User.ID
	terminalMap := make(map[int]string)
	if err = GetV4Redis(RUser(userID), &terminalMap); err != nil {
		return c.InternalServerError(err.Error())
	}

	d := ExpireTimeToken
	_, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if err1 := SetV2RedisPipe(pipe, RUser(userID), &terminalMap, d); err1 != nil {
			return err1
		}
		pipe.Set(ctx, RToken(token), userID, d)
		return nil
	})
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(gin.H{"id": terminalID, "name": _terminal.Name, "token": token, "onlyOnce": false})
}
