package patient

import (
	"errors"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	"strings"
)

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
		return Patient{}, err
	}

	return data, nil
}

func (s *Service) GetByDNI(dni string) (Patient, error) {
	data, err := s.repository.GetByDNI(dni)
	if err != nil {
		return Patient{}, err
	}

	return data, nil
}

func (s *Service) Create(patient Patient) (Patient, error) {

	patientExist, err := s.repository.GetByDNI(patient.DNI)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			break
		}
	}

	if patientExist.DNI != "" {
		return Patient{}, internal.ErDniAlreadyExists
	}

	Normalize(&patient)

	patientCreated, err := s.repository.Create(patient)
	if err != nil {
		return Patient{}, err
	}

	return patientCreated, nil
}

func (s *Service) Update(patient Patient) (Patient, error) {

	patientSearched, err := s.repository.GetByID(patient.ID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Patient{}, internal.ErNotFound

		default:
			return Patient{}, internal.ErServiceUnavailable
		}
	}

	Normalize(&patient)

	if patient.DNI != patientSearched.DNI {
		patientExist, err := s.repository.GetByDNI(patient.DNI)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErNotFound):
				break
			}
		}

		if patientExist.DNI != "" {
			return Patient{}, internal.ErDniAlreadyExists
		}
	}

	patientUpdated, err := s.repository.Update(patient)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Patient{}, internal.ErNotFound

		default:
			return Patient{}, internal.ErServiceUnavailable
		}
	}

	return patientUpdated, nil
}

func (s *Service) Patch(patient Patient) (Patient, error) {

	patientSearched, err := s.repository.GetByID(patient.ID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Patient{}, internal.ErNotFound

		default:
			return Patient{}, internal.ErServiceUnavailable
		}
	}

	CompareTo(&patient, patientSearched)

	patientUpdated, err := s.repository.Update(patient)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Patient{}, internal.ErNotFound

		default:
			return Patient{}, internal.ErServiceUnavailable
		}
	}

	return patientUpdated, nil
}

func (s *Service) Delete(id uint) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return internal.ErNotFound

		default:
			return internal.ErServiceUnavailable
		}
	}

	err = s.repository.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return internal.ErNotFound

		default:
			return internal.ErServiceUnavailable
		}
	}

	return nil
}

// Custom functions for the service
// Normalize and CompareTo are used to avoid empty fields in the database

func Normalize(patient *Patient) {
	patient.Name = strings.ToLower(patient.Name)
	patient.Lastname = strings.ToLower(patient.Lastname)
	patient.Address = strings.ToLower(patient.Address)
	patient.DNI = strings.ToLower(patient.DNI)
	patient.Email = strings.ToLower(patient.Email)
}

func CompareTo(a *Patient, b Patient) {
	if a.Name == "" {
		a.Name = b.Name
	}

	if a.Lastname == "" {
		a.Lastname = b.Lastname
	}

	if a.Address == "" {
		a.Address = b.Address
	}

	if a.DNI == "" {
		a.DNI = b.DNI
	}

	if a.Email == "" {
		a.Email = b.Email
	}

	if a.AdmissionDate.IsZero() {
		a.AdmissionDate = b.AdmissionDate
	}
}
