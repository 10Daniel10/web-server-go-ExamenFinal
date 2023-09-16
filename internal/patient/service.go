package patient

import "github.com/10Daniel10/web-server-go-ExamenFinal/internal"

type Repository interface {
	Create(patient Patient) (Patient, error)
	GetAll() ([]Patient, error)
	GetByID(id uint) (Patient, error)
	GetByDNI(dni string) (Patient, error)
	Update(patient Patient) (Patient, error)
	Delete(id uint) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(patient Patient) (Patient, error) {

	exists, err := s.repository.GetByDNI(patient.DNI)
	if err != nil {
		if err.Error() == "record not found" {
			return Patient{}, internal.ErrNotFound
		}
	}

	if exists.DNI == patient.DNI {
		return Patient{}, internal.ErrDuplicate
	}

	data, err := s.repository.Create(patient)
	if err != nil {
		return Patient{}, err
	}

	return data, nil
}

func (s *Service) GetAll() ([]Patient, error) {
	data, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) GetByID(id uint) (Patient, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Patient{}, internal.ErrNotFound
		}
	}

	return data, nil
}

func (s *Service) GetByDNI(dni string) (Patient, error) {
	data, err := s.repository.GetByDNI(dni)
	if err != nil {
		if err.Error() == "record not found" {
			return Patient{}, internal.ErrNotFound
		}
	}

	return data, nil
}

func (s *Service) Update(patient Patient) (Patient, error) {
	return s.repository.Update(patient)
}

func (s *Service) Delete(id uint) error {
	return s.repository.Delete(id)
}
