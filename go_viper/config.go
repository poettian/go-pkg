package go_viper

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfig() {
	configPath, _  := os.Getwd()
	viper.SetConfigName("app_dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath + "/go_viper")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
