package internal

import "errors"

var (
	/* General errors */

	ErNotFound           = errors.New("not found")
	ErServiceUnavailable = errors.New("service unavailable, try again later")

	/* Dentist errors */

	ErLicenseAlreadyExists = errors.New("license already exists")

	/* Patient errors */

	ErDniAlreadyExists = errors.New("dni already exists")
)
