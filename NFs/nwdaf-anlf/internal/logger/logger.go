package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	log           *logrus.Logger
	AppLog        *logrus.Entry
	InitLog       *logrus.Entry
	CfgLog        *logrus.Entry
	ProducerLog   *logrus.Entry
	AniLog        *logrus.Entry
	NfPer         *logrus.Entry
	HandlerLog    *logrus.Entry
	DataRepoLog   *logrus.Entry
	UtilLog       *logrus.Entry
	HttpLog       *logrus.Entry
	ConsumerLog   *logrus.Entry
	GinLog        *logrus.Entry
	PcmLog        *logrus.Entry
	MlModelTrainingLog *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	AppLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "CFG"})
	ProducerLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Producer"})
	AniLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "AnalyticsInfo"})
	NfPer = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "NfProfileProvition"})
	PcmLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "PCM"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "HDLR"})
	DataRepoLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "DRepo"})
	UtilLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Util"})
	HttpLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "HTTP"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Consumer"})
	GinLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "GIN"})
	MlModelTrainingLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "MMTraining"})
}

func LogFileHook(logNfPath string, log5gcPath string) error {
	if fullPath, err := logger_util.CreateFree5gcLogFile(log5gcPath); err == nil {
		if fullPath != "" {
			free5gcLogHook, hookErr := logger_util.NewFileHook(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
			if hookErr != nil {
				return hookErr
			}
			log.Hooks.Add(free5gcLogHook)
		}
	} else {
		return err
	}

	log.SetLevel(logrus.InfoLevel)

	if fullPath, err := logger_util.CreateNfLogFile(logNfPath, "nwdaf.log"); err == nil {
		selfLogHook, hookErr := logger_util.NewFileHook(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
		if hookErr != nil {
			return hookErr
		}
		log.Hooks.Add(selfLogHook)
	} else {
		return err
	}

	return nil
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(enable bool) {
	log.SetReportCaller(enable)
}
