package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miracle-kang/plc-gw/internal/app"
)

// ListGateways godoc
// @Summary List gateways
// @Description List gateways
// @ID manager-list-gateways
// @Tags manager
// @Accept  json
// @Produce json
// @Success	200 {object} BaseResponse{Data=[]GatewayDto}
// @Router 	/manager/gateways [GET]
func ListGateways(c *gin.Context) {
	gateways := app.ListGateways()
	res := make([]*GatewayDto, len(gateways))
	for i, g := range gateways {
		plcs := make([]string, 0, len(g.PLCs))
		for k := range g.PLCs {
			plcs = append(plcs, k)
		}
		res[i] = &GatewayDto{
			SN:     g.SN,
			IP:     g.IP,
			Online: g.Online,
			PLCs:   plcs,
		}
	}
	c.JSON(http.StatusOK, SuccessWithData("查询数据成功", res))
}

// ListGatewayPLCs godoc
// @Summary List gateway PLCs
// @Description List gateway PLCs
// @ID manager-list-gateway-plcs
// @Tags manager
// @Accept  json
// @Produce json
// @Param sn path string true "Gateway SN"
// @Success	200 {object} BaseResponse{Data=[]PLCDto}
// @Router 	/manager/gateways/{sn}/plcs [GET]
func ListGatewayPLCs(c *gin.Context) {
	sn := c.Param("sn")
	gateway := app.GetGateway(sn)
	if gateway == nil {
		c.JSON(http.StatusNotFound, Failed("查询数据失败，网关SN号错误"))
		return
	}
	plcs := make([]*PLCDto, 0, len(gateway.PLCs))
	for _, v := range gateway.PLCs {
		rTags := make([]*TagDto, 0, len(v.RTags))
		for k, v := range v.RTags {
			rTags = append(rTags, &TagDto{k, v})
		}
		wTags := make([]*TagDto, 0, len(v.WTags))
		for k, v := range v.WTags {
			wTags = append(wTags, &TagDto{k, v})
		}
		plcs = append(plcs, &PLCDto{
			SN:         gateway.SN,
			Name:       v.Name,
			ConnState:  v.ConnState,
			RTags:      rTags,
			WTags:      wTags,
			LastReport: v.LastReport,
			LastWrite:  v.LastWrite,
		})
	}

	c.JSON(http.StatusOK, SuccessWithData("查询数据成功", plcs))
}
