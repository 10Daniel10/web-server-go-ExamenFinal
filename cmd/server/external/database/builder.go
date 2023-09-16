package database

import (
	"fmt"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/patient"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConnectionParams struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func Connect(params ConnectionParams) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		params.User, params.Password, params.Host, params.Port, params.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&dentist.Dentist{}, &patient.Patient{}, &appointment.Appointment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
