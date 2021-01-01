package scache

import (
	"errors"
	"io/ioutil"
	"os"
	"time"
)

const extension string = ".scache"

type cacheValueHolder struct {
	Expiry  time.Time
	Channel chan bool
}

// Object => holds the state information
type Object struct {
	scacheKeys map[string]cacheValueHolder
}

// Init => initializes storage cache
func Init() (resp Object) {
	scacheKeys := make(map[string]cacheValueHolder)

	resp.scacheKeys = scacheKeys

	return
}

// Set => sets the key value related data to the storage
func (obj Object) Set(key string, value []byte, expiry time.Duration) (err error) {
	ch := make(chan bool, 1)

	if obj.Has(key) {
		obj.Remove(key)
	}

	obj.scacheKeys[key] = cacheValueHolder{Expiry: time.Now().Add(expiry), Channel: ch}

	err = ioutil.WriteFile(key+extension, value, 0755)

	if err == nil {
		go func() {
			for {
				select {
				case <-time.After(expiry):
					os.Remove(key + extension)
					delete(obj.scacheKeys, key)
					return
				case <-obj.scacheKeys[key].Channel:
					return
				}
			}
		}()
	}

	return
}

// Get => gets the key value related data from storage
func (obj Object) Get(key string) (resp string, err error) {
	if obj.Has(key) {
		val, err := ioutil.ReadFile(key + extension)
		resp = string(val)
		return resp, err
	}

	err = errors.New("no such key exists")

	return
}

// Remove => helps remove the cache on user request
func (obj Object) Remove(key string) (err error) {
	if obj.Has(key) {
		err = os.Remove(key + extension)
		obj.scacheKeys[key].Channel <- true
		delete(obj.scacheKeys, key)
	}
	return
}

// Has => returns a bool if the key exists in the cache or not
func (obj Object) Has(key string) (resp bool) {
	return obj.scacheKeys[key].Expiry.Year() > 1
}

// Flush => clears all of the cache completely
func (obj Object) Flush() {
	for key := range obj.scacheKeys {
		obj.Remove(key)
	}
}
