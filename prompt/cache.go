package prompt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Cache represents a cache data store
type Cache struct {
	file     string
	data     map[string]interface{}
	modified bool
}

func newCache() (*Cache, error) {
	var err error
	c := Cache{}
	user := os.Getenv("USER")
	if user == "" {
		return nil, fmt.Errorf("Unable to get current user")
	}
	c.file = fmt.Sprintf("/var/tmp/goprompt-%s", user)

	if _, oserr := os.Stat(c.file); !os.IsNotExist(oserr) {
		err = c.load()
	}
	return &c, err
}

func (c *Cache) save() error {
	if !c.modified {
		return nil
	}

	if c.file == "" {
		return fmt.Errorf("Cache not initialized")
	}

	b, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("Unable to marshal cache: %v", err)
	}
	err = ioutil.WriteFile(c.file, b, 0600)
	if err != nil {
		return fmt.Errorf("Unable to save cache: %v", err)
	}
	return nil

}

func (c *Cache) load() error {
	d, err := ioutil.ReadFile(c.file)
	if err != nil {
		return fmt.Errorf("Unable to load cache: %v", err)
	}
	err = json.Unmarshal(d, &c.data)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal cache: %v", err)
	}
	return nil
}

func (c *Cache) get(key string) (interface{}, bool) {
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) set(key string, value interface{}) error {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
	c.modified = true
	return nil
}
