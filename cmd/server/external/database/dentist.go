package database

import (
	"errors"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	model "github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
	"gorm.io/gorm"
)

type DentistRepository struct {
	db *gorm.DB
}

func NewDentistRepository(db *gorm.DB) *DentistRepository {
	return &DentistRepository{db: db}
}

func (d *DentistRepository) Create(dentist model.Dentist) (model.Dentist, error) {
	query := d.db.Create(&dentist)
	if query.Error != nil {
		return model.Dentist{}, query.Error
	}
	return dentist, nil
}

func (d *DentistRepository) GetAll() ([]model.Dentist, error) {
	var data []model.Dentist
	query := d.db.Find(&data)
	if query.Error != nil {
		return nil, internal.ErServiceUnavailable
	}
	return data, nil
}

func (d *DentistRepository) GetByID(id uint) (model.Dentist, error) {
	var data model.Dentist
	query := d.db.First(&data, id)
	if query.Error != nil {
		switch {
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return model.Dentist{}, internal.ErNotFound
		}
		return model.Dentist{}, internal.ErServiceUnavailable
	}
	return data, nil
}

func (d *DentistRepository) GetByLicense(license string) (model.Dentist, error) {
	var data model.Dentist

	query := d.db.Where("license = ?", license).First(&data)
	if query.Error != nil {
		switch {
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return data, internal.ErNotFound
		}
		return model.Dentist{}, internal.ErServiceUnavailable
	}

	return data, nil
}

func (d *DentistRepository) Update(dentist model.Dentist) (model.Dentist, error) {
	query := d.db.Save(&dentist)
	if query.Error != nil {
		return model.Dentist{}, query.Error
	}
	return dentist, nil
}

func (d *DentistRepository) Delete(id uint) error {
	query := d.db.Delete(&model.Dentist{}, id)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
