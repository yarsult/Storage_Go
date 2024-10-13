package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	kind Kind
	st   string
	d    string
}

type Kind string

const (
	KindString Kind = "S"
	KindInt    Kind = "D"
)

type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	defer logger.Sync()
	logger.Info("Created new storage")
	return Storage{inner: make(map[string]Value),
		logger: logger}, nil
}

func (r Storage) Set(key, val string) {
	var val1 Value
	if _, err := strconv.Atoi(val); err == nil {
		val1 = Value{kind: KindInt, d: val}
	} else {
		val1 = Value{kind: KindString, st: val}
	}
	r.inner[key] = val1
	r.logger.Info("key has been set")
}

func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		return nil
	}
	r.logger.Info("val got")
	if res.kind == KindString {
		return &(res).st
	}
	return &(res).d
}

func (r Storage) GetKind(key string) string {
	res := r.inner[key]
	return string(res.kind)
}
