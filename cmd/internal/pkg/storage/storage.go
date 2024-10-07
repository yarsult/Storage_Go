package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	kind string
	st   string
	d    string
}

type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	// defer logger.Sync()
	// logger.Info("Created new storage")
	return Storage{inner: make(map[string]Value),
		logger: logger}, nil
}
func (r Storage) Set(key, val string) {
	var val1 Value
	if _, err := strconv.Atoi(val); err == nil {
		val1 = Value{kind: "D", st: "", d: val}
	} else {
		val1 = Value{kind: "S", st: val, d: ""}
	}
	r.inner[key] = val1
	// r.logger.Info("key has been set")
	// r.logger.Sync()
}
func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		return nil
	}
	// r.logger.Info("val got")
	// r.logger.Sync()
	if res.kind == "S" {
		return &(res).st
	} else {
		return &(res).d
	}
}

func (r Storage) GetKind(key string) string {
	res := r.inner[key]
	return res.kind
}
