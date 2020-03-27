package units

import (
	log "github.com/Code-Fight/golog"
	"github.com/spf13/viper"
)

func ConfigInit() {
	viper.SetDefault("port",2048)
	viper.SetDefault("log","info")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Errorf("Fatal error config file: %s \n", err)
	}

}

//get port by config.json
func GetPort() string {
	return viper.GetString("port")
}

//get log level  by config.json
func GetLog() string {
	return viper.GetString("log")
}
