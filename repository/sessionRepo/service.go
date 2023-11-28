package sessionRepo

type SessionRepo interface {
	CreateSession(*Session) (*Session, error)
	GetSessionById(id string) (*Session, error)
	GetSessionByCode(code string) (*Session, error)
}
