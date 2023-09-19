package appointment

import (
	"errors"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
)

type Repository interface {
	GetAll() ([]Appointment, error)
	GetByID(id uint) (Appointment, error)
	GetByDNI(dni string) (Appointment, error)
	Create(appointment Appointment) (Appointment, error)
	Update(appointment Appointment) (Appointment, error)
	Delete(id uint) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll() ([]Appointment, error) {
	data, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetByID(id uint) (Appointment, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Appointment{}, internal.ErNotFound

		default:
			return Appointment{}, internal.ErServiceUnavailable
		}
	}
	return data, nil
}

func (s *Service) GetByDNI(dni string) (Appointment, error) {
	data, err := s.repository.GetByDNI(dni)
	if err != nil {
		return Appointment{}, err
	}
	return data, nil
}

func (s *Service) Create(appointment Appointment) (Appointment, error) {

	appointmentCreated, err := s.repository.Create(appointment)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Appointment{}, internal.ErNotFound

		default:
			return Appointment{}, internal.ErServiceUnavailable
		}
	}

	return appointmentCreated, nil
}

func (s *Service) Update(appointment Appointment) (Appointment, error) {
	appointmentSearched, err := s.repository.GetByID(appointment.ID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Appointment{}, internal.ErNotFound

		default:
			return Appointment{}, internal.ErServiceUnavailable
		}
	}

	appointmentUpdated, err := s.repository.Update(appointmentSearched)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Appointment{}, internal.ErNotFound

		default:
			return Appointment{}, internal.ErServiceUnavailable
		}
	}

	return appointmentUpdated, nil
}

func (s *Service) Patch(appointment Appointment) (Appointment, error) {

	appointmentSearched, err := s.repository.GetByID(appointment.ID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Appointment{}, internal.ErNotFound

		default:
			return Appointment{}, internal.ErServiceUnavailable
		}
	}

	CompareTo(&appointment, appointmentSearched)

	appointmentUpdated, err := s.repository.Update(appointmentSearched)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			return Appointment{}, internal.ErNotFound

		default:
			return Appointment{}, internal.ErServiceUnavailable
		}
	}

	return appointmentUpdated, nil
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

func CompareTo(a *Appointment, b Appointment) {
	if a.PatientID == 0 {
		a.PatientID = b.PatientID
	}
	if a.DentistID == 0 {
		a.DentistID = b.DentistID
	}
	if a.Date.IsZero() {
		a.Date = b.Date
	}
	if a.Description == "" {
		a.Description = b.Description
	}
}
