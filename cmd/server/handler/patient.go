package handler

import (
	"errors"
	"fmt"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/patient"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type PatientResponse struct {
	Id            uint      `json:"id"`
	Name          string    `json:"name"`
	LastName      string    `json:"last_name"`
	Address       string    `json:"address"`
	DNI           string    `json:"dni"`
	Email         string    `json:"email"`
	AdmissionDate time.Time `json:"admission_date"`
}

type PatientPost struct {
	Name          string    `json:"name" binding:"required"`
	LastName      string    `json:"last_name" binding:"required"`
	Address       string    `json:"address" binding:"required"`
	DNI           string    `json:"dni" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	AdmissionDate time.Time `json:"admission_date" binding:"required"`
}

type PatientPut struct {
	Name          string    `json:"name" binding:"required"`
	LastName      string    `json:"last_name" binding:"required"`
	Address       string    `json:"address" binding:"required"`
	DNI           string    `json:"dni" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	AdmissionDate time.Time `json:"admission_date" binding:"required"`
}

type PatientPatch struct {
	Name          string    `json:"name"`
	LastName      string    `json:"last_name"`
	Address       string    `json:"address"`
	DNI           string    `json:"dni"`
	Email         string    `json:"email"`
	AdmissionDate time.Time `json:"admission_date"`
}

type PatientHandler struct {
	service *patient.Service
}

func NewPatientHandler(service *patient.Service) *PatientHandler {
	return &PatientHandler{service: service}
}

func (dh *PatientHandler) GetAll(ctx *gin.Context) {
	data, err := dh.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	if len(data) == 0 {
		ctx.JSON(http.StatusOK, data)
		return
	}

	var body []PatientResponse
	for index, item := range data {
		body[index] = PatientResponse{
			Id:            item.ID,
			Name:          item.Name,
			LastName:      item.Lastname,
			Address:       item.Address,
			DNI:           item.DNI,
			Email:         item.Email,
			AdmissionDate: item.AdmissionDate,
		}
	}

	ctx.JSON(http.StatusOK, body)
}

func (dh *PatientHandler) GetById(ctx *gin.Context) {
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

	data, err := dh.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, internal.ErNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   "patient not found",
				Path:      ctx.Request.URL.Path,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	body := PatientResponse{
		Id:            data.ID,
		Name:          data.Name,
		LastName:      data.Lastname,
		Address:       data.Address,
		DNI:           data.DNI,
		Email:         data.Email,
		AdmissionDate: data.AdmissionDate,
	}

	ctx.JSON(http.StatusOK, body)
}

func (dh *PatientHandler) GetByDNI(ctx *gin.Context) {
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

	data, err := dh.service.GetByDNI(dniQuery)
	if err != nil {
		if errors.Is(err, internal.ErNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   "patient not found",
				Path:      ctx.Request.URL.Path,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}
}

func (dh *PatientHandler) Create(ctx *gin.Context) {
	bodyToBind := PatientPost{}
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

	patientToCreate := patient.Patient{
		Name:          bodyToBind.Name,
		Lastname:      bodyToBind.LastName,
		Address:       bodyToBind.Address,
		DNI:           bodyToBind.DNI,
		Email:         bodyToBind.Email,
		AdmissionDate: bodyToBind.AdmissionDate,
	}

	data, err := dh.service.Create(patientToCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	body := PatientResponse{
		Id:            data.ID,
		Name:          data.Name,
		LastName:      data.Lastname,
		Address:       data.Address,
		DNI:           data.DNI,
		Email:         data.Email,
		AdmissionDate: data.AdmissionDate,
	}

	ctx.JSON(http.StatusCreated, body)
}

func (dh *PatientHandler) Update(ctx *gin.Context) {
}

func (dh *PatientHandler) Patch(ctx *gin.Context) {
}

func (dh *PatientHandler) Delete(ctx *gin.Context) {
}
