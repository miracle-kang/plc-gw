package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miracle-kang/plc-gw/internal/app"
)

// ReportData godoc
// @Summary Gateway report data
// @Description Gateway report data
// @ID report-data
// @Tags report
// @Accept  json
// @Produce json
// @Param command body ReportDataCommand true "Report Data"
// @Success	200 {object} ReportDataResponse
// @Router 	/report [POST]
func ReportData(c *gin.Context) {
	var command ReportDataCommand
	if err := c.Bind(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(command)
	log.Println("Report  >>> Received data:", string(data))

	// Report RTag values
	gateway := app.AcquireGateway(command.SN, c.ClientIP())
	plc := gateway.AcquirePLC(command.DeviceName)

	var retTags map[string]interface{}
	if command.GetDataError != nil {
		plc.Disconnect()
		retTags = map[string]interface{}{"PlcGatewayServerConState": rand.Intn(100)}
	} else {
		retTags = plc.ReportTags(command.Data)
	}

	// Response WTag values
	res := &ReportDataResponse{
		DeviceName: plc.Name,
		Data:       retTags,
	}
	c.JSON(http.StatusOK, res)

	data, _ = json.Marshal(res)
	log.Println("Report  >>> Response data:", string(data))
}
