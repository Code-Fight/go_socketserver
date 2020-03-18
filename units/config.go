package units

import (
	"fmt"
	"github.com/spf13/viper"
)

//get port by config.json
func GetPort() string {
	viper.SetDefault("port",2048)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return viper.GetString("port")
}
