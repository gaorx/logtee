package logtee

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type Config map[string]interface{}

func LoadConfig(filename string) (Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var o map[string]interface{}
	err = json.Unmarshal(content, &o)
	if err != nil {
		return nil, err
	}
	return Config(o), nil
}

func (conf Config) Has(k string) bool {
	_, ok := conf[k]
	return ok
}

func (conf Config) Interface(k string, def interface{}) interface{} {
	v, ok := conf[k]
	if !ok {
		return def
	}
	return v
}

func (conf Config) Str(k, def string) string {
	v := conf.Interface(k, nil)
	if v == nil {
		return def
	}
	if v1, ok := v.(string); ok {
		return v1
	} else {
		return def
	}
}

func (conf Config) Int(k string, def int) int {
	v := conf.Interface(k, nil)
	if v == nil {
		return def
	}
	switch v1 := v.(type) {
	case int:
		return v1
	case int8:
		return int(v1)
	case int16:
		return int(v1)
	case int32:
		return int(v1)
	case int64:
		return int(v1)
	case float32:
		return int(v1)
	case float64:
		return int(v1)
	case json.Number:
		v2, err := v1.Int64()
		if err != nil {
			return def
		}
		return int(v2)
	case string:
		v2, err := strconv.Atoi(v1)
		if err != nil {
			return def
		}
		return v2
	default:
		return def
	}
}

func (conf Config) Bool(k string, def bool) bool {
	v := conf.Interface(k, nil)
	if v == nil {
		return def
	}
	if v1, ok := v.(bool); ok {
		return v1
	} else {
		return def
	}
}

func (conf Config) Object(k string, def map[string]interface{}) map[string]interface{} {
	v := conf.Interface(k, nil)
	if v == nil {
		return def
	}
	switch v1 := v.(type) {
	case map[string]interface{}:
		return v1
	case Config:
		return map[string]interface{}(v1)
	default:
		return def
	}
}

func (conf Config) Sub(k string, def Config) Config {
	v := conf.Object(k, nil)
	if v == nil {
		return def
	}
	return Config(v)
}
