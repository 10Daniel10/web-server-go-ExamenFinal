package database

import (
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
	"gorm.io/gorm"
)

type DentistRepository struct {
	db *gorm.DB
}

func NewDentistRepository(db *gorm.DB) *DentistRepository {
	return &DentistRepository{db: db}
}

func (dr *DentistRepository) Create(dentist dentist.Dentist) (dentist.Dentist, error) {
	dr.db.Create(&dentist)
	return dentist, nil
}

func (dr *DentistRepository) GetAll() ([]dentist.Dentist, error) {
	var data []dentist.Dentist
	query := dr.db.Find(&data)
	if query.Error != nil {
		return nil, query.Error
	}
	return data, nil
}

func (dr *DentistRepository) GetByID(id uint) (dentist.Dentist, error) {
	var data dentist.Dentist
	query := dr.db.First(&data, id)
	if query.Error != nil {
		return data, query.Error
	}
	return data, nil
}

func (dr *DentistRepository) GetByLicense(license string) (dentist.Dentist, error) {
	var data dentist.Dentist
	query := dr.db.Where("license = ?", license).First(&data)
	if query.Error != nil {
		return data, query.Error
	}
	return data, nil
}

func (dr *DentistRepository) Update(dentist dentist.Dentist) (dentist.Dentist, error) {
	dr.db.Save(&dentist)
	return dentist, nil
}

func (dr *DentistRepository) Delete(id uint) error {
	var data dentist.Dentist
	query := dr.db.Delete(&data, id)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
