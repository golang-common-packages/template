package config

import (
	"github.com/golang-common-packages/template/common/service/cachestore"
	"github.com/golang-common-packages/template/common/service/datastore"
	"github.com/golang-common-packages/template/common/service/monitor"
	"github.com/golang-common-packages/template/common/util/condition"
	"github.com/golang-common-packages/template/common/util/hash"
	"github.com/golang-common-packages/template/common/util/otp"
	"github.com/golang-common-packages/template/model"

	"github.com/golang-common-packages/cloud-storage"
	"github.com/golang-common-packages/echo-jwt-middleware"
	"github.com/golang-common-packages/email"
)

// Environment model for variable environment
type Environment struct {
	Config    *model.Root
	Database  datastore.Datastore
	Cache     cachestore.Cachestore
	Storage   cloudStorage.Filestore
	Email     email.IMailClient
	Monitor   monitor.MonitorStore
	JWT       jwtMiddleware.Assertion
	Condition condition.Storage
	Hash      hash.Storage
	OTP       otp.Storage
}
