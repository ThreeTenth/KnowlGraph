package main

// import (
// 	"encoding/binary"
// 	"encoding/json"
// 	"math/rand"
// 	"net/http"

// 	"github.com/duo-labs/webauthn/protocol"
// 	"github.com/duo-labs/webauthn/webauthn"
// 	"github.com/gorilla/sessions"
// 	"github.com/pkg/errors"
// )

// // WebauthnSession is web authn session name
// const WebauthnSession = "webauthn-session"

// // ErrMarshal is returned if unexpected data is present in a webauthn session.
// var ErrMarshal = errors.New("error unmarshaling data")

// var webAuthn *webauthn.WebAuthn
// var sessionStore *Store

// func init() {
// 	var err error
// 	webAuthn, err = webauthn.New(&webauthn.Config{
// 		RPDisplayName: "KnowlGraph",             // Display Name for your site
// 		RPID:          "localhost",              // Generally the domain name for your site
// 		RPOrigin:      "http://localhost:20010", // The origin URL for WebAuthn requests
// 		// RPIcon: "https://duo.com/logo.png", // Optional icon URL for your site
// 	})

// 	if err != nil {
// 		panicIfErrNotNil(err)
// 	}

// 	sessionStore = NewStore()
// }

// var device *WebAuthnDevice

// func beginRegistration(c *Context) error {
// 	// init user
// 	device = NewUser("", "")

// 	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
// 		credCreationOpts.CredentialExcludeList = device.CredentialExcludeList()
// 		credCreationOpts.AuthenticatorSelection.AuthenticatorAttachment = protocol.Platform
// 	}

// 	// generate PublicKeyCredentialCreationOptions, session data
// 	options, sessionData, err := webAuthn.BeginRegistration(
// 		device,
// 		registerOptions,
// 	)

// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	// store session data as marshaled JSON
// 	err = sessionStore.SaveWebauthnSession("registration", sessionData, c.Request, c.Writer)
// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	return c.Ok(&options)
// }

// func finishRegistration(c *Context) error {

// 	// load the session data
// 	sessionData, err := sessionStore.GetWebauthnSession("registration", c.Request, c.Writer)
// 	if err != nil {
// 		c.BadRequest(err.Error())
// 	}
// 	c.SetCookie(WebauthnSession, "", -1, "/", "", !config.Debug, false)

// 	credential, err := webAuthn.FinishRegistration(device, sessionData, c.Request)
// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	device.AddCredential(*credential)

// 	if setAuthorization(c.Context, 1) != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	return c.Ok(true)
// }

// func beginLogin(c *Context) error {

// 	// get username
// 	// webUser := NewUser("", "")
// 	loginOption := func(credCreationOpts *protocol.PublicKeyCredentialRequestOptions) {
// 		credCreationOpts.UserVerification = protocol.VerificationPreferred
// 	}

// 	// generate PublicKeyCredentialRequestOptions, session data
// 	options, sessionData, err := webAuthn.BeginLogin(device, loginOption)
// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	// store session data as marshaled JSON
// 	err = sessionStore.SaveWebauthnSession("authentication", sessionData, c.Request, c.Writer)
// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	return c.Ok(&options)
// }

// func finishLogin(c *Context) error {
// 	// load the session data
// 	sessionData, err := sessionStore.GetWebauthnSession("authentication", c.Request, c.Writer)
// 	if err != nil {
// 		return c.BadRequest(err.Error())
// 	}
// 	c.SetCookie(WebauthnSession, "", -1, "/", "", !config.Debug, false)

// 	// in an actual implementation, we should perform additional checks on
// 	// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
// 	// and then increment the credentials counter
// 	_, err = webAuthn.FinishLogin(device, sessionData, c.Request)
// 	if err != nil {
// 		return c.BadRequest(err.Error())
// 	}

// 	return c.Ok(true)
// }

// // WebAuthnDevice represents the user model
// type WebAuthnDevice struct {
// 	id          uint64
// 	name        string
// 	displayName string
// 	credentials []webauthn.Credential
// }

// // NewUser creates and returns a new User
// func NewUser(name string, displayName string) *WebAuthnDevice {

