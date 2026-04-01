package users

type Service struct {
	Repository RepositoryInstance
}

func NewService(repo RepositoryInstance) *Service {
	return &Service{
		Repository: repo,
	}
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.Repository.GetByEmail(email)
}

func (s *Service) Get(id int64) (*User, error) {
	return s.Repository.Get(id)
}

func (s *Service) Create(u *User) error {
	return s.Repository.Create(u)
}

func (s *Service) Update(u *User) error {
	return s.Repository.Update(u)
}
