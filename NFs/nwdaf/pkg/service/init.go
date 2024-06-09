package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/nwdaf/internal/context"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type NWDAF struct {
	KeyLogPath string
}

type (
	// Commands information.
	Commands struct {
		config string
	}
)

var commands Commands

var cliCmd = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
	cli.StringFlag{
		Name:  "log, l",
		Usage: "Output NF log to `FILE`",
	},
	cli.StringFlag{
		Name:  "log5gc, lc",
		Usage: "Output free5gc log to `FILE`",
	},
}

func (*NWDAF) GetCliCmd() (flags []cli.Flag) {
	return cliCmd
}

func (udr *NWDAF) Initialize(c *cli.Context) error {
	commands = Commands{
		config: c.String("config"),
	}

	if commands.config != "" {
		if err := factory.InitConfigFactory(commands.config); err != nil {
			return err
		}
	} else {
		if err := factory.InitConfigFactory(util.NwdafDefaultConfigPath); err != nil {
			return err
		}
	}

	udr.SetLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	if _, err := factory.NwdafConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (udr *NWDAF) SetLogLevel() {
	if factory.NwdafConfig.Logger == nil {
		logger.InitLog.Warnln("UDR config without log level setting!!!")
		return
	}

	if factory.NwdafConfig.Logger.UDR != nil {
		if factory.NwdafConfig.Logger.UDR.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.UDR.DebugLevel); err != nil {
				logger.InitLog.Warnf("UDR Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.UDR.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				logger.InitLog.Infof("UDR Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			logger.InitLog.Infoln("UDR Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NwdafConfig.Logger.UDR.ReportCaller)
	}
}

func (udr *NWDAF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range udr.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (udr *NWDAF) Start() {
	// get config file info
	config := factory.NwdafConfig

	logger.InitLog.Infof("NWDAF Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	logger.InitLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)
	// router.Use(cors.New(cors.Config{
	// 	AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders: []string{
	// 		"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host",
	// 		"Token", "X-Requested-With",
	// 	},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowAllOrigins:  true,
	// 	MaxAge:           86400,
	// }))
	//
	// datarepository.AddService(router)

	pemPath := util.NwdafDefaultPemPath
	keyPath := util.NwdafDefaultKeyPath
	sbi := config.Configuration.Sbi
	if sbi.Tls != nil {
		pemPath = sbi.Tls.Pem
		keyPath = sbi.Tls.Key
	}

	self := context.NWDAF_Self()
	util.InitNwdafContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)
	profile := consumer.BuildNFInstance(self)
	var newNrfUri string
	var err error
	newNrfUri, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, profile.NfInstanceId, profile)
	if err == nil {
		self.NrfUri = newNrfUri
	} else {
		logger.InitLog.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		<-signalChannel
		udr.Terminate()
		os.Exit(0)
	}()

	server, err := httpwrapper.NewHttp2Server(addr, udr.KeyLogPath, router)
	if server == nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		logger.InitLog.Warnf("Initialize HTTP server: %+v", err)
	}

	serverScheme := factory.NwdafConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(pemPath, keyPath)
	}

	if err != nil {
		logger.InitLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (udr *NWDAF) Exec(c *cli.Context) error {
	// UDR.Initialize(cfgPath, c)

	logger.InitLog.Traceln("args:", c.String("udrcfg"))
	args := udr.FilterCli(c)
	logger.InitLog.Traceln("filter: ", args)
	command := exec.Command("./udr", args...)

	if err := udr.Initialize(c); err != nil {
		return err
	}

	var stdout io.ReadCloser
	if readCloser, err := command.StdoutPipe(); err != nil {
		logger.InitLog.Fatalln(err)
	} else {
		stdout = readCloser
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var stderr io.ReadCloser
	if readCloser, err := command.StderrPipe(); err != nil {
		logger.InitLog.Fatalln(err)
	} else {
		stderr = readCloser
	}
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var err error
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		if errormessage := command.Start(); err != nil {
			fmt.Println("command.Start Fails!")
			err = errormessage
		}
		wg.Done()
	}()

	wg.Wait()
	return err
}

func (udr *NWDAF) Terminate() {
	logger.InitLog.Infof("Terminating UDR...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}
	logger.InitLog.Infof("UDR terminated")
}
