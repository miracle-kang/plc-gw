package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/c2h5oh/datasize"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/leemcloughlin/logfile"
	"github.com/miracle-kang/plc-gw/api"
	"github.com/miracle-kang/plc-gw/config"
	"github.com/miracle-kang/plc-gw/internal/app"
	"github.com/miracle-kang/plc-gw/internal/pkg"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go startup()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	if len(os.Args) > 1 && pkg.Contains(service.ControlAction, os.Args[1]) {
		svcControl(os.Args[1])
		return
	}
	if service.Interactive() {
		flag.Usage = Usage
		startup()
	} else {
		gin.SetMode(gin.ReleaseMode)
		svcStartup()
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -h, --help      Print current help page")

	flag.PrintDefaults()

	// "start", "stop", "restart", "install", "uninstall"
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Service Commands:")
	fmt.Fprintln(os.Stderr, "  start           Startup the service")
	fmt.Fprintln(os.Stderr, "  stop            Stop the service")
	fmt.Fprintln(os.Stderr, "  restart         Restart the service")
	fmt.Fprintln(os.Stderr, "  install         Install as service")
	fmt.Fprintln(os.Stderr, "  uninstall       Uninstall the service")
}

func newSvc() service.Service {
	var option service.KeyValue
	if runtime.GOOS == "windows" {
		option = service.KeyValue{"DelayedAutoStart": true}
	} else {
		option = service.KeyValue{}
	}
	svcConfig := &service.Config{
		Name:        "plc-gw",
		DisplayName: "PLC Gateway Service",
		Description: "PLC Gateway Service",
		Option:      option,
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func svcControl(action string) {
	s := newSvc()
	err := service.Control(s, action)
	if err != nil {
		log.Fatal(err)
	}
}

func svcStartup() {
	s := newSvc()
	l, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		l.Error(err)
	}
}

func startup() {
	config := config.LoadConfig()

	// Initialize logger
	initLogger(config.Logging)

	if err := app.LoadPLCs(config.PLC); err != nil {
		log.Fatalf("load plcs failed: %v\n", err)
	}

	log.Println("PLC Gateway starting...")
	if err := api.Startup(config.Server.Port); err != nil {
		log.Fatalf("Error startup web server: %v\n", err)
	}
}

func initLogger(logConfig config.LoggingConfig) {
	if logConfig.Path != "" {
		if err := os.MkdirAll(logConfig.Path, os.ModePerm); err != nil {
			log.Fatalf("mkdir %s failed: %v\n", logConfig.Path, err)
		}
		logFileName := logConfig.Path + "/plc-gw.log"

		var s datasize.ByteSize
		s.UnmarshalText([]byte(logConfig.MaxSize))
		log.Printf("Initialize logger to file %s, max size %.2fM, max file %d\n", logFileName, s.MBytes(), logConfig.MaxFile)

		logFile, err := logfile.New(&logfile.LogFile{
			FileName:     logFileName,
			MaxSize:      int64(s.Bytes()),
			FileMode:     logfile.Defaults.FileMode,
			OldVersions:  logConfig.MaxFile,
			FlushSeconds: logfile.Defaults.FlushSeconds,
			CheckSeconds: logfile.Defaults.CheckSeconds,
			Flags:        logfile.FileOnly,
		})
		if err != nil {
			log.Fatalf("Error creating log file: %v", err)
		}
		log.SetOutput(logFile)
	}
}
