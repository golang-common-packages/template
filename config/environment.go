package config

import (
	"github.com/golang-microservices/template/common/service/cachestore"
	"github.com/golang-microservices/template/common/service/datastore"
	"github.com/golang-microservices/template/common/service/email"
	//"github.com/golang-microservices/template/common/service/filestore"
	"github.com/golang-microservices/template/common/service/monitor"
	"github.com/golang-microservices/template/common/util/condition"
	"github.com/golang-microservices/template/common/util/hash"
	"github.com/golang-microservices/template/common/util/jwt"
	"github.com/golang-microservices/template/common/util/otp"
	"github.com/golang-microservices/template/model"

	"github.com/golang-microservices/cloud-storage"
)

// Environment stuct for variable environment
type Environment struct {
	Config    *model.Root
	Database  datastore.Datastore
	Cache     cachestore.Cachestore
	Storage   cloudStorage.Filestore
	Email     email.Mailstore
	Monitor   monitor.MonitorStore
	JWT       jwt.Storage
	Condition condition.Storage
	Hash      hash.Storage
	OTP       otp.Storage
}
