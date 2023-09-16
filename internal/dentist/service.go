package dentist

import (
	"errors"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	"strings"
)

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
		return Dentist{}, err
	}

	return data, nil
}

func (s *Service) GetByLicense(license string) (Dentist, error) {
	data, err := s.repository.GetByLicense(license)
	if err != nil {
		return Dentist{}, err
	}

	return data, nil
}

func (s *Service) Create(dentist Dentist) (Dentist, error) {

	dentistExist, err := s.repository.GetByLicense(dentist.License)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			break
		}
	}

	if dentistExist.License != "" {
		return Dentist{}, internal.ErLicenseAlreadyExists
	}

	Normalize(&dentist)

	dentistCreated, err := s.repository.Create(dentist)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			break
		}
	}

	return dentistCreated, nil
}

func (s *Service) Update(dentist Dentist) (Dentist, error) {

	dentistSearched, err := s.repository.GetByID(dentist.ID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Dentist{}, internal.ErNotFound

		default:
			return Dentist{}, internal.ErServiceUnavailable
		}
	}

	Normalize(&dentist)

	if dentistSearched.License != dentist.License {
		dentistExist, err := s.repository.GetByLicense(dentist.License)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErNotFound):
				break
			}
		}

		if dentistExist.License != "" {
			return Dentist{}, internal.ErLicenseAlreadyExists
		}
	}

	dentistUpdated, err := s.repository.Update(dentist)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Dentist{}, internal.ErNotFound

		default:
			return Dentist{}, internal.ErServiceUnavailable
		}
	}

	return dentistUpdated, nil
}

func (s *Service) Patch(dentist Dentist) (Dentist, error) {

	dentistSearched, err := s.repository.GetByID(dentist.ID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Dentist{}, internal.ErNotFound

		default:
			return Dentist{}, internal.ErServiceUnavailable
		}
	}

	CompareTo(&dentist, dentistSearched)

	if dentistSearched.License != dentist.License {
		dentistExist, err := s.repository.GetByLicense(dentist.License)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErNotFound):
				break
			}
		}

		if dentistExist.License != "" {
			return Dentist{}, internal.ErLicenseAlreadyExists
		}
	}

	dentistUpdated, err := s.repository.Update(dentist)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Dentist{}, internal.ErNotFound

		default:
			return Dentist{}, internal.ErServiceUnavailable
		}
	}

	return dentistUpdated, nil
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

func Normalize(dentist *Dentist) {
	dentist.Name = strings.ToLower(dentist.Name)
	dentist.Lastname = strings.ToLower(dentist.Lastname)
	dentist.License = strings.ToLower(dentist.License)

	dentist.Name = strings.TrimSpace(dentist.Name)
	dentist.Lastname = strings.TrimSpace(dentist.Lastname)
	dentist.License = strings.TrimSpace(dentist.License)
}

func CompareTo(a *Dentist, b Dentist) {
	if a.Name == "" {
		a.Name = b.Name
	}

	if a.Lastname == "" {
		a.Lastname = b.Lastname
	}

	if a.License == "" {
		a.License = b.License
	}
}
