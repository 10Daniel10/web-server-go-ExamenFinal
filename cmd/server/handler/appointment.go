package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppointmentResponse struct {
	Id          uint      `json:"id"`
	PatientID   uint      `json:"patient_id"`
	DentistID   uint      `json:"dentist_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

type AppointmentDetailResponse struct {
	Id          uint            `json:"id"`
	Patient     PatientResponse `json:"patient"`
	Dentist     DentistResponse `json:"dentist"`
	Date        time.Time       `json:"date"`
	Description string          `json:"description"`
}

type AppointmentPost struct {
	PatientDNI     string `json:"patient_dni" binding:"required"`
	DentistLicense string `json:"dentist_license" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Description    string `json:"description" binding:"required"`
}

type AppointmentPut struct {
	PatientID   uint   `json:"patient_id" binding:"required"`
	DentistID   uint   `json:"dentist_id" binding:"required"`
	Date        string `json:"date" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type AppointmentPatch struct {
	PatientID   uint   `json:"patient_id"`
	DentistID   uint   `json:"dentist_id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type AppointmentService interface {
	GetAll() ([]appointment.Appointment, error)
	GetByID(id uint) (appointment.Appointment, error)
	GetByDNI(dni string) (appointment.Appointment, error)
	Create(appointment appointment.Appointment) (appointment.Appointment, error)
	Update(appointment appointment.Appointment) (appointment.Appointment, error)
	Patch(appointment appointment.Appointment) (appointment.Appointment, error)
	Delete(id uint) error
}

type AppointmentHandler struct {
	service        AppointmentService
	patientService PatientService
	dentistService DentistService
}

func NewAppointmentHandler(service AppointmentService, patient PatientService, dentist DentistService) *AppointmentHandler {
	return &AppointmentHandler{service: service, patientService: patient, dentistService: dentist}
}

func (a *AppointmentHandler) GetAll(ctx *gin.Context) {
	appointments, err := a.service.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErServiceUnavailable):
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
	}

	if len(appointments) == 0 {
		ctx.JSON(http.StatusOK, appointments)
		return
	}

	var body []AppointmentResponse
	for _, currentAppointment := range appointments {
		body = append(body, AppointmentResponse{
			Id:          currentAppointment.ID,
			PatientID:   currentAppointment.PatientID,
			DentistID:   currentAppointment.DentistID,
			Date:        currentAppointment.Date,
			Description: currentAppointment.Description,
		})
	}

	ctx.JSON(http.StatusOK, body)
}

func (a *AppointmentHandler) GetById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param must be a number greater than 0",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	data, err := a.service.GetByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("appointment with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})

			return

		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
	}

	body := AppointmentResponse{
		Id:          data.ID,
		PatientID:   data.PatientID,
		DentistID:   data.DentistID,
		Date:        data.Date,
		Description: data.Description,
	}

	ctx.JSON(http.StatusOK, body)
}

func (a *AppointmentHandler) GetByDNI(ctx *gin.Context) {
	dniQuery := ctx.Query("dni")
	if dniQuery == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "value of 'dni' query param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	appointmentSearched, err := a.service.GetByDNI(dniQuery)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("appointment for patient with dni %s %s", dniQuery, err.Error()),
				Path:      ctx.Request.URL.Path,
			})

		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	patient, err := a.patientService.GetByID(appointmentSearched.PatientID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return

		case errors.Is(err, internal.ErServiceUnavailable):
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   "internal server error, please try again later",
				Path:      ctx.Request.URL.Path,
			})
		}
	}

	dentist, err := a.dentistService.GetByID(appointmentSearched.DentistID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return

		case errors.Is(err, internal.ErServiceUnavailable):
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   "internal server error, please try again later",
				Path:      ctx.Request.URL.Path,
			})
		}
	}

	body := AppointmentDetailResponse{
		Id: appointmentSearched.ID,
		Patient: PatientResponse{
			Id:            patient.ID,
			Name:          patient.Name,
			LastName:      patient.Lastname,
			Address:       patient.Address,
			DNI:           patient.DNI,
			Email:         patient.Email,
			AdmissionDate: patient.AdmissionDate,
		},
		Dentist: DentistResponse{
			Id:       dentist.ID,
			Name:     dentist.Name,
			Lastname: dentist.Lastname,
			License:  dentist.License,
		},
		Date:        appointmentSearched.Date,
		Description: appointmentSearched.Description,
	}

	ctx.JSON(http.StatusOK, body)
}

