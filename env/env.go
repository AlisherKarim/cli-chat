package env

import (
	"github.com/alisherkarim/cli-chat/types"
	tea "github.com/charmbracelet/bubbletea"
)

type Env struct {
	CurrentProgram *tea.Program
	currentUser types.User
	currentSession types.Session
}

func (env *Env) GetUser() (user types.User) {
	return env.currentUser
}

func (env *Env) SetUser(username, email string) {
	env.currentUser.Username = username
	env.currentUser.Email = email
}

func (env *Env) SetSession(access_token string) {
	env.currentSession.AccessToken = access_token
}

func (env *Env) GetSession() types.Session {
	return env.currentSession
}