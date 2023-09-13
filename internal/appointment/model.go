package appointment

import (
	"time"
)

type Appointment struct {
	ID          uint      `gorm:"primaryKey"`
	PatientID   uint      `gorm:"not null"`
	DentistID   uint      `gorm:"not null"`
	Date        time.Time `gorm:"not null;type:datetime(3)"`
	Description string    `gorm:"type:longtext"`
}
