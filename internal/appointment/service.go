package appointment

import (
	"time"

	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
)

type Repository interface {
	GetByID(id uint) (Appointment, error)
	//Update(appointment Appointment) (Appointment, error)
	Delete(id uint) error
	Create(appointment Appointment) (Appointment, error)
	GetAll() ([]Appointment, error)
}

type AppointmentPost struct {
	PatientDNI     uint      `json:"patient_dni" binding:"required"`
	DentistLicense uint      `json:"dentist_license" binding:"required"`
	Date           time.Time `json:"date" binding:"required"`
	Description    string    `json:"description" binding:"required"`
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
		if err.Error() == "record not found" {
			return Appointment{}, internal.ErNotFound
		}
	}
	return data, nil
}

func (s *Service) Create(appointment AppointmentPost) (Appointment, error) {
	//TODO:
	/*
		dentist, err := s.dentistService.GetByLicense(dentistLicense)
		if err != nil {
			if err.Error() == "record not found" {
				return Appointment{}, internal.ErNotFound
			}
		}

		patient, err := s.patientService.GetByDNI(patientDNI)
		if err != nil {
			if err.Error() == "record not found" {
				return Appointment{}, internal.ErNotFound
			}
		}

		appointmentToCreate := Appointment{
			PatientID:   patient.ID,
			DentistID:   dentist.ID,
			Date:        appointment.Date,
			Description: appointment.Description,
		}

		data, err := s.repository.Create(appointmentToCreate)
		if err != nil {
			return Appointment{}, err
		}

		return data, nil
	*/
	return Appointment{}, nil
}

/*
	func (s *Service) Update(appointment Appointment) (Appointment, error) {
		return s.repository.Update(appointment)
	}
*/
func (s *Service) Delete(id uint) error {
	return s.repository.Delete(id)
}
