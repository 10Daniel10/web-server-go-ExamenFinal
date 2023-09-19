package handler

import (
	"errors"
	"fmt"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type DentistResponse struct {
	Id       uint   `json:"id"`
	Lastname string `json:"last_name"`
	Name     string `json:"name"`
	License  string `json:"license"`
}

type DentistPost struct {
	LastName string `json:"last_name" binding:"required"`
	Name     string `json:"name" binding:"required"`
	License  string `json:"license" binding:"required"`
}

type DentistPut struct {
	LastName string `json:"last_name" binding:"required"`
	Name     string `json:"name" binding:"required"`
	License  string `json:"license" binding:"required"`
}

type DentistPatch struct {
	LastName string `json:"last_name"`
	Name     string `json:"name"`
	License  string `json:"license"`
}

type DentistService interface {
	GetAll() ([]dentist.Dentist, error)
	GetByID(id uint) (dentist.Dentist, error)
	GetByLicense(license string) (dentist.Dentist, error)
	Create(dentist dentist.Dentist) (dentist.Dentist, error)
	Update(dentist dentist.Dentist) (dentist.Dentist, error)
	Patch(dentist dentist.Dentist) (dentist.Dentist, error)
	Delete(id uint) error
}

type DentistHandler struct {
	service DentistService
}

func NewDentistHandler(service DentistService) *DentistHandler {
	return &DentistHandler{service: service}
}

func (d *DentistHandler) GetAll(ctx *gin.Context) {
	dentists, err := d.service.GetAll()
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
		return
	}

	if len(dentists) == 0 {
		ctx.JSON(http.StatusOK, dentists)
		return
	}

	var body []DentistResponse
	for _, currentDentist := range dentists {
		body = append(body, DentistResponse{
			Id:       currentDentist.ID,
			Lastname: currentDentist.Lastname,
			Name:     currentDentist.Name,
			License:  currentDentist.License,
		})
	}

	ctx.JSON(http.StatusOK, body)
}

func (d *DentistHandler) GetById(ctx *gin.Context) {
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

	dentistSearched, err := d.service.GetByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})

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

	body := DentistResponse{
		Id:       dentistSearched.ID,
		Lastname: dentistSearched.Lastname,
		Name:     dentistSearched.Name,
		License:  dentistSearched.License,
	}

	ctx.JSON(http.StatusOK, body)
}

func (d *DentistHandler) GetByLicense(ctx *gin.Context) {
	licenseQuery := ctx.Query("license")
	if licenseQuery == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "value of 'license' query param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	dentistSearched, err := d.service.GetByLicense(licenseQuery)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with license %s %s", licenseQuery, err.Error()),
				Path:      ctx.Request.URL.Path,
			})

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

	body := DentistResponse{
		Id:       dentistSearched.ID,
		Lastname: dentistSearched.Lastname,
		Name:     dentistSearched.Name,
		License:  dentistSearched.License,
	}

	ctx.JSON(http.StatusOK, body)
}

func (d *DentistHandler) Create(ctx *gin.Context) {
	dentistToCreated := DentistPost{}
	err := ctx.ShouldBindJSON(&dentistToCreated)
	if err != nil {

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s field is %s",
				extractJSONTag(err.Field(), dentistToCreated), err.Tag()))
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

	dentistToCreate := dentist.Dentist{
		Lastname: dentistToCreated.LastName,
		Name:     dentistToCreated.Name,
		License:  dentistToCreated.License,
	}

	dentistCreated, err := d.service.Create(dentistToCreate)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErLicenseAlreadyExists):
			ctx.JSON(http.StatusConflict, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusConflict,
				Message:   fmt.Sprintf("dentistService with license %s already exists", dentistToCreate.License),
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

	body := DentistResponse{
		Id:       dentistCreated.ID,
		Lastname: dentistCreated.Lastname,
		Name:     dentistCreated.Name,
		License:  dentistCreated.License,
	}

	ctx.JSON(http.StatusCreated, body)
}

func (d *DentistHandler) Update(ctx *gin.Context) {
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

	dentistToUpdate := DentistPut{}
	err = ctx.ShouldBindJSON(&dentistToUpdate)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s field is %s",
				extractJSONTag(err.Field(), dentistToUpdate), err.Tag()))
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

	dentistUpdated := dentist.Dentist{
		ID:       uint(id),
		Lastname: dentistToUpdate.LastName,
		Name:     dentistToUpdate.Name,
		License:  dentistToUpdate.License,
	}

	dentistUpdated, err = d.service.Update(dentistUpdated)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		case errors.Is(err, internal.ErLicenseAlreadyExists):
			ctx.JSON(http.StatusConflict, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusConflict,
				Message:   fmt.Sprintf("dentistService with license %s already exists", dentistToUpdate.License),
				Path:      ctx.Request.URL.Path,
			})
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

	body := DentistResponse{
		Id:       dentistUpdated.ID,
		Lastname: dentistUpdated.Lastname,
		Name:     dentistUpdated.Name,
		License:  dentistUpdated.License,
	}

	ctx.JSON(http.StatusOK, body)
}

func (d *DentistHandler) Patch(ctx *gin.Context) {
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

	dentistToUpdate := DentistPatch{}
	err = ctx.ShouldBindJSON(&dentistToUpdate)
	if err != nil {
		return
	}

	dentistUpdated := dentist.Dentist{
		ID:       uint(id),
		Lastname: dentistToUpdate.LastName,
		Name:     dentistToUpdate.Name,
		License:  dentistToUpdate.License,
	}

	dentistUpdated, err = d.service.Patch(dentistUpdated)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with id %d %s", id, err.Error()),
				Path:      ctx.Request.URL.Path,
			})
			return
		case errors.Is(err, internal.ErLicenseAlreadyExists):
			ctx.JSON(http.StatusConflict, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusConflict,
				Message:   fmt.Sprintf("dentistService with license %s already exists", dentistToUpdate.License),
				Path:      ctx.Request.URL.Path,
			})
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

	body := DentistResponse{
		Id:       dentistUpdated.ID,
		Lastname: dentistUpdated.Lastname,
		Name:     dentistUpdated.Name,
		License:  dentistUpdated.License,
	}

	ctx.JSON(http.StatusOK, body)
}

func (d *DentistHandler) Delete(ctx *gin.Context) {
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

	err = d.service.Delete(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, internal.ErNotFound):
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   fmt.Sprintf("dentistService with id %d %s", id, err.Error()),
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
