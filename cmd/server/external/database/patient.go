package database

import (
	"errors"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	model "github.com/10Daniel10/web-server-go-ExamenFinal/internal/patient"
	"gorm.io/gorm"
)

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (dr *PatientRepository) Create(patient model.Patient) (model.Patient, error) {
	query := dr.db.Create(&patient)
	if query.Error != nil {
		return patient, query.Error
	}
	return patient, nil
}

func (dr *PatientRepository) GetAll() ([]model.Patient, error) {
	var data []model.Patient
	query := dr.db.Find(&data)
	if query.Error != nil {
		return nil, internal.ErServiceUnavailable
	}
	return data, nil
}

func (dr *PatientRepository) GetByID(id uint) (model.Patient, error) {
	var data model.Patient
	query := dr.db.First(&data, id)
	if query.Error != nil {
		switch {
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return model.Patient{}, internal.ErNotFound
		}
		return model.Patient{}, internal.ErServiceUnavailable
	}
	return data, nil
}

func (dr *PatientRepository) GetByDNI(dni string) (model.Patient, error) {
	var data model.Patient

	query := dr.db.Where("dni = ?", dni).First(&data)
	if query.Error != nil {
		switch {
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return data, internal.ErNotFound
		}
		return model.Patient{}, internal.ErServiceUnavailable
	}

	return data, nil
}

func (dr *PatientRepository) Update(patient model.Patient) (model.Patient, error) {
	query := dr.db.Save(&patient)
	if query.Error != nil {
		return model.Patient{}, query.Error
	}
	return patient, nil
}

func (dr *PatientRepository) Delete(id uint) error {
	var data model.Patient
	query := dr.db.Delete(&data, id)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
