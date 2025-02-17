package models

type UserAgent struct {
	Device         string `json:"device"`
	UserAgent      string `json:"useragent"`
	Accept         string `json:"accept"`
	AcceptEncoding string `json:"accept_encoding"`
	AcceptLanguage string `json:"accept_language"`
	Referer        string `json:"referer"`
}
