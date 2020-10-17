package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config ...
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Init()
}

type viperConfig struct {
}

// NewViperConfig ...
func NewViperConfig() Config {
	v := &viperConfig{}
	v.Init()
	return v
}

func (v *viperConfig) Init() {
	viper.SetEnvPrefix(`go-clean`)
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(`.`, `_`)
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType(`json`)
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

}

func (v *viperConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (v *viperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

func (v *viperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}
