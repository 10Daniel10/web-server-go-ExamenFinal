package dentist

type Repository interface {
	Create(dentist Dentist) (Dentist, error)
	GetAll() ([]Dentist, error)
	GetByID(id int) (Dentist, error)
	Update(dentist Dentist) (Dentist, error)
	Delete(id int) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(dentist Dentist) (Dentist, error) {
	return s.repository.Create(dentist)
}

func (s *Service) GetAll() ([]Dentist, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id int) (Dentist, error) {
	return s.repository.GetByID(id)
}

func (s *Service) Update(dentist Dentist) (Dentist, error) {
	return s.repository.Update(dentist)
}

func (s *Service) Delete(id int) error {
	return s.repository.Delete(id)
}
