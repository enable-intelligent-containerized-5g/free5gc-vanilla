/*
 * NWDAF Configuration Factory
 */

package factory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/nwdaf/internal/logger"
)

var NwdafConfig Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := os.ReadFile(f); err != nil {
		return err
	} else {
		NwdafConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NwdafConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := NwdafConfig.GetVersion()

	if currentVersion != NWDAF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("Config version is [%s], but expected is [%s]",
			currentVersion, NWDAF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
