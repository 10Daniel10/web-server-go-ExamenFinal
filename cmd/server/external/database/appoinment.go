package database

import (
	model "github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"
	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (ar *AppointmentRepository) Create(appointment model.Appointment) (model.Appointment, error) {
	ar.db.Create(&appointment)
	return appointment, nil
}

func (ar *AppointmentRepository) GetAll() ([]model.Appointment, error) {
	var data []model.Appointment
	query := ar.db.Find(&data)
	if query.Error != nil {
		return nil, query.Error
	}
	return data, nil
}

func (ar *AppointmentRepository) GetByID(id uint) (model.Appointment, error) {
	var data model.Appointment
	query := ar.db.First(&data, id)
	if query.Error != nil {
		return data, query.Error
	}
	return data, nil
}

func (ar *AppointmentRepository) Update(appointment model.Appointment) (model.Appointment, error) {
	ar.db.Save(&appointment)
	return appointment, nil
}

func (ar *AppointmentRepository) Delete(id uint) error {
	var data model.Appointment
	query := ar.db.Delete(&data, id)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
