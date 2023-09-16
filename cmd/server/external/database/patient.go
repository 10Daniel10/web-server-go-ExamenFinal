package database

import (
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/patient"
	"gorm.io/gorm"
)

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (dr *PatientRepository) Create(patient patient.Patient) (patient.Patient, error) {
	dr.db.Create(&patient)
	return patient, nil
}

func (dr *PatientRepository) GetAll() ([]patient.Patient, error) {
	var data []patient.Patient
	query := dr.db.Find(&data)
	if query.Error != nil {
		return nil, query.Error
	}
	return data, nil
}

func (dr *PatientRepository) GetByID(id uint) (patient.Patient, error) {
	var data patient.Patient
	query := dr.db.First(&data, id)
	if query.Error != nil {
		return data, query.Error
	}
	return data, nil
}

func (dr *PatientRepository) GetByDNI(dni string) (patient.Patient, error) {
	var data patient.Patient
	query := dr.db.Where("dni = ?", dni).First(&data)
	if query.Error != nil {
		return data, query.Error
	}
	return data, nil
}

func (dr *PatientRepository) Update(patient patient.Patient) (patient.Patient, error) {
	dr.db.Save(&patient)
	return patient, nil
}

func (dr *PatientRepository) Delete(id uint) error {
	var data patient.Patient
	query := dr.db.Delete(&data, id)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
