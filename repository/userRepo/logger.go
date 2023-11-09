package userRepo

import (
	"time"

	"github.com/go-kit/log"
)

type logLevel int

const (
	None logLevel = iota
	Info
	Error
	Debug
)

type userRepoLogger struct {
	logger log.Logger
	next   UserRepo
	level  logLevel
}

func (url *userRepoLogger) Create(up *UserPayload) (out *User, err error) {
	defer url.Log("create user", *up, *out, err, time.Now())
	return url.next.Create(up)
}

func (url *userRepoLogger) Read(id string) (out *User, err error) {
	defer url.Log("read user", id, *out, err, time.Now())
	return url.next.Read(id)
}

func (url *userRepoLogger) Search(us *UserSearch) (out []User, err error) {
	defer url.Log("search user", *us, out, err, time.Now())
	return url.next.Search(us)

}

func (url *userRepoLogger) Update(up *User) (out *User, err error) {
	defer url.Log("update user", *up, *out, err, time.Now())
	return url.next.Update(up)
}

func (url *userRepoLogger) Delete(id string) (err error) {
	defer url.Log("delete user", "id: "+string(id), "", err, time.Now())
	return url.next.Delete(id)
}

func (url *userRepoLogger) Log(method string, input interface{}, output interface{}, err error, begin time.Time) {

	if url.level == logLevel(None) {
		return
	}
	url.logger.Log("method", method, "dateTime", begin, "duration", time.Since(begin))
	if url.level == logLevel(Info) {
		return
	}
	if err != nil {
		url.logger.Log("error", err.Error())
	}
	if url.level == logLevel(Error) {
		return
	}
	url.logger.Log("in", input, "out", output)
	return
}
