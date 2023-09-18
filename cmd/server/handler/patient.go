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
	Name          string `json:"name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Address       string `json:"address" binding:"required"`
	DNI           string `json:"dni" binding:"required"`
	Email         string `json:"email" binding:"required"`
	AdmissionDate string `json:"admission_date" binding:"required"`
}

type PatientPut struct {
	Name          string `json:"name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Address       string `json:"address" binding:"required"`
	DNI           string `json:"dni" binding:"required"`
	Email         string `json:"email" binding:"required"`
	AdmissionDate string `json:"admission_date" binding:"required"`
}

type PatientPatch struct {
	Name          string `json:"name"`
	LastName      string `json:"last_name"`
	Address       string `json:"address"`
	DNI           string `json:"dni"`
	Email         string `json:"email"`
	AdmissionDate string `json:"admission_date"`
}

type PatientService interface {
	GetAll() ([]patient.Patient, error)
	GetByID(id uint) (patient.Patient, error)
	GetByDNI(dni string) (patient.Patient, error)
	Create(patient patient.Patient) (patient.Patient, error)
	Update(patient patient.Patient) (patient.Patient, error)
	Patch(patient patient.Patient) (patient.Patient, error)
	Delete(id uint) error
}

type PatientHandler struct {
	service PatientService
}

func NewPatientHandler(service PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

func (p *PatientHandler) GetAll(ctx *gin.Context) {
	data, err := p.service.GetAll()
	if err != nil {
		if errors.Is(err, internal.ErServiceUnavailable) {
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusServiceUnavailable,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return
		}
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

func (p *PatientHandler) GetById(ctx *gin.Context) {
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

	data, err := p.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, internal.ErNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patient with id %d not found", id),
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

func (p *PatientHandler) GetByDNI(ctx *gin.Context) {
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

	patientSearched, err := p.service.GetByDNI(dniQuery)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patient with dni %s not found", dniQuery),
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

	body := PatientResponse{
		Id:            patientSearched.ID,
		Name:          patientSearched.Name,
		LastName:      patientSearched.Lastname,
		Address:       patientSearched.Address,
		DNI:           patientSearched.DNI,
		Email:         patientSearched.Email,
		AdmissionDate: patientSearched.AdmissionDate,
	}

	ctx.JSON(http.StatusOK, body)
}

func (p *PatientHandler) Create(ctx *gin.Context) {
	patientToPost := PatientPost{}
	err := ctx.ShouldBindJSON(&patientToPost)
	if err != nil {

		var errs []string
		switch {
		case err.(validator.ValidationErrors) != nil:
			for _, err := range err.(validator.ValidationErrors) {
				errs = append(errs, fmt.Sprintf("%s field is %s",
					extractJSONTag(err.Field(), patientToPost), err.Tag()))
			}
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
	admissionDate, err := time.Parse(time.RFC3339, patientToPost.AdmissionDate)
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

	patientToCreate := patient.Patient{
		Name:          patientToPost.Name,
		Lastname:      patientToPost.LastName,
		Address:       patientToPost.Address,
		DNI:           patientToPost.DNI,
		Email:         patientToPost.Email,
		AdmissionDate: admissionDate,
	}

	patientCreated, err := p.service.Create(patientToCreate)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErDniAlreadyExists):
			ctx.JSON(http.StatusConflict, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusConflict,
				Message:   fmt.Sprintf("dni %s already exists", patientToCreate.DNI),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return
		}
	}

	body := PatientResponse{
		Id:            patientCreated.ID,
		Name:          patientCreated.Name,
		LastName:      patientCreated.Lastname,
		Address:       patientCreated.Address,
		DNI:           patientCreated.DNI,
		Email:         patientCreated.Email,
		AdmissionDate: patientCreated.AdmissionDate,
	}

	ctx.JSON(http.StatusCreated, body)
}

func (p *PatientHandler) Update(ctx *gin.Context) {
	var err error
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

	patientToUpdate := PatientPut{}
	err = ctx.ShouldBindJSON(&patientToUpdate)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s field is %s",
				extractJSONTag(err.Field(), patientToUpdate), err.Tag()))
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
	admissionDate, err := time.Parse(time.RFC3339, patientToUpdate.AdmissionDate)
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

	patientUpdated := patient.Patient{
		ID:            uint(id),
		Name:          patientToUpdate.Name,
		Lastname:      patientToUpdate.LastName,
		Address:       patientToUpdate.Address,
		DNI:           patientToUpdate.DNI,
		Email:         patientToUpdate.Email,
		AdmissionDate: admissionDate,
	}

	patientUpdated, err = p.service.Update(patientUpdated)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patient with id %d not found", id),
				Path:      ctx.Request.URL.Path,
			})
			return
		case errors.Is(err, internal.ErDniAlreadyExists):
			ctx.JSON(http.StatusConflict, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusConflict,
				Message:   fmt.Sprintf("dni %s already exists", patientUpdated.DNI),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return
		}
	}

	body := PatientResponse{
		Id:            patientUpdated.ID,
		Name:          patientUpdated.Name,
		LastName:      patientUpdated.Lastname,
		Address:       patientUpdated.Address,
		DNI:           patientUpdated.DNI,
		Email:         patientUpdated.Email,
		AdmissionDate: patientUpdated.AdmissionDate,
	}

	ctx.JSON(http.StatusOK, body)
}

func (p *PatientHandler) Patch(ctx *gin.Context) {
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

	patientToUpdateParsed := PatientPatch{}
	err = ctx.ShouldBindJSON(&patientToUpdateParsed)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s field is %s",
				extractJSONTag(err.Field(), patientToUpdateParsed), err.Tag()))
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

	var admissionDate time.Time
	if patientToUpdateParsed.AdmissionDate != "" {
		timeLayout := "RFC3339"
		admissionDate, err = time.Parse(time.RFC3339, patientToUpdateParsed.AdmissionDate)
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

	patientToUpdate := patient.Patient{
		ID:            uint(id),
		Name:          patientToUpdateParsed.Name,
		Lastname:      patientToUpdateParsed.LastName,
		Address:       patientToUpdateParsed.Address,
		DNI:           patientToUpdateParsed.DNI,
		Email:         patientToUpdateParsed.Email,
		AdmissionDate: admissionDate,
	}

	patientUpdated, err := p.service.Patch(patientToUpdate)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patient with id %d not found", id),
				Path:      ctx.Request.URL.Path,
			})
			return
		case errors.Is(err, internal.ErDniAlreadyExists):
			ctx.JSON(http.StatusConflict, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusConflict,
				Message:   fmt.Sprintf("dni %s already exists", patientToUpdate.DNI),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return
		}
	}

	body := PatientResponse{
		Id:            patientUpdated.ID,
		Name:          patientUpdated.Name,
		LastName:      patientUpdated.Lastname,
		Address:       patientUpdated.Address,
		DNI:           patientUpdated.DNI,
		Email:         patientUpdated.Email,
		AdmissionDate: patientUpdated.AdmissionDate,
	}

	ctx.JSON(http.StatusOK, body)
}

func (p *PatientHandler) Delete(ctx *gin.Context) {
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

	err = p.service.Delete(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("patient with id %d not found", id),
				Path:      ctx.Request.URL.Path,
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Path:      ctx.Request.URL.Path,
			})
			return
		}
	}

	ctx.JSON(http.StatusNoContent, nil)
}
