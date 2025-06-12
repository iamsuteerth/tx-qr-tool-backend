package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/dto"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/service"
	"github.com/iamsuteerth/tx-qr-tool-backend/utils"
)

type RegistrationController struct {
	service service.RegistrationService
}

func NewRegistrationController(service service.RegistrationService) *RegistrationController {
	return &RegistrationController{service: service}
}

func (rc *RegistrationController) Register(c *gin.Context) {
	requestID := utils.GetRequestID(c)

	var req dto.CreateRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleErrorResponse(c, utils.NewBadRequestError("INVALID_JSON", "Invalid JSON format", err), requestID)
		return
	}

	response, err := rc.service.CreateRegistration(&req)
	if err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	utils.SendCreatedResponse(c, "Registration created successfully", requestID, response)
}

func (rc *RegistrationController) DownloadCSV(c *gin.Context) {
	requestID := utils.GetRequestID(c)

	csvData, err := rc.service.GenerateCSV()
	if err != nil {
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	filename := fmt.Sprintf("registrations_%s.csv", time.Now().Format("2006-01-02_15-04-05"))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Length", fmt.Sprintf("%d", len(csvData)))

	c.Data(http.StatusOK, "text/csv", csvData)
}
