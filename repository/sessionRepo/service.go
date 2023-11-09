package sessionRepo

type SessionRepo interface {
	CreateSession(*Session) (*Session, error)
	GetSessionById(id string) (*Session, error)
	GetSessionByCodec(code string) (*Session, error)
	CreateClient(*Client) (*Client, error)
	GetClientById(id string) (*Client, error)
	UpdateClient(*Client) (*Client, error)
}
