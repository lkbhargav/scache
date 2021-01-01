package scache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

const extension string = ".scache"
const characters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var defaultPath string = ""

type cacheValueHolder struct {
	Expiry  time.Time
	Channel chan bool
}

// Object => holds the state information
type Object struct {
	scacheKeys map[string]cacheValueHolder
	prefix     string
	path       string
}

// Init => initializes storage cache
func Init(path string) (resp Object) {
	scacheKeys := make(map[string]cacheValueHolder)

	if path == "" {
		path = defaultPath
	} else if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	resp.scacheKeys = scacheKeys
	resp.path = path
	resp.prefix = randomPrefix(10)

	return
}

// Set => sets the key value related data to the storage
func (obj Object) Set(key string, value []byte, expiry time.Duration) (err error) {
	ch := make(chan bool, 1)

	if obj.Has(key) {
		obj.Remove(key)
	}

	obj.scacheKeys[key] = cacheValueHolder{Expiry: time.Now().Add(expiry), Channel: ch}

	err = ioutil.WriteFile(getFilePath(obj, key), value, 0755)

	if err == nil {
		go func() {
			for {
				select {
				case <-time.After(expiry):
					os.Remove(getFilePath(obj, key))
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
		val, err := ioutil.ReadFile(getFilePath(obj, key))
		resp = string(val)
		return resp, err
	}

	err = errors.New("no such key exists")

	return
}

// Remove => helps remove the cache on user request
func (obj Object) Remove(key string) (err error) {
	if obj.Has(key) {
		err = os.Remove(getFilePath(obj, key))
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

func randomPrefix(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func getFilePath(o Object, key string) string {
	return fmt.Sprintf("%v%v_%v%v", o.path, o.prefix, key, extension)
}
