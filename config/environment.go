package config

import (
	"github.com/golang-common-packages/template/common/service/datastore"
	"github.com/golang-common-packages/template/common/service/monitor"
	"github.com/golang-common-packages/template/common/util/condition"
	"github.com/golang-common-packages/template/common/util/otp"
	"github.com/golang-common-packages/template/model"

	"github.com/golang-common-packages/caching"
	"github.com/golang-common-packages/cloud-storage"
	"github.com/golang-common-packages/echo-jwt-middleware"
	"github.com/golang-common-packages/email"
	"github.com/golang-common-packages/hash"
)

// Environment model for variable environment
type Environment struct {
	Config    *model.Root
	Database  datastore.Datastore
	Cache     caching.ICaching
	Storage   cloudStorage.Filestore
	Email     email.IMailClient
	Monitor   monitor.MonitorStore
	JWT       jwtMiddleware.Assertion
	Condition condition.Storage
	Hash      hash.IHash
	OTP       otp.Storage
}
