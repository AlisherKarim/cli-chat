package env

import (
	"github.com/alisherkarim/cli-chat/types"
)

type Env struct {
	user types.User
}

func (env *Env) GetUser() (user types.User) {
	return env.user
}

func (env *Env) SetUser(username, email string) {
	env.user.Username = username
	env.user.Email = email
}