package main

type UserService struct {
	Repo UserRepository
}

func (s *UserService) GetUser(id int) (string, error) {
	return s.Repo.GetUser(id)
}

func NewUserService(repo UserRepository) UserService {
	return UserService{Repo: repo}
}
