/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 14:26
 * @Description:
 */

package conf

import (
	"blog/pkg/logger"
	"blog/pkg/redis"
	"blog/pkg/storage/elasticsearch"
	"blog/pkg/storage/mysql"
	"blog/pkg/storage/orm"
	"blog/pkg/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig
	Logger        logger.Config
	ORM           orm.Config
	Mysql         mysql.Config
	Redis         redis.Config
	Elasticsearch elasticsearch.Config
}

type AppConfig struct {
	Name            string
	Version         bool
	Mode            string
	Addr            string
	Host            string
	Cdn             string
	DisableDingDing bool
}

var (
	Conf = &Config{}
)

func Init(configPath string) *Config {
	viper.SetConfigType("yml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); utils.IsNotNil(err) {
		panic(err)
	}
	if err := viper.Unmarshal(&Conf); utils.IsNotNil(err) {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&Conf); utils.IsNotNil(err) {
			panic(err)
		}
	})
	return Conf
}
