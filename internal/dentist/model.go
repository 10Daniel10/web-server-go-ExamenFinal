package dentist

import model "github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"

type Dentist struct {
	ID           uint                `gorm:"primaryKey"`
	Lastname     string              `gorm:"not null;type:varchar(60)"`
	Name         string              `gorm:"not null;type:varchar(60)"`
	License      string              `gorm:"not null;unique;type:varchar(40)"`
	Appointments []model.Appointment `gorm:"foreignKey:DentistID"`
}
