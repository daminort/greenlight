package users

type RepositoryInstance interface {
	GetByEmail(email string) (*User, error)
	Get(id int64) (*User, error)
	Create(u *User) error
	Update(u *User) error
}

type ServiceInstance interface {
	GetByEmail(email string) (*User, error)
	Get(id int64) (*User, error)
	Create(u *User) error
	Update(u *User) error
}
