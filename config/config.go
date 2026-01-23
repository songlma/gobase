package config

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/songlma/gobase/logger"
	"github.com/spf13/viper"
	//"path/filepath"
)

func Init(ctx context.Context, configFiles []string, watch bool) {
	if len(configFiles) == 0 {
		log.Fatalf("config dir is empty")
	}

	for _, file := range configFiles {
		if file == "" {
			log.Fatalf("config dir is empty")
		}
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		if err != nil {
			log.Fatalf("config fetch fail:%v \n", err)
		}
		if watch {
			viper.WatchConfig()
			viper.OnConfigChange(func(e fsnotify.Event) {
				logger.Infof(ctx, "Config file changed %s ", e.String())
				for _, file := range configFiles {
					if file == "" {
						log.Fatalf("config dir is empty")
					}
					viper.SetConfigFile(file)
					err := viper.MergeInConfig()
					if err != nil {
						log.Fatalf("config fetch fail:%v \n", err)
					}
				}
			})
		}

	}

}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func PutString(key string, value string) {
	viper.Set(key, value)
}
func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func PutInt64(key string, value int64) {
	viper.Set(key, value)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func PutInt(key string, value int) {
	viper.Set(key, value)
}
func GetBool(key string) bool {
	return viper.GetBool(key)
}
func UnmarshalKey(key string, rawVal interface{}) error {
	return viper.UnmarshalKey(key, rawVal)
}
func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}
