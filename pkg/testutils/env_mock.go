package testutils

import (
	"os"
	"sync"
)

type EnvMock struct {
	mtx   sync.Mutex
	store map[string]string
}

func NewEnvMock() *EnvMock {
	return &EnvMock{
		store: make(map[string]string),
	}
}

func (mock *EnvMock) Set(key string, value string) {
	mock.mtx.Lock()
	defer mock.mtx.Unlock()
	mock.store[key] = value
}

func (mock *EnvMock) Get(key string) string {
	mock.mtx.Lock()
	defer mock.mtx.Unlock()
	return mock.store[key]
}

func (mock *EnvMock) Apply(key string) {
	mock.mtx.Lock()
	defer mock.mtx.Unlock()
	for key, value := range mock.store {
		os.Setenv(key, value)
	}
}

func (mock *EnvMock) Restore(originalEnv *map[string]string) {
	mock.mtx.Lock()
	defer mock.mtx.Unlock()
	os.Clearenv()
	mock.store = *originalEnv
	for key, value := range *originalEnv {
		os.Setenv(key, value)
	}
}
