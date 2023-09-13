package handler

import (
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
	"github.com/gin-gonic/gin"
)

type DentistResponse struct {
	Id       int    `json:"id"`
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
	service dentist.Service
}

func NewDentistHandler(service dentist.Service) *DentistHandler {
	return &DentistHandler{service: service}
}

func (dh *DentistHandler) Create(ctx *gin.Context) {
}

func (dh *DentistHandler) GetAll(ctx *gin.Context) {
}

func (dh *DentistHandler) GetById(ctx *gin.Context) {
}

func (dh *DentistHandler) GetByLicense(ctx *gin.Context) {
}

func (dh *DentistHandler) Update(ctx *gin.Context) {
}

func (dh *DentistHandler) Patch(ctx *gin.Context) {
}

func (dh *DentistHandler) Delete(ctx *gin.Context) {
}
