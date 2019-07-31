package main

import (
	"encoding/json"
	"reflect"

	"go.etcd.io/bbolt"
)

// Config is the configuration for cast.
type Config struct {
	Interfaces []string
}

// DefaultConfig contains the default values for config entries.
var DefaultConfig = Config{
	Interfaces: nil,
}

// GetConfig loads a config from the database.
func GetConfig(db *bbolt.DB) (conf Config, err error) {
	tx, err := db.Begin(false)
	if err != nil {
		return conf, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("config"))
	if b == nil {
		return DefaultConfig, nil
	}

	v := reflect.ValueOf(&conf).Elem()
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		f := v.Field(i)

		c := b.Get([]byte(name))
		if c == nil {
			f.Set(reflect.ValueOf(DefaultConfig).FieldByName(name))
			continue
		}

		err = json.Unmarshal(c, f.Addr().Interface())
		if err != nil {
			return conf, err
		}
	}

	return conf, nil
}
