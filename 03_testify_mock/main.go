package main

import "fmt"

func main() {
	repo := NewInmemoryRepository()
	service := NewUserService(repo)

	response, err := service.GetUser(1)
	fmt.Println("response", response)
	fmt.Println("error", err)
	fmt.Println("-------------------")
	response, err = service.GetUser(4)
	fmt.Println("response", response)
	fmt.Println("error", err)
}
