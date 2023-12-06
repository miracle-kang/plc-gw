package app

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/miracle-kang/plc-gw/config"
	"github.com/miracle-kang/plc-gw/internal/repo"
)

type Gateway struct {
	SN         string
	IP         string
	Online     bool
	PLCs       map[string]*PLC
	LastReport *time.Time
}

type PLC struct {
	Gateway    *Gateway
	Name       string
	ConnState  bool
	RTags      map[string]interface{}
	WTags      map[string]interface{}
	LastReport *time.Time
	LastWrite  *time.Time
}

// SN->Gateway
var gatewayMp = map[string]*Gateway{}

// Get Gateway by SN, if not exist, create a new one
func AcquireGateway(sn, ip string) *Gateway {
	if gatewayMp[sn] == nil {
		gatewayMp[sn] = &Gateway{
			SN:     sn,
			IP:     ip,
			Online: false,
			PLCs:   map[string]*PLC{},
		}
	}
	// Update gateway ip
	if ip != "" && gatewayMp[sn].IP != ip {
		gatewayMp[sn].IP = ip
	}
	return gatewayMp[sn]
}

func ListGateways() []*Gateway {
	res := make([]*Gateway, 0, len(gatewayMp))
	for _, v := range gatewayMp {
		res = append(res, v)
	}
	return res
}

func GetGateway(sn string) *Gateway {
	return gatewayMp[sn]
}

func DeleteGateway(sn string) bool {
	gateway := gatewayMp[sn]
	res := false
	for name := range gateway.PLCs {
		if gateway.DeletePLC(name) {
			res = true
		}
	}
	res = res || len(gatewayMp) == 0
	delete(gatewayMp, sn)

	log.Printf("Gateway %s(%s)\n deleted", gateway.SN, gateway.IP)
	return res
}

func LoadPLCs(cfg config.PLCConfig) error {
	plcs, err := repo.QueryPLCs()
	if err != nil {
		return err
	}
	for _, plc := range plcs {
		gateway := AcquireGateway(plc.SN, "offline")
		gateway.PLCs[plc.Name] = &PLC{
			Gateway:    gateway,
			Name:       plc.Name,
			RTags:      plc.RTags,
			WTags:      plc.WTags,
			LastReport: plc.LastReport,
			LastWrite:  plc.LastWrite,
		}
	}
	log.Println("Loadded PLCs:", len(plcs))

	// Initialize gateway checker
	return initGatewayChecker(cfg)
}

func initGatewayChecker(cfg config.PLCConfig) error {
	if cfg.CheckInterval < 1 {
		return errors.New("Gateway check interval must greater or equal than 1")
	}
	if cfg.TimeoutSeconds < 5 {
		return errors.New("Gateway check timeout seconds must greater or equal than 5")
	}
	interval := time.Duration(int64(cfg.CheckInterval)) * time.Second
	timeout := time.Duration(int64(cfg.TimeoutSeconds)) * time.Second
	ticker := time.NewTicker(interval)
	go func() {
		for {
			now := <-ticker.C
			for _, g := range gatewayMp {
				if g.LastReport == nil {
					g.Online = false
					continue
				}
				if now.Sub(*g.LastReport) > timeout {
					if g.Online {
						log.Printf("Gateway %s(%s) is offline, last report: %v\n", g.SN, g.IP, g.LastReport)
					}
					g.Online = false
				}
			}
		}
	}()
	log.Printf("Initialized gateway checker every %.0f seconds, check timeout %.0f seconds\n",
		interval.Seconds(), timeout.Seconds())
	return nil
}

// Get PLC by Name, if not exist, create a new one
func (g *Gateway) AcquirePLC(name string) *PLC {
	if g.PLCs[name] == nil {
		g.PLCs[name] = &PLC{
			Gateway: g,
			Name:    name,
			RTags:   map[string]interface{}{},
			WTags:   map[string]interface{}{},
		}
	}
	return g.PLCs[name]
}

func (g *Gateway) DeletePLC(name string) bool {
	if g.PLCs[name] == nil {
		return false
	}

	err := repo.DeletePLC(g.SN, name)
	if err != nil {
		log.Println("Failed to delete plc file", err)
		return false
	}
	delete(g.PLCs, name)
	log.Printf("PLC %s deleted\n", name)
	return true
}

func (p *PLC) Disconnect() {
	p.ConnState = false

	t := time.Now()
	p.LastReport = &t

	repo.SavePLC(&repo.PLCDO{
		SN:         p.Gateway.SN,
		Name:       p.Name,
		ConnState:  p.ConnState,
		RTags:      p.RTags,
		WTags:      p.WTags,
		LastReport: p.LastReport,
		LastWrite:  p.LastWrite,
	})
}

// Report PLC Tags, overwrite the old tags
// Return the write tags
func (p *PLC) ReportTags(tags map[string]interface{}) map[string]interface{} {
	for k, v := range tags {
		p.RTags[k] = v
	}
	p.ConnState = true
	t := time.Now()
	p.LastReport = &t

	if !p.Gateway.Online {
		log.Printf("Gateway %s(%s) is upline", p.Gateway.SN, p.Gateway.IP)
	}
	p.Gateway.Online = true
	p.Gateway.LastReport = &t

	repo.SavePLC(&repo.PLCDO{
		SN:         p.Gateway.SN,
		Name:       p.Name,
		ConnState:  p.ConnState,
		RTags:      p.RTags,
		WTags:      p.WTags,
		LastReport: p.LastReport,
		LastWrite:  p.LastWrite,
	})

	// Return the write tags
	wTags := make(map[string]interface{}, len(p.WTags)+1)
	wTags["PlcGatewayServerConState"] = rand.Intn(100)
	for k, v := range p.WTags {
		wTags[k] = v
	}
	return wTags
}

// Write PLC Tags, overwrite the old tags
func (p *PLC) WriteTags(tags map[string]uint32) {
	for k, v := range tags {
		p.WTags[k] = v
	}
	t := time.Now()
	p.LastWrite = &t

	repo.SavePLC(&repo.PLCDO{
		SN:         p.Gateway.SN,
		Name:       p.Name,
		ConnState:  p.ConnState,
		RTags:      p.RTags,
		WTags:      p.WTags,
		LastReport: p.LastReport,
		LastWrite:  p.LastWrite,
	})
}

// Write PLC Tag
func (p *PLC) WriteTag(tag string, value interface{}) {
	p.WTags[tag] = value
	t := time.Now()
	p.LastWrite = &t

	repo.SavePLC(&repo.PLCDO{
		SN:         p.Gateway.SN,
		Name:       p.Name,
		ConnState:  p.ConnState,
		RTags:      p.RTags,
		WTags:      p.WTags,
		LastReport: p.LastReport,
		LastWrite:  p.LastWrite,
	})
}

// Clear RTags and WTags
func (p *PLC) ClearTags() {
	for k := range p.RTags {
		delete(p.RTags, k)
	}
	for k := range p.WTags {
		delete(p.WTags, k)
	}

	repo.SavePLC(&repo.PLCDO{
		SN:         p.Gateway.SN,
		Name:       p.Name,
		ConnState:  p.ConnState,
		RTags:      p.RTags,
		WTags:      p.WTags,
		LastReport: p.LastReport,
		LastWrite:  p.LastWrite,
	})
}
