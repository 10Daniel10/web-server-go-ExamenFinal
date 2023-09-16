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
	LastName string `json:"last_name"`
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

type DentistHandler struct {
	service *dentist.Service
}

func NewDentistHandler(service *dentist.Service) *DentistHandler {
	return &DentistHandler{service: service}
}

func (dh *DentistHandler) GetAll(ctx *gin.Context) {
	data, err := dh.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	if len(data) == 0 {
		ctx.JSON(http.StatusOK, data)
		return
	}

	var body []DentistResponse
	for index, item := range data {
		body[index] = DentistResponse{
			Id:       item.ID,
			LastName: item.Lastname,
			Name:     item.Name,
			License:  item.License,
		}
	}

	ctx.JSON(http.StatusOK, body)
}

func (dh *DentistHandler) GetById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, ResponseError{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ResponseError{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "id param must be a number greater than 0",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	data, err := dh.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, ResponseError{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   "dentist not found",
				Path:      ctx.Request.URL.Path,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	body := DentistResponse{
		Id:       data.ID,
		LastName: data.Lastname,
		Name:     data.Name,
		License:  data.License,
	}

	ctx.JSON(http.StatusOK, body)
}

func (dh *DentistHandler) GetByLicense(ctx *gin.Context) {
	licenseQuery := ctx.Query("license")
	if licenseQuery == "" {
		ctx.JSON(http.StatusBadRequest, ResponseError{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "value of 'license' query param is required",
			Path:      ctx.Request.URL.Path,
		})
		return
	}

	data, err := dh.service.GetByLicense(licenseQuery)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, ResponseError{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    http.StatusNotFound,
				Message:   "dentist not found",
				Path:      ctx.Request.URL.Path,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

}

func (dh *DentistHandler) Create(ctx *gin.Context) {
	bodyToBind := DentistPost{}
	err := ctx.ShouldBindJSON(&bodyToBind)
	if err != nil {

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("'%s' field is: %s",
				extractJSONTag(err.Field(), bodyToBind), err.Tag()))
		}

		ctx.JSON(http.StatusBadRequest, ResponseError{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    http.StatusBadRequest,
			Message:   "invalid body",
			Path:      ctx.Request.URL.Path,
			Errors:    errs,
		})
		return
	}

	dentistToCreate := dentist.Dentist{
		Lastname: bodyToBind.LastName,
		Name:     bodyToBind.Name,
		License:  bodyToBind.License,
	}

	data, err := dh.service.Create(dentistToCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}

	body := DentistResponse{
		Id:       data.ID,
		LastName: data.Lastname,
		Name:     data.Name,
		License:  data.License,
	}

	ctx.JSON(http.StatusCreated, body)
}

func (dh *DentistHandler) Update(ctx *gin.Context) {
}

func (dh *DentistHandler) Patch(ctx *gin.Context) {
}

func (dh *DentistHandler) Delete(ctx *gin.Context) {
}
