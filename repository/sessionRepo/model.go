package sessionRepo

import "time"

type Client struct {
	Id          string
	Name        string
	RedirectUri string
	Secret      string
}

type Session struct {
	Id            string
	RemoteAddr    string
	ClientId      string
	Code          string
	CodeChallenge string
	State         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
