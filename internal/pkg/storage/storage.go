package storage

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"proj1/internal/pkg/saving"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type SliceValue struct {
	Kind       Kind
	Expires_at int64
	St_sl      []string
	D_sl       []string
	St         string
	D          string
	Mint       map[string]int
	Mstr       map[string]string
}

type SliceStorage struct {
	inner  map[string]SliceValue
	logger *zap.Logger
	mu     sync.RWMutex
}

type Kind string

const (
	KindString   Kind = "S"
	KindInt      Kind = "D"
	KindSliceInt Kind = "SD"
	KindSliceStr Kind = "SS"
	KindMapInt   Kind = "MI"
	KindMapStr   Kind = "MS"
)

func NewSliceStorage() (SliceStorage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return SliceStorage{}, err
	}
	defer logger.Sync()
	logger.Info("Created new storage")
	return SliceStorage{inner: make(map[string]SliceValue),
		logger: logger}, nil
}

func (s *SliceStorage) Set(key, val string) error {
	var val1 SliceValue
	if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
		val1 = SliceValue{Kind: KindString, St: strings.Trim(val, `"`)}
	} else {
		if _, err := strconv.Atoi(val); err != nil {
			return errors.New("uncorrect string")
		}
		val1 = SliceValue{Kind: KindInt, D: val}
	}
	s.inner[key] = val1
	s.logger.Info("key has been set")
	return nil
}

func (s *SliceStorage) Get(key string) (string, bool) {
	res, ok := s.inner[key]
	if !ok {
		return "", false
	}
	s.logger.Info("val got")
	if res.Kind == KindString {
		return res.St, true
	}
	return res.D, true
}

func (s *SliceStorage) GetKind(key string) string {
	res := s.inner[key]
	return string(res.Kind)
}

func (s *SliceStorage) HSet(key string, maps []map[string]string) (int, error) {
	other_types := []Kind{KindInt, KindSliceInt, KindSliceStr, KindString}
	if slices.Contains(other_types, s.inner[key].Kind) {
		s.logger.Info("uncorrect indexes")
		return 0, errors.New("no such key")
	}
	final1 := make(map[string]string)
	final2 := make(map[string]int)
	for _, x := range maps {
		for k, v := range x {
			if res, err := strconv.Atoi(v); err == nil {
				final2[k] = res
				//
			} else {
				final1[k] = v
				//
			}
		}
	}
	if len(final1) > 0 {
		s.inner[key] = SliceValue{Kind: KindMapStr, Mstr: final1}
		return len(final1), nil
	}
	s.inner[key] = SliceValue{Kind: KindMapInt, Mint: final2}
	return len(final2), nil
}

func (s *SliceStorage) HGet(key string, field string) (*string, error) {
	res := s.inner[key]
	res1, ok1 := res.Mint[field]
	res2, ok2 := res.Mstr[field]
	if res.Kind == "" || (!ok1 && !ok2) {
		s.logger.Info("uncorrect key or field")
		return nil, nil
	}
	other_types := []Kind{KindInt, KindSliceInt, KindSliceStr, KindString}
	if slices.Contains(other_types, s.inner[key].Kind) {
		s.logger.Info("uncorrect indexes")
		return nil, errors.New("no such key")
	}
	if ok1 {
		fin := strconv.Itoa(res1)
		return &fin, nil
	}
	return &res2, nil
}

func (s *SliceStorage) createIfEmpty(values []string) SliceValue {
	if _, err := strconv.Atoi(values[0]); err == nil {
		return SliceValue{Kind: KindSliceInt}
	}
	return SliceValue{Kind: KindSliceStr}
}

func (s *SliceStorage) defineKind(key string) []string {
	if s.inner[key].Kind == KindSliceInt {
		return s.inner[key].D_sl
	}
	return s.inner[key].St_sl
}

func (s *SliceStorage) addToAppropriate(key string, values []string, cur SliceValue) {
	if cur.Kind == KindSliceInt {
		for _, x := range values {
			if _, err := strconv.Atoi(x); err != nil {
				s.logger.Info("uncorrect values")
				return
			}
		}
		s.inner[key] = SliceValue{Kind: KindSliceInt, D_sl: values}
	} else {
		s.inner[key] = SliceValue{Kind: KindSliceStr, St_sl: values}
	}
}

func (s *SliceStorage) LPush(key string, values []string) {
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

func (s *SliceStorage) RPush(key string, values []string) {
	val, ok := s.inner[key]
	var tmp []string
	tmp = append(tmp, values...)
	if !ok {
		s.addToAppropriate(key, tmp, s.createIfEmpty(values))

	} else {
		s.addToAppropriate(key, append(s.defineKind(key), tmp...), val)
	}
}

func (s *SliceStorage) RAddToSet(key string, values []string) {
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

func (s *SliceStorage) LPop(key string, indexes ...int) []string {
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

func (s *SliceStorage) RPop(key string, indexes ...int) []string {
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

func (s *SliceStorage) LSet(key string, index int, elem string) (string, error) {
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

func (s *SliceStorage) LGet(key string, index int) (string, error) {
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

func (s *SliceStorage) SaveToFile(filename string) error {
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

func (s *SliceStorage) CheckIfExpired(key string) bool {
	if time.Now().UnixMilli() >= s.inner[key].Expires_at && s.inner[key].Expires_at != 0 {
		s.logger.Info("expired")
		delete(s.inner, key)
		return true
	}
	return false
}

func (s *SliceStorage) Expire(key string, seconds int64) int {
	if res, ok := s.inner[key]; ok {
		res.Expires_at = time.Now().UnixMilli() + seconds
		s.inner[key] = res
		return 1
	}
	return 0
}

func (s *SliceStorage) Clean(file string) {
	var expiredKeys []string
	s.mu.RLock()
	for key, val := range s.inner {
		if time.Now().UnixMilli() >= val.Expires_at && s.inner[key].Expires_at != 0 {
			expiredKeys = append(expiredKeys, key)
		}
	}
	s.mu.RUnlock()
	s.mu.Lock()
	for _, key := range expiredKeys {
		s.logger.Info("Deleting expired key: " + key)
		delete(s.inner, key)
	}
	s.mu.Unlock()
	s.SaveToFile(file)
}

func (s *SliceStorage) IWantToSleepFor(closeChan chan struct{}, interval time.Duration, file string) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(interval):
			s.Clean(file)
		}
	}
}
