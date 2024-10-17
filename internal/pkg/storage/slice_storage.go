package storage

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"proj1/internal/pkg/saving"
	"slices"
	"strconv"

	"go.uber.org/zap"
)

type SliceValue struct {
	Kind Kind
	St   []string
	D    []string
}

type SliceStorage struct {
	inner  map[string]SliceValue
	logger *zap.Logger
}

func NewSliceStorage() (SliceStorage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return SliceStorage{}, err
	}
	defer logger.Sync()
	logger.Info("Created new storage for slices")
	return SliceStorage{inner: make(map[string]SliceValue),
		logger: logger}, nil
}

func (s SliceStorage) createIfEmpty(values []string) SliceValue {
	if _, err := strconv.Atoi(values[0]); err == nil {
		return SliceValue{Kind: KindInt}
	}
	return SliceValue{Kind: KindString}
}

func (s SliceStorage) defineKind(key string) []string {
	if s.inner[key].Kind == KindInt {
		return s.inner[key].D
	}
	return s.inner[key].St
}

func (s SliceStorage) addToAppropriate(key string, values []string, cur SliceValue) {
	if cur.Kind == KindInt {
		s.inner[key] = SliceValue{Kind: KindInt, D: values}
	} else {
		s.inner[key] = SliceValue{Kind: KindString, St: values}
	}
}

func (s SliceStorage) LPush(key string, values []string) {
	val, ok := s.inner[key]
	var tmp []string
	tmp = append(tmp, values...)
	slices.Reverse(tmp)
	if !ok {
		s.addToAppropriate(key, tmp, s.createIfEmpty(values))

	} else {
		s.addToAppropriate(key, append(tmp, s.defineKind(key)...), val)
	}
}

func (s SliceStorage) RPush(key string, values []string) {
	val, ok := s.inner[key]
	var tmp []string
	tmp = append(tmp, values...)
	if !ok {
		s.addToAppropriate(key, tmp, s.createIfEmpty(values))

	} else {
		s.addToAppropriate(key, append(s.defineKind(key), tmp...), val)
	}
}

func (s SliceStorage) RAddToSet(key string, values []string) {
	val, ok := s.inner[key]
	var tmp []string
	if !ok {
		tmp = append(tmp, values...)
		s.addToAppropriate(key, tmp, s.createIfEmpty(values))
	} else {
		myValues := s.defineKind(key)
		for _, x := range values {
			if !slices.Contains(myValues, x) {
				tmp = append(tmp, x)
			}
		}
		s.addToAppropriate(key, append(s.defineKind(key), tmp...), val)
	}
}

func (s SliceStorage) LPop(key string, indexes ...int) []string {
	var start int
	end := indexes[0]
	if len(indexes) == 2 {
		start = indexes[0]
		end = indexes[1] + 1
	}
	val, ok := s.inner[key]
	if !ok {
		return []string{}
	}
	res := s.defineKind(key)
	if start < 0 {
		start += len(res)
	}
	if end < 0 {
		end += len(res)
	}
	if end > len(res) {
		if math.Abs((float64(start))) > float64(len(res)) {
			s.logger.Info("uncorrect indexes")
			return []string{}
		}
		s.addToAppropriate(key, res[:start], val)
		return res[start:]
	}
	s.addToAppropriate(key, slices.Concat(res[:start], res[end:]), val)
	return res[start:end]
}

func (s SliceStorage) RPop(key string, indexes ...int) []string {
	var start int
	end := indexes[0]
	val, ok := s.inner[key]
	if !ok {
		return []string{}
	}
	res := s.defineKind(key)
	slices.Reverse(res)
	if len(indexes) == 2 {
		start = indexes[0]
		end = indexes[1] + 1
	}
	if start < 0 {
		start += len(res)
	}
	if end < 0 {
		end += len(res)
	}
	if end > len(res) {
		if math.Abs((float64(start))) > float64(len(res)) {
			s.logger.Info("uncorrect indexes")
			return []string{}
		}
		tmp := res[:start]
		slices.Reverse(tmp)
		s.addToAppropriate(key, tmp, val)
		return res[start:]
	}
	tmp := slices.Concat(res[:start], res[end:])
	slices.Reverse(tmp)
	s.addToAppropriate(key, tmp, val)
	return res[start:end]
}

func (s SliceStorage) LSet(key string, index int, elem string) (string, error) {
	_, ok := s.inner[key]
	if !ok {
		s.logger.Info("no such key")
		return "", errors.New("no such key")
	}
	res := s.defineKind(key)
	if index > len(res)-1 {
		s.logger.Info("out of range")
		return "", errors.New("slice bounds ot of range")
	}
	res[index] = elem
	return "ok", nil
}

func (s SliceStorage) LGet(key string, index int) (string, error) {
	_, ok := s.inner[key]
	if !ok {
		s.logger.Info("no such key")
		return "", errors.New("no such key")
	}
	res := s.defineKind(key)
	if index > len(res)-1 {
		s.logger.Info("out of range")
		return "", errors.New("slice bounds ot of range")
	}
	return res[index], nil
}

func (s SliceStorage) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(s.inner, "", "  ")
	if err != nil {
		s.logger.Error("Failed to marshal SliceStorage to JSON", zap.Error(err))
		return err
	}
	err = saving.WriteAtomic(filename, data)
	if err != nil {
		s.logger.Error("Failed to write JSON to file", zap.Error(err))
		return err
	}

	s.logger.Info("SliceStorage successfully saved to file", zap.String("filename", filename))
	return nil
}

func (s *SliceStorage) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		s.logger.Error("Failed to read file", zap.Error(err))
		return err
	}
	var inner map[string]SliceValue
	if err = json.Unmarshal(data, &inner); err != nil {
		s.logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return err
	}
	s.inner = inner
	s.logger.Info("SliceStorage successfully loaded from file", zap.String("filename", filename))
	return nil
}
