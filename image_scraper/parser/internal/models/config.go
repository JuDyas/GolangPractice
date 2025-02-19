package models

import "time"

const TTL = 1 * time.Minute

type UserAgent struct {
	Device         string `json:"device"`
	UserAgent      string `json:"useragent"`
	Accept         string `json:"accept"`
	AcceptEncoding string `json:"accept_encoding"`
	AcceptLanguage string `json:"accept_language"`
	Referer        string `json:"referer"`
}
