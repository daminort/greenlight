package movies

import (
	"greenlight.damian.net/internal/pkg/filters"
)

type Service struct {
	Repository RepositoryInstance
}

func NewService(repo RepositoryInstance) *Service {
	return &Service{
		Repository: repo,
	}
}

func (s *Service) GetList(params GetMoviesParams) ([]Movie, *filters.Meta, error) {
	return s.Repository.GetList(params)
}

func (s *Service) Get(id int64) (*Movie, error) {
	return s.Repository.Get(id)
}

func (s *Service) Create(movie *Movie) error {
	return s.Repository.Create(movie)
}

func (s *Service) Update(movie *Movie) error {
	return s.Repository.Update(movie)
}

func (s *Service) Delete(id int64) error {
	return s.Repository.Delete(id)
}
