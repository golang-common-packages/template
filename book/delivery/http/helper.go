package bookHttpDelivery

import (
	"net/http"

	"github.com/golang-common-packages/template/constant"
	"github.com/sirupsen/logrus"
)

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case constant.ErrInternalServerError:
		return http.StatusInternalServerError
	case constant.ErrNotFound:
		return http.StatusNotFound
	case constant.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
