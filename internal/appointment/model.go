package appointment

import "time"

type Appointment struct {
	ID          int
	PatientID   int
	DentistID   int
	Date        time.Time
	Description string
}
