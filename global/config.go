package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"path"
)

var SensitiveWords []string
var MessageQueenLen int

func initConfig() {
	viper.SetConfigName("chatroom")
	viper.AddConfigPath(path.Join(RootDir, "config"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	SensitiveWords = viper.GetStringSlice("SensitiveWords")
	MessageQueenLen = viper.GetInt("MessageQueenLen")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		SensitiveWords = viper.GetStringSlice("SensitiveWords")
	})
}
