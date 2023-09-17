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

type AppointmentPut struct {
	PatientID   uint      `json:"patient_id" binding:"required"`
	DentistID   uint      `json:"dentist_id" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

type AppointmentPatch struct {
	PatientID   uint      `json:"patient_id" binding:"required"`
	DentistID   uint      `json:"dentist_id" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

type AppointmentHandler struct {
	service *appointment.Service
}

func NewAppointmentHandler(service *appointment.Service) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

// OK
func (ah *AppointmentHandler) GetAll(ctx *gin.Context) {
	data, err := ah.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	var body []AppointmentResponse

	for _, item := range data {
		body = append(body, AppointmentResponse{
			Id:          item.ID,
			PatientID:   item.PatientID,
			DentistID:   item.DentistID,
			Date:        item.Date,
			Description: item.Description,
		})
	}

	ctx.JSON(http.StatusOK, body)
}

// OK
func (ah *AppointmentHandler) GetById(ctx *gin.Context) {
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

	data, err := ah.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, internal.ErNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   "appointment not found",
				Path:      ctx.Request.URL.Path,
			})
			return
		}
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

	ctx.JSON(http.StatusOK, body)
}

// TODO
func (ah *AppointmentHandler) Create(ctx *gin.Context) {
	bodyToBind := appointment.AppointmentPost{}
	err := ctx.ShouldBindJSON(&bodyToBind)
	if err != nil {

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("'%s' field is: %s",
				extractJSONTag(err.Field(), bodyToBind), err.Tag()))
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

	appointmentToCreate := appointment.AppointmentPost{
		PatientDNI:     bodyToBind.PatientDNI,
		DentistLicense: bodyToBind.DentistLicense,
		Date:           bodyToBind.Date,
		Description:    bodyToBind.Description,
	}

	data, err := ah.service.Create(appointmentToCreate)
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

/*
func (ah *AppointmentHandler) Update(ctx *gin.Context) {
}

func (ah *AppointmentHandler) Patch(ctx *gin.Context) {
}
*/
func (ah *AppointmentHandler) Delete(ctx *gin.Context) {
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

	err = ah.service.Delete(uint(id))
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
