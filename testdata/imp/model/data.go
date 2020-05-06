package model

type User struct {
	FirstName      string
	Email          string
	FavoriteColors []string
	RawContent     string
	Id             int
}

type Story struct {
	StoryId  int
	UserId   int
	UserName string
}
