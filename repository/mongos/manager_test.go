package mongos

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"testing"
)

func TestFuncName(t *testing.T) {
	logrus.Infof("Geeter func name:%s", runtime.FuncForPC(reflect.ValueOf(Greeter).Pointer()).Name())
}

func TestFuncName2(t *testing.T) {
	logrus.Infof("Geeter func name:%s", runtime.FuncForPC(reflect.ValueOf(&MongoManager{}).MethodByName("SaveDocuments").Pointer()).Name())
}
