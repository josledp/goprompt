package prompt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//Config is the struct to fetch the config
type Config struct {
	file   string
	params parameters
}

type parameters struct {
	CustomTemplate string                 `json:"custom_template"`
	Options        map[string]interface{} `json:"options"`
}

//NewConfig loads and returns the config
func NewConfig(file string) (Config, error) {
	var err error
	c := Config{file: file}
	if _, osErr := os.Stat(c.file); os.IsNotExist(osErr) {
		err = os.MkdirAll(path.Dir(c.file), 0755)
		if err != nil {
			return c, fmt.Errorf("unable to create config path %s: %v", path.Dir(c.file), err)
		}
		c.params = parameters{}
		err = c.save()
	} else {
		err = c.load()
		if err != nil {
			return c, fmt.Errorf("Error loading config file %s: %v", file, err)
		}
	}
	return c, err
}

func (c Config) save() error {
	data, err := json.Marshal(c.params)
	if err != nil {
		return fmt.Errorf("unable to marshal config %s: %v", c.file, err)
	}
	err = ioutil.WriteFile(c.file, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write config %s: %v", c.file, err)
	}
	return nil
}

func (c *Config) load() error {
	data, err := ioutil.ReadFile(c.file)
	if err != nil {
		return fmt.Errorf("unable to read config %s: %v", c.file, err)
	}
	err = json.Unmarshal(data, &c.params)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config %s: %v", c.file, err)
	}
	return nil
}

func (c *Config) GetCustomTemplate() (string, bool) {

	return c.params.CustomTemplate, c.params.CustomTemplate != ""
}

func (c *Config) GetOptions() (map[string]interface{}, bool) {
	return c.params.Options, c.params.Options != nil
}