func (a *AppointmentHandler) Create(ctx *gin.Context) {
	appointmentToPost := AppointmentPost{}
	err := ctx.ShouldBindJSON(&appointmentToPost)
	if err != nil {

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("'%s' field is: %s",
				extractJSONTag(err.Field(), appointmentToPost), err.Tag()))
		}

		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "invalid body",
			Path:      ctx.Request.URL.Path,
			Errors:    errs,
		})
		return
	}

	timeLayout := "RFC3339"
	date, err := time.Parse(time.RFC3339, appointmentToPost.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "invalid body",
			Path:      ctx.Request.URL.Path,
			Errors: []string{
				fmt.Sprintf("admission_date field is invalid"),
				fmt.Sprintf("admission_date field must be in format %s", timeLayout),
			},
		})
		return
	}

	patientExist, err := a.patientService.GetByDNI(appointmentToPost.PatientDNI)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patientService with dni %s %s", appointmentToPost.PatientDNI, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	dentistExist, err := a.dentistService.GetByLicense(appointmentToPost.DentistLicense)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with license %s %s", appointmentToPost.DentistLicense, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	appointmentToCreate := appointment.Appointment{
		PatientID:   patientExist.ID,
		DentistID:   dentistExist.ID,
		Date:        date,
		Description: appointmentToPost.Description,
	}

	data, err := a.service.Create(appointmentToCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	body := AppointmentResponse{
		Id:          data.ID,
		PatientID:   data.PatientID,
		DentistID:   data.DentistID,
		Date:        data.Date,
		Description: data.Description,
	}

	ctx.JSON(http.StatusCreated, body)
}

func (a *AppointmentHandler) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param must be a number greater than 0",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	appointmentToPut := AppointmentPut{}
	err = ctx.ShouldBindJSON(&appointmentToPut)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "invalid body",
			Path:      ctx.Request.URL.Path,
			Errors:    []string{err.Error()},
		})
		return
	}

	timeLayout := "RFC3339"
	date, err := time.Parse(time.RFC3339, appointmentToPut.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "invalid body",
			Path:      ctx.Request.URL.Path,
			Errors: []string{
				fmt.Sprintf("admission_date field is invalid"),
				fmt.Sprintf("admission_date field must be in format %s", timeLayout),
			},
		})
		return
	}

	_, err = a.patientService.GetByID(appointmentToPut.PatientID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patientService with id %d %s", appointmentToPut.PatientID, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	_, err = a.dentistService.GetByID(appointmentToPut.DentistID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with id %d %s", appointmentToPut.DentistID, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	appointmentToUpdate := appointment.Appointment{
		ID:          uint(id),
		PatientID:   appointmentToPut.PatientID,
		DentistID:   appointmentToPut.DentistID,
		Date:        date,
		Description: appointmentToPut.Description,
	}

	appointmentUpdated, err := a.service.Update(appointmentToUpdate)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("appointment with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return

		case errors.Is(err, internal.ErServiceUnavailable):
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
	}

	body := AppointmentResponse{
		Id:          appointmentUpdated.ID,
		PatientID:   appointmentUpdated.PatientID,
		DentistID:   appointmentUpdated.DentistID,
		Date:        appointmentUpdated.Date,
		Description: appointmentUpdated.Description,
	}

	ctx.JSON(http.StatusOK, body)
}

func (a *AppointmentHandler) Patch(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param must be a number greater than 0",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	appointmentToPatch := AppointmentPatch{}
	err = ctx.ShouldBindJSON(&appointmentToPatch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "invalid body",
			Path:      ctx.Request.URL.Path,
			Errors:    []string{err.Error()},
		})
		return
	}

	_, err = a.patientService.GetByID(appointmentToPatch.PatientID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patientService with id %d %s", appointmentToPatch.PatientID, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	_, err = a.dentistService.GetByID(appointmentToPatch.DentistID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with id %d %s", appointmentToPatch.DentistID, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	timeLayout := "RFC3339"
	var date time.Time
	if appointmentToPatch.Date != "" {
		date, err = time.Parse(time.RFC3339, appointmentToPatch.Date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusBadRequest,
				Message:   "invalid body",
				Path:      ctx.Request.URL.Path,
				Errors: []string{
					fmt.Sprintf("admission_date field is invalid"),
					fmt.Sprintf("admission_date field must be in format %s", timeLayout),
				},
			})
			return
		}
	}

	appointmentToUpdate := appointment.Appointment{
		ID:          uint(id),
		PatientID:   appointmentToPatch.PatientID,
		DentistID:   appointmentToPatch.DentistID,
		Date:        date,
		Description: appointmentToPatch.Description,
	}

	appointmentUpdated, err := a.service.Patch(appointmentToUpdate)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("appointment with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return

		case errors.Is(err, internal.ErServiceUnavailable):
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
	}

	body := AppointmentResponse{
		Id:          appointmentUpdated.ID,
		PatientID:   appointmentUpdated.PatientID,
		DentistID:   appointmentUpdated.DentistID,
		Date:        appointmentUpdated.Date,
		Description: appointmentUpdated.Description,
	}

	ctx.JSON(http.StatusOK, body)
}

func (a *AppointmentHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param must be a number greater than 0",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	err = a.service.Delete(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("appointment with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
