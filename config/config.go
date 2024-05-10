package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

type StartConfigStruct struct {
	Aws struct {
		AwsRegion string `mapstructure:"aws_region"`

		// 登录key, 没有的话则走机器的
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`

		Profile string `mapstructure:"profile"`

		// 数据表前缀
		DynamoTablePrefix string `mapstructure:"dynamo_table_prefix"`

		// 内部开发时可以填dynamo local的地址, 正式时填interface endpoint
		// 不填则是常规的公网访问dynamo
		DynamoEndpoint string `mapstructure:"dynamo_endpoint"`

		// 尝试创建表
		CreateTableAnyway bool `mapstructure:"create_table_anyway"`
	} `mapstructure:"aws"`

	Log struct {
		LogrusLevel string `mapstructure:"logrus_level"`
	} `mapstructure:"log"`

	SwitchController struct {
		IsDebugMode   bool      `mapstructure:"is_debug_mode"`
		SetServerTime time.Time `mapstructure:"server_time"`
	} `mapstructure:"switch_controller"`
}

var (
	StartConfig = new(StartConfigStruct)

	configReadOnce sync.Once
	configReadErr  error

	HolidayConfig = new(HolidayData)
)

func init() {
	StartConfig.Log.LogrusLevel = "trace"
}

func (d *StartConfigStruct) OnLoaded() error {
	return nil
}

func ReadStartUpConfig() error {
	configReadOnce.Do(func() {
		configReadErr = doReadStartUpConfig()
	})
	return configReadErr
}

func doReadStartUpConfig() error {

	file := os.Getenv("STARSERVER_CONFIG")
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		viper.SetConfigName("server")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "读fight启动配置文件出错")
	}

	err := viper.Unmarshal(StartConfig)
	if err != nil {
		return errors.Wrap(err, "解析fight启动配置出错")
	}

	if err := StartConfig.OnLoaded(); err != nil {
		return errors.Wrap(err, "解析server.yaml启动配置出错")
	}

	return nil
}

func LoadHolidayConfig() (err error) {

	HolidayConfig, err = LoadHolidaysJson()

	return err
}
