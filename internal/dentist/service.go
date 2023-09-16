package dentist

import "github.com/10Daniel10/web-server-go-ExamenFinal/internal"

type Repository interface {
	Create(dentist Dentist) (Dentist, error)
	GetAll() ([]Dentist, error)
	GetByID(id uint) (Dentist, error)
	GetByLicense(license string) (Dentist, error)
	Update(dentist Dentist) (Dentist, error)
	Delete(id uint) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Dentist, error) {
	data, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) GetByID(id uint) (Dentist, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Dentist{}, internal.ErrNotFound
		}
	}

	return data, nil
}

func (s *Service) GetByLicense(license string) (Dentist, error) {
	data, err := s.repository.GetByLicense(license)
	if err != nil {
		if err.Error() == "record not found" {
			return Dentist{}, internal.ErrNotFound
		}
	}

	return data, nil
}

func (s *Service) Create(dentist Dentist) (Dentist, error) {

	exists, err := s.repository.GetByLicense(dentist.License)
	if err != nil {
		if err.Error() == "record not found" {
			return Dentist{}, internal.ErrNotFound
		}
	}

	if exists.License == dentist.License {
		return Dentist{}, internal.ErrDuplicate
	}

	data, err := s.repository.Create(dentist)
	if err != nil {
		return Dentist{}, err
	}

	return data, nil
}

func (s *Service) Update(dentist Dentist) (Dentist, error) {
	return s.repository.Update(dentist)
}

func (s *Service) Delete(id uint) error {
	return s.repository.Delete(id)
}
