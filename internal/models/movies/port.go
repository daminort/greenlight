package movies

import "greenlight.damian.net/internal/filters"

type RepositoryInstance interface {
	GetList(params GetMoviesParams) ([]Movie, *filters.Meta, error)
	Get(id int64) (*Movie, error)
	Create(movie *Movie) error
	Update(movie *Movie) error
	Delete(id int64) error
}

type ServiceInstance interface {
	GetList(params GetMoviesParams) ([]Movie, *filters.Meta, error)
	Get(id int64) (*Movie, error)
	Create(movie *Movie) error
	Update(movie *Movie) error
	Delete(id int64) error
}
