package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/miracle-kang/plc-gw/config"
	"github.com/miracle-kang/plc-gw/internal/pkg"
)

type PLCDO struct {
	SN         string
	Name       string
	ConnState  bool
	RTags      map[string]interface{}
	WTags      map[string]interface{}
	LastReport *time.Time
	LastWrite  *time.Time
}

func SavePLC(plc *PLCDO) error {
	config := config.LoadedConfig
	if exists, _ := pkg.ExistsFile(config.PLC.BasePath); !exists {
		os.MkdirAll(config.PLC.BasePath, 0755)
	}
	filename := plc.SN + "-" + plc.Name + ".json"
	data, err := json.Marshal(plc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(config.PLC.BasePath, filename), data, 0644)
}

func DeletePLC(sn, name string) error {
	config := config.LoadedConfig
	if exists, _ := pkg.ExistsFile(config.PLC.BasePath); !exists {
		return errors.New("base path not exists")
	}

	filepath := filepath.Join(config.PLC.BasePath, sn+"-"+name+".json")
	if exists, _ := pkg.ExistsFile(filepath); !exists {
		return fmt.Errorf("PLC file '%s' not exists", filepath)
	}
	err := os.Remove(filepath)
	if err == nil {
		log.Printf("PLC file '%s' deleted", filepath)
	}
	return err
}

// id->plc
func QueryPLCs() ([]*PLCDO, error) {
	config := config.LoadedConfig
	if b, _ := pkg.ExistsFile(config.PLC.BasePath); !b {
		return []*PLCDO{}, nil
	}
	files, err := ioutil.ReadDir(config.PLC.BasePath)
	if err != nil {
		return nil, err
	}
	plcs := make([]*PLCDO, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(filepath.Join(config.PLC.BasePath, file.Name()))
		if err != nil {
			return nil, err
		}
		var plc PLCDO
		err = json.Unmarshal(data, &plc)
		if err != nil {
			return nil, err
		}
		plcs = append(plcs, &plc)
	}
	return plcs, nil
}
