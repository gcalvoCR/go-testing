package main

type UserRepository interface {
	GetUser(id int) (string, error)
}
