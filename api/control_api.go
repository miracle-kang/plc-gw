package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/miracle-kang/plc-gw/internal/app"
)

// ReadTags godoc
// @Summary Read tag value by name list
// @Description Read tag value by name list
// @ID control-read-tag
// @Tags control
// @Accept  json
// @Produce json
// @Param command body ReadTagCommand true "Read Tags"
// @Success	200 {object} BaseResponse{Data=[]TagDto}
// @Router 	/control/readTagsValueByNameList [POST]
func ReadTags(c *gin.Context) {
	var command ReadTagCommand
	if err := c.Bind(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(command)
	log.Println("Control >>> Read tags request: ", string(data))

	gateway := app.GetGateway(command.SN)
	if gateway == nil {
		c.JSON(http.StatusNotFound, Failed("查询数据失败，指定网关不存在"))
		return
	}
	if !gateway.Online {
		c.JSON(http.StatusBadRequest, Failed("查询数据失败，可能因网关不在线导致上报失败，检查网关是否断电或断网！"))
		return
	}
	plc := gateway.PLCs[command.DeviceName]
	if plc == nil {
		c.JSON(http.StatusNotFound, Failed("查询数据失败，指定PLC不存在"))
		return
	}

	tags := make([]*TagDto, len(command.Data))
	for i, tagName := range command.Data {
		var val interface{} = 0
		if plc.ConnState {
			val = plc.RTags[tagName]
		}
		tags[i] = &TagDto{
			Name:  tagName,
			Value: val,
		}
	}
	var res BaseResponse
	if plc.ConnState && len(plc.RTags) > 0 {
		res = SuccessWithData("查询数据成功", tags)
	} else if len(plc.RTags) == 0 {
		res = FailedWithData("数据区为null，查询数据失败", tags)
	} else {
		// PLC Not connected
		res = FailedWithData("查询数据失败", tags)
	}
	c.JSON(http.StatusOK, res)

	data, _ = json.Marshal(res)
	log.Println("Control >>> Read tags response:", string(data))
}

// WriteTag godoc
// @Summary write tag value by name
// @Description write tag value by name
// @ID control-write-tag
// @Tags control
// @Accept  json
// @Produce json
// @Param command body WriteTagCommand true "Write Tag"
// @Success	200 {object} BaseResponse
// @Router 	/control/writeTagValueByName [POST]
func WriteTag(c *gin.Context) {
	var command WriteTagCommand
	if err := c.Bind(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(command)
	log.Println("Control >>> Write tag request: ", string(data))

	gateway := app.GetGateway(command.SN)
	if gateway == nil {
		c.JSON(http.StatusNotFound, Failed("写入数据失败，指定网关不存在"))
		return
	}
	if !gateway.Online {
		c.JSON(http.StatusBadRequest, Failed("写入数据失败，可能因网关不在线，检查网关是否断电或断网！"))
		return
	}
	plc := gateway.PLCs[command.DeviceName]
	if plc == nil {
		c.JSON(http.StatusNotFound, Failed("写入数据失败，指定PLC不存在"))
		return
	}

	var res BaseResponse
	if plc.ConnState {
		plc.WriteTag(command.WTagName, command.WTagValue)
		res = Success("写入数据成功")
	} else {
		// PLC Not connected
		res = Failed("写入数据失败")
	}
	c.JSON(http.StatusOK, res)

	data, _ = json.Marshal(res)
	log.Println("Control >>> Write tag response:", string(data))
}

// ClearGateway godoc
// @Summary clear gateway or plcs tags
// @Description clear gateway or plcs tags
// @ID control-clear-gateway
// @Tags control
// @Accept  json
// @Produce json
// @Param command body ClearGatewayCommand true "Clear Gateway"
// @Success	200 {object} BaseResponse
// @Router 	/control/clearGateway [POST]
func ClearGateway(c *gin.Context) {
	var command ClearGatewayCommand
	if err := c.Bind(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(command)
	log.Println("Control >>> Clear gateway request: ", string(data))

	gateway := app.GetGateway(command.SN)
	if gateway == nil {
		c.JSON(http.StatusNotFound, Failed("清空数据失败，指定SN序列号网关不存在"))
		return
	}
	// Check any plc exists
	exists := false
	for _, v := range command.DeviceNames {
		if gateway.PLCs[v] != nil {
			exists = true
		}
	}
	if !exists && len(command.DeviceNames) > 0 {
		c.JSON(http.StatusNotFound, Failed("清空数据失败，指定SN序列号网关下挂的PLC设备不存在"))
		return
	}
	if len(command.DeviceNames) == 0 {
		for _, p := range gateway.PLCs {
			log.Printf("Control >>> Clearing PLC [%s] Tags\n", p.Name)
			p.ClearTags()
		}
	} else {
		for _, v := range command.DeviceNames {
			plc := gateway.PLCs[v]
			if plc == nil {
				continue
			}
			log.Printf("Control >>> Clearing PLC [%s] Tags\n", plc.Name)
			plc.ClearTags()
		}
	}
	res := Success("清空数据成功")
	c.JSON(http.StatusOK, res)

	data, _ = json.Marshal(res)
	log.Println("Control >>> Clear gateway response:", string(data))
}

// ClearGatewayFile godoc
// @Summary clear gateway or plcs files
// @Description clear gateway or plcs files
// @ID control-clear-gateway-file
// @Tags control
// @Accept  json
// @Produce json
// @Param command body ClearGatewayCommand true "Clear Gateway"
// @Success	200 {object} BaseResponse
// @Router 	/control/clearGatewayFile [POST]
func ClearGatewayFile(c *gin.Context) {
	var command ClearGatewayCommand
	if err := c.Bind(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(command)
	log.Println("Control >>> Clear gateway file request: ", string(data))

	gateway := app.GetGateway(command.SN)
	if gateway == nil {
		c.JSON(http.StatusNotFound, Failed("清空数据失败，指定SN序列号网关不存在"))
		return
	}
	// Check any plc exists
	exists := false
	for _, v := range command.DeviceNames {
		if gateway.PLCs[v] != nil {
			exists = true
		}
	}
	if !exists && len(command.DeviceNames) > 0 {
		c.JSON(http.StatusNotFound, Failed("清空数据失败，指定SN序列号网关下挂的PLC设备不存在"))
		return
	}

	deleted := false
	if len(command.DeviceNames) > 0 {
		for _, name := range command.DeviceNames {
			del := gateway.DeletePLC(name)
			if del {
				deleted = true
			}
			log.Printf("Control >>> Deleted PLC [%s] file, result: %v\n", name, del)
		}
	} else {
		deleted = app.DeleteGateway(command.SN)
	}
	var res BaseResponse
	if deleted {
		res = Success("清除数据文件成功")
	} else {
		res = Failed("清除数据文件失败")
	}
	c.JSON(http.StatusOK, res)

	data, _ = json.Marshal(res)
	log.Println("Control >>> Clear gateway file response:", string(data))

	sort.Ints([]int{1, 2, 3})
}
