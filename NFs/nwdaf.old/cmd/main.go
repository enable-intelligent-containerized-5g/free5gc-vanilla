package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/asaskevich/govalidator"
	"github.com/urfave/cli"

	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/util"
	nwdaf_service "github.com/free5gc/nwdaf/pkg/service"
	"github.com/free5gc/util/version"
)

var NWDAF = &nwdaf_service.NWDAF{}

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.AppLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	app := cli.NewApp()
	app.Name = "nwdaf"
	app.Usage = "5G Unified Data Repository (NWDAF)"
	app.Action = action
	app.Flags = NWDAF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("NWDAF Run error: %v\n", err)
	}
}

func action(c *cli.Context) error {
	if err := initLogFile(c.String("log"), c.String("log5gc")); err != nil {
		logger.AppLog.Errorf("%+v", err)
		return err
	}

	if err := NWDAF.Initialize(c); err != nil {
		switch errType := err.(type) {
		case govalidator.Errors:
			validErrs := err.(govalidator.Errors).Errors()
			for _, validErr := range validErrs {
				logger.CfgLog.Errorf("%+v", validErr)
			}
		default:
			logger.CfgLog.Errorf("%+v", errType)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return fmt.Errorf("Failed to initialize !!")
	}

	logger.AppLog.Infoln(c.App.Name)
	logger.AppLog.Infoln("NWDAF version: ", version.GetVersion())

	NWDAF.Start()

	return nil
}

func initLogFile(logNfPath, log5gcPath string) error {
	NWDAF.KeyLogPath = util.NwdafDefaultKeyLogPath

	if err := logger.LogFileHook(logNfPath, log5gcPath); err != nil {
		return err
	}

	if logNfPath != "" {
		nfDir, _ := filepath.Split(logNfPath)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
			return err
		}
		_, name := filepath.Split(util.NwdafDefaultKeyLogPath)
		NWDAF.KeyLogPath = filepath.Join(tmpDir, name)
	}

	return nil
}
