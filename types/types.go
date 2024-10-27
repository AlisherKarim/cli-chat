package types

type User struct {
	Username string `json:"username"`
	Email string `json:"email"`
}

type Session struct {
	AccessToken string `json:"access_token"`
}

type ChatListItem struct {
	Name string `json:"name"`
	Id string `json:"id"`
}