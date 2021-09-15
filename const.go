package main

import "time"

const (
	// TokenStateIdle is that the token state is idle, and the token is unavailable
	// at this time, and it needs to wait for subsequent authorization.
	// When the idle time exceeds 1 minute, the token becomes invalid.
	TokenStateIdle = 1
	// TokenStateActivated means that the current token is activated and
	// can be used only after authorization.
	// If it is not authorized within 1 minute, the token will become invalid.
	TokenStateActivated = 2
	// TokenStateAuthorized means that the token is authorized,
	// this state is hidden and invisible in the system.
	// When the token is authorized, the token will be associated with the user.
	TokenStateAuthorized = 4
)

const (
	// Version is current version number
	Version = 13
	// VersionName is current version name
	VersionName = "0.0.13"
	// DBVersion is current database version number
	DBVersion = 1
)

const (
	// ExpireTimeToken is token expire time
	ExpireTimeToken = 30 * 24 * time.Hour
	// ExpireTimeTokenOnce is temporary authorization, only once
	ExpireTimeTokenOnce = 2 * time.Hour
	// ExpireTimeChallenge is challenge expire time
	ExpireTimeChallenge = 1 * time.Minute
	// ExpireTimeChallengeConfirm is authorized party's final confirmation time
	ExpireTimeChallengeConfirm = 15 * time.Second
	// ExpireTimeIPInfo is ipinfo expire time
	ExpireTimeIPInfo = 7 * 24 * time.Hour
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
	// GinKeyChallenge means gin context keys `challenge` field
	GinKeyChallenge = "challenge"
)

// Redis Key
const (
	// SF is static file short word
	SF = "sf:"
	// SFVersion is current static file version code
	SFVersion = "sf_version"
)
