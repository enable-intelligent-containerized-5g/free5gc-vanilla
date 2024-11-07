package util

import (
	"time"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/internal/sbi/producer"
)

func BackgroundTask() {
	for {
		// logger.InitLog.Info("Start Remove Inactive Nfs")
		err := producer.RemoveInactiveNfs()
		if err == nil {
			logger.InitLog.Info("Remove Inactive Nfs successful")
		}
		time.Sleep(60 * time.Second)
	}
}