package database

import (
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
)

type DentistDatabase struct {
}

func (db *DentistDatabase) Create(dentist dentist.Dentist) (dentist.Dentist, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DentistDatabase) GetAll() ([]dentist.Dentist, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DentistDatabase) GetByID(id int) (dentist.Dentist, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DentistDatabase) Update(dentist dentist.Dentist) (dentist.Dentist, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DentistDatabase) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
