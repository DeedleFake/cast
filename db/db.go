// Package db provides utilities for dealing with the database.
package db

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/etcd-io/bbolt"
)

var (
	interBucket = []byte("inter")
)

// Open opens the database at the given path.
func Open(p string) (*bbolt.DB, error) {
	return bbolt.Open(p, 0644, nil)
}

// ForEach loops over each entry for a given interface, starting after
// the given time.
//
// If f returns false, the function immediately returns.
func ForEach(db *bbolt.DB, inter string, after time.Time, f func(t time.Time, in, out int) bool) error {
	tx, err := db.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket(interBucket).Bucket([]byte(inter))
	if b == nil {
		return nil
	}

	a := []byte(strconv.FormatInt(after.Unix(), 10))

	c := b.Cursor()
	for k, v := c.Seek(a); k != nil; k, v = c.Next() {
		kn, _ := strconv.ParseInt(string(k), 10, 64)
		t := time.Unix(kn, 0)

		var data inout
		err = json.Unmarshal(v, &data)
		if err != nil {
			return err
		}

		if !f(t, data.in, data.out) {
			break
		}
	}

	return nil
}

// Add adds a record for a given interface.
func Add(db *bbolt.DB, inter string, t time.Time, in, out int) (err error) {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	b, err := tx.CreateBucketIfNotExists(interBucket)
	if err != nil {
		return err
	}
	b, err = b.CreateBucketIfNotExists([]byte(inter))
	if err != nil {
		return err
	}

	k := []byte(strconv.FormatInt(t.Unix(), 10))

	data, err := json.Marshal(&inout{in: in, out: out})
	if err != nil {
		return err
	}

	return b.Put(k, data)
}

type inout struct {
	in, out int
}
