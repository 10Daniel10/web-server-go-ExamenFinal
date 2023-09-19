package database

import (
	"errors"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	model "github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"
	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewOtherAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (a *AppointmentRepository) GetAll() ([]model.Appointment, error) {
	var data []model.Appointment
	query := a.db.Find(&data)
	if query.Error != nil {
		return nil, internal.ErServiceUnavailable
	}
	return data, nil
}

func (a *AppointmentRepository) GetByID(id uint) (model.Appointment, error) {
	var data model.Appointment
	query := a.db.First(&data, id)
	if query.Error != nil {
		switch {
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return model.Appointment{}, internal.ErNotFound
		}
		return model.Appointment{}, internal.ErServiceUnavailable
	}
	return data, nil
}

func (a *AppointmentRepository) GetByDNI(dni string) (model.Appointment, error) {
	var data model.Appointment
	query := a.db.
		Model(&model.Appointment{}).
		Select("appointments.*").
		Joins("JOIN patients ON appointments.patient_id = patients.id").
		Where("patients.dni = ?", dni).Scan(&data)

	if query.Error != nil {
		switch {
		default:
			return model.Appointment{}, internal.ErServiceUnavailable
		}
	}
	if query.RowsAffected == 0 {
		return model.Appointment{}, internal.ErNotFound
	}

	return data, nil
}

func (a *AppointmentRepository) Create(appointment model.Appointment) (model.Appointment, error) {
	query := a.db.Create(&appointment)
	if query.Error != nil {
		switch {
		default:
			return model.Appointment{}, internal.ErServiceUnavailable
		}
	}
	return appointment, nil
}

func (a *AppointmentRepository) Update(appointment model.Appointment) (model.Appointment, error) {
	query := a.db.Save(&appointment)
	if query.Error != nil {
		return model.Appointment{}, query.Error
	}

	return appointment, nil
}

func (a *AppointmentRepository) Delete(id uint) error {
	query := a.db.Delete(&model.Appointment{}, id)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
