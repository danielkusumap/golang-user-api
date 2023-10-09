package models

type User struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
