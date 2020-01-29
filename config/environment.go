package config

import (
	"github.com/golang-common-packages/template/common/service/datastore"
	"github.com/golang-common-packages/template/model"

	"github.com/golang-common-packages/caching"
	"github.com/golang-common-packages/cloud-storage"
	"github.com/golang-common-packages/echo-jwt-middleware"
	"github.com/golang-common-packages/email"
	"github.com/golang-common-packages/hash"
	"github.com/golang-common-packages/monitoring"
	"github.com/golang-common-packages/otp"
)

// Environment model for variable environment
type Environment struct {
	Config   *model.Root
	Database datastore.Datastore
	Cache    caching.ICaching
	Storage  cloudStorage.Filestore
	Email    email.IMailClient
	Monitor  monitoring.IMonitoring
	JWT      jwtMiddleware.Assertion
	Hash     hash.IHash
	OTP      otp.IOTP
}
