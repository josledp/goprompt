package prompt

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//Config is the struct to fetch the config
type Config struct {
	params parameters
}

type parameters struct {
	Template       string                 `json:"template"`
	CustomTemplate string                 `json:"custom_template"`
	Options        map[string]interface{} `json:"options"`
}

//NewConfigFromFile loads the config from a file and returns the config
func NewConfigFromFile(file string) (*Config, error) {
	var err error
	var c *Config
	if _, osErr := os.Stat(file); os.IsNotExist(osErr) {
		err = os.MkdirAll(path.Dir(file), 0755)
		if err != nil {
			return nil, fmt.Errorf("unable to create config path %s: %v", path.Dir(file), err)
		}
		c = &Config{params: parameters{}}
		rw, err := os.Create(file)
		if err != nil {
			return nil, fmt.Errorf("unable to create initial configuration file: %v", err)
		}
		defer rw.Close()

		err = c.save(rw)
		if err != nil {
			return nil, fmt.Errorf("unable to save initial configuration file: %v", err)
		}
	} else {
		rw, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("error opening config file %s: %v", file, err)
		}
		defer rw.Close()
		return NewConfig(rw)
	}
	return c, nil
}

//NewConfig returns a new Config struct from a io.ReadWriteCloser
func NewConfig(r io.Reader) (*Config, error) {
	c := &Config{}
	err := c.load(r)
	if err != nil {
		return c, fmt.Errorf("error loading config: %v", err)
	}
	return c, nil
}
func (c Config) save(w io.Writer) error {
	data, err := json.Marshal(c.params)
	if err != nil {
		return fmt.Errorf("unable to marshal config: %v", err)
	}
	n, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write config: %v", err)
	}
	if n != len(data) {
		return fmt.Errorf("not all data could be saved: %v", err)
	}
	return nil
}

func (c *Config) load(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("unable to read config: %v", err)
	}
	err = json.Unmarshal(data, &c.params)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config: %v", err)
	}
	return nil
}

//GetTemplate returns the configured predefined template
func (c *Config) GetTemplate() (string, bool) {
	return c.params.Template, c.params.Template != ""
}

//GetCustomTemplate returns the configured custom templatefile
func (c *Config) GetCustomTemplate() (string, bool) {
	return c.params.CustomTemplate, c.params.CustomTemplate != ""
}

//GetOptions return the configured options
func (c *Config) GetOptions() (map[string]interface{}, bool) {
	return c.params.Options, c.params.Options != nil
}