// 	user := &WebAuthnDevice{}
// 	user.id = randomUint64()
// 	user.name = name
// 	user.displayName = displayName
// 	// user.credentials = []webauthn.Credential{}

// 	return user
// }

// func randomUint64() uint64 {
// 	buf := make([]byte, 8)
// 	rand.Read(buf)
// 	return binary.LittleEndian.Uint64(buf)
// }

// // WebAuthnID returns the user's ID
// func (u WebAuthnDevice) WebAuthnID() []byte {
// 	buf := make([]byte, binary.MaxVarintLen64)
// 	binary.PutUvarint(buf, uint64(u.id))
// 	return buf
// }

// // WebAuthnName returns the user's username
// func (u WebAuthnDevice) WebAuthnName() string {
// 	return u.name
// }

// // WebAuthnDisplayName returns the user's display name
// func (u WebAuthnDevice) WebAuthnDisplayName() string {
// 	return u.displayName
// }

// // WebAuthnIcon is not (yet) implemented
// func (u WebAuthnDevice) WebAuthnIcon() string {
// 	return ""
// }

// // AddCredential associates the credential to the user
// func (u *WebAuthnDevice) AddCredential(cred webauthn.Credential) {
// 	u.credentials = append(u.credentials, cred)
// }

// // WebAuthnCredentials returns credentials owned by the user
// func (u WebAuthnDevice) WebAuthnCredentials() []webauthn.Credential {
// 	return u.credentials
// }

// // CredentialExcludeList returns a CredentialDescriptor array filled
// // with all the user's credentials
// func (u WebAuthnDevice) CredentialExcludeList() []protocol.CredentialDescriptor {
// 	credentialExcludeList := []protocol.CredentialDescriptor{}
// 	for _, cred := range u.credentials {
// 		descriptor := protocol.CredentialDescriptor{
// 			Type:         protocol.PublicKeyCredentialType,
// 			CredentialID: cred.ID,
// 		}
// 		credentialExcludeList = append(credentialExcludeList, descriptor)
// 	}

// 	return credentialExcludeList
// }

// // =======================================================================

// // Store is a wrapper around sessions.CookieStore which provides some helper
// // methods related to webauthn operations.
// type Store struct {
// 	*sessions.CookieStore
// }

// // NewStore returns a new session store.
// func NewStore(keyPairs ...[]byte) *Store {
// 	// Generate a default encryption key if one isn't provided
// 	if len(keyPairs) == 0 {
// 		key := []byte(RandNdigMbitString(32))
// 		keyPairs = append(keyPairs, key)
// 	}
// 	store := &Store{
// 		sessions.NewCookieStore(keyPairs...),
// 	}
// 	return store
// }

// // Set stores a value to the session with the provided key.
// func (store *Store) Set(key string, value interface{}, r *http.Request, w http.ResponseWriter) error {
// 	session, err := store.Get(r, WebauthnSession)
// 	if err != nil {
// 		return err
// 	}
// 	session.Values[key] = value
// 	session.Save(r, w)
// 	return nil
// }

// // SaveWebauthnSession marhsals and saves the webauthn data to the provided
// // key given the request and responsewriter
// func (store *Store) SaveWebauthnSession(key string, data *webauthn.SessionData, r *http.Request, w http.ResponseWriter) error {
// 	marshaledData, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	return store.Set(key, marshaledData, r, w)
// }

// // GetWebauthnSession unmarshals and returns the webauthn session information
// // from the session cookie.
// func (store *Store) GetWebauthnSession(key string, r *http.Request, w http.ResponseWriter) (webauthn.SessionData, error) {
// 	sessionData := webauthn.SessionData{}
// 	session, err := store.Get(r, WebauthnSession)
// 	if err != nil {
// 		return sessionData, err
// 	}
// 	assertion, ok := session.Values[key].([]byte)
// 	if !ok {
// 		return sessionData, ErrMarshal
// 	}
// 	err = json.Unmarshal(assertion, &sessionData)
// 	if err != nil {
// 		return sessionData, err
// 	}
// 	// Delete the value from the session now that it's been read
// 	delete(session.Values, key)
// 	return sessionData, nil
// }
