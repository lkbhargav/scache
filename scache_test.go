package scache

import (
	"testing"
	"time"
)

var obj Object
var obj2 Object

func init() {
	obj = Init("/Users/lkbhargav/Desktop")
	obj2 = Init("")
}

func TestForBasicSet(t *testing.T) {
	expectedVal := "Bhargav"
	err := obj.Set("name", []byte(expectedVal), 5*time.Second)

	if err != nil {
		t.Errorf("Unexpected error trying to set the key/value. Error: %v", err)
	}

	time.Sleep(1 * time.Second)

	val, err := obj.Get("name")

	if err != nil {
		t.Errorf("Unexpected error trying to read the file. Error: %v", err)
	}

	if val != expectedVal {
		t.Errorf("Expected %v, but instead got %v", expectedVal, val)
	}

	time.Sleep(6 * time.Second)

	val, err = obj.Get("name")

	if err == nil || val != "" {
		t.Errorf("Expected value and error returned to be empty and \"no such key exists\". Instead found value to be %v and error to be %v respectively", val, err)
	}
}

func TestForSetAndRemove(t *testing.T) {
	key := "clover.com_20210101"
	obj.Set(key, []byte("Clover test data"), 5*time.Second)

	time.Sleep(1 * time.Second)

	found := obj.Has(key)

	if !found {
		t.Errorf("Expected found to be true but found it to be false")
	}

	obj.Remove(key)

	time.Sleep(1 * time.Second) // needs a max of one min for the remove to work as it has to kill a goroutine

	found = obj.Has(key)

	if found {
		t.Errorf("Expected found to be false but found it to be true")
	}
}
func TestForFlush(t *testing.T) {
	obj.Set("name", []byte("Bhargav"), 1800*time.Second)
	obj.Set("age", []byte("27"), 1800*time.Second)

	time.Sleep(3 * time.Second)

	obj.Flush()

	time.Sleep(1 * time.Second)

	val, err := obj.Get("name")

	if val != "" {
		t.Errorf("Expected the value to be empty")
	}

	if err.Error() != "no such key exists" {
		t.Errorf("Invalid error message")
	}

	val, err = obj.Get("age")

	if val != "" {
		t.Errorf("Expected the value to be empty")
	}

	if err.Error() != "no such key exists" {
		t.Errorf("Invalid error message")
	}
}

func TestForLongDurationKeys(t *testing.T) {
	defer obj.Flush()
	obj.Set("name", []byte("Monil"), 2*time.Second)
	obj.Set("age", []byte("27"), 8*time.Second)
	obj.Set("something", []byte("something"), 5*time.Second)

	time.Sleep(3 * time.Second)

	obj.Set("name", []byte("Bhargav"), 1800*time.Second)

	time.Sleep(8 * time.Second)

	val, err := obj.Get("name")

	if err != nil || val == "" {
		t.Errorf("Expected value: %v, while we got %v; Expected error: %v, while we got %v", "Bhargav", val, nil, err)
	}
}

func TestForDefaultPaths(t *testing.T) {
	defer obj2.Flush()
	obj2.Set("name", []byte("Bhargav"), 5*time.Second)

	time.Sleep(3 * time.Second)

	val, err := obj2.Get("name")

	if err != nil || val == "" {
		t.Errorf("Expected value: %v, while we got %v; Expected error: %v, while we got %v", "Bhargav", val, nil, err)
	}
}

func TestForSettingDuplicateKeysAndGet(t *testing.T) {
	defer obj.Flush()
	obj.Set("name", []byte("Monil"), 2*time.Second)
	obj.Set("name", []byte("Bhargav"), 1800*time.Second)

	time.Sleep(3 * time.Second)

	val, err := obj.Get("name")

	if err != nil || val == "" {
		t.Errorf("Expected value: %v, while we got %v; Expected error: %v, while we got %v", "Bhargav", val, nil, err)
	}
}
