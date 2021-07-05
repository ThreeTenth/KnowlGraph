package main

import "time"

const (
	// Version is current version number
	Version = 7
	// VersionName is current version name
	VersionName = "0.0.7"
	// DBVersion is current database version number
	DBVersion = 1
)

const (
	// ExpireTimeToken is token expire time
	ExpireTimeToken = 30 * 24 * time.Hour
)

// HTML template
const (
	// TplIndexHTML is index html template
	TplIndexHTML = "index.html"
)

const (
	// DefaultLanguage is article default language
	DefaultLanguage = "en"
)

// HTTP URL query key
const (
	// QueryLang is http url query `lang` key
	QueryLang = "lang"
	// QueryCode is http url query `code` key
	QueryCode = "code"
)

// HTTP Header key
const (
	// HeaderAccept means http header `accept` field
	HeaderAccept = "Accept"
	// HeaderContentType means http header `Content-Type` field
	HeaderContentType = "Content-Type"
	// HeaderAuthorization means http header `Authorization` field
	HeaderAuthorization = "Authorization"
	// HeaderAcceptLanguage means http header `Accept-Language`field
	HeaderAcceptLanguage = "Accept-Language"
)

// HTTP cookies key
const (
	// CookieToken means http cookie `access_token`field
	CookieToken = "access_token"
	// CookieUserLang means http cookie `user-lang`field
	CookieUserLang = "user-lang"
)

// Gin set keys
const (
	// GinKeyUserID means gin context keys `access_token` field
	GinKeyUserID = "user_id"
)

// Redis Key
const (
	// SF is static file short word
	SF = "sf:"
	// SFVersion is current static file version code
	SFVersion = "sf_version"
)
