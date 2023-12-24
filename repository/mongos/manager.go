package mongos

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"sync"
)

var manager MongoManager

// Interceptor 定义拦截器接口
type Interceptor interface {
	Before()
	After()
}

// MongoInterceptor 用于包装函数并应用拦截器
type MongoInterceptor struct {
	interceptor Interceptor
	function    reflect.Value
}

// MongoManager 管理 Mongo 拦截器的注册和应用
type MongoManager struct {
	interceptors map[interface{}]*MongoInterceptor
	mu           sync.RWMutex

	Client     *mongo.Client
	DataBase   *mongo.Database
	Collection *mongo.Collection
}

// newMongoManager 创建一个 MongoManager 实例
func newMongoManager() *MongoManager {
	return &MongoManager{
		interceptors: make(map[interface{}]*MongoInterceptor),
	}
}

// Register 注册函数和拦截器
func (m *MongoManager) Register(fn interface{}, interceptor Interceptor) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.interceptors[fn] = &MongoInterceptor{
		interceptor: interceptor,
		function:    reflect.ValueOf(m).MethodByName(reflect.ValueOf(fn).String()),
	}
}

// Apply 执行包装的函数，并在前后应用拦截器
func (a *MongoInterceptor) Apply(params []reflect.Value) []reflect.Value {
	a.interceptor.Before()
	defer a.interceptor.After()

	return a.function.Call(params)
}

// Call 执行注册的函数，应用拦截器
func (m *MongoManager) Call(fn interface{}, params ...interface{}) {
	m.mu.RLock()
	interceptor, ok := m.interceptors[fn]
	m.mu.RUnlock()

	if !ok {
		panic("Function not registered with MongoManager")
	}

	args := make([]reflect.Value, len(params))
	for i, p := range params {
		args[i] = reflect.ValueOf(p)
	}

	interceptor.Apply(args)
}

// AuthInterceptor 实现拦截器接口，用于记录日志
type AuthInterceptor struct{}

// Before 在调用前记录日志
func (auth *AuthInterceptor) Before() {
	logrus.Info("Before function call - Logging")
}

// After 在调用后记录日志
func (auth *AuthInterceptor) After() {
	logrus.Info("After function call - Logging")
}

// Greeter 是一个示例函数
func Greeter(name string) error {
	fmt.Printf("Hello, %s!\n", name)
	return nil
}

func Register() *MongoManager {
	m := newMongoManager()

	// 注册 SaveDocuments 函数和 LoggingInterceptor 拦截器

	m.Register(Greeter, &AuthInterceptor{})

	return m
}
