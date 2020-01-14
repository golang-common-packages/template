package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/spf13/viper"

	"github.com/golang-microservices/template/model"
)

// Load function will read config from environment or config file.
func Load(prefix string, fileNames ...string) model.Root {
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Name, type and path of config file
	viper.SetConfigType("yaml")
	viper.SetConfigName("main")
	viper.AddConfigPath("config/")
	viper.AddConfigPath("../../config/")

	for _, fileName := range fileNames {
		viper.SetConfigName(fileName)
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			panic(err)
		}
	}
	var c model.Root
	bindEnvs(c)
	viper.Unmarshal(&c)
	return c
}

// bindEnvs function will bind ymal file to struc model
func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
