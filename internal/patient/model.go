package patient

import (
	model "github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"
	"time"
)

type Patient struct {
	ID            uint                `gorm:"primaryKey"`
	Name          string              `gorm:"not null;type:varchar(60)"`
	Lastname      string              `gorm:"not null;type:varchar(60)"`
	Address       string              `gorm:"not null;type:varchar(120)"`
	DNI           string              `gorm:"not null;unique;type:varchar(20)"`
	Email         string              `gorm:"not null;type:varchar(80)"`
	AdmissionDate time.Time           `gorm:"not null;type:datetime(3)"`
	Appointments  []model.Appointment `gorm:"foreignKey:PatientID"`
}
