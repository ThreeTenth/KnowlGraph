package main

// HTML template
const (
	// TplIndexHTML is index html template
	TplIndexHTML = "index.html"
	TplMainHTML  = "main.html"
)

const (
	// DefaultLanguage is story default language
	DefaultLanguage = "en"
)

// HTTP Header key
const (
	// HeaderAccept means http header `accept` field
	HeaderAccept = "Accept"
	// HeaderContentType means http header `Content-Type` field
	HeaderContentType = "Content-Type"
	// HeaderAuthorization means http header `Authorization` field
	HeaderAuthorization = "Authorization"
)

// HTTP cookies key
const (
	// CookieToken means http cookie `access_token`field
	CookieToken = "access_token"
)

// Gin set keys
const (
	// GinKeyUserID means gin context keys `access_token` field
	GinKeyUserID = "user_id"
)
