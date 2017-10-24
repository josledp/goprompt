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
	Templates map[string]string                 `json:"templates"`
	Options   map[string]map[string]interface{} `json:"options"`
}

var defaultTemplates map[string]string

func init() {
	defaultTemplates = map[string]string{
		"Evermeet": "<(%python%) ><%aws%|><%user%@><%hostname%> <%lastcommand% ><%path%>< %git%><%userchar%> ",
		"Fedora":   "[ <(%python%) ><%aws%|><%user%@><%hostname%> <%lastcommand% ><%path%>< %git%> ]<%userchar%> ",
	}
}

//New loads and returns the config
func NewConfig(file string) (Config, error) {
	var err error
	c := Config{file: file}
	if _, osErr := os.Stat(c.file); os.IsNotExist(osErr) {
		err = os.MkdirAll(path.Dir(c.file), 0755)
		if err != nil {
			return c, fmt.Errorf("unable to create config path %s: %v", path.Dir(c.file), err)
		}
		c.params = parameters{
			Templates: defaultTemplates,
			Options: map[string]map[string]interface{}{
				"Evermeet": map[string]interface{}{
					"path.fullpath": 1,
				},
				"Fedora": map[string]interface{}{
					"path.fullpath": 0,
				},
			},
		}
		err = c.save()
	} else {
		err = c.load()
		if _, ok := c.params.Templates["Evermeet"]; !ok {
			if c.params.Templates == nil {
				c.params.Templates = make(map[string]string)
			}
			for k, v := range defaultTemplates {
				c.params.Templates[k] = v
			}
			err = c.save()
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

//GetTemplates return current configured templates
func (c Config) GetTemplates() []string {
	t := make([]string, len(c.params.Templates))
	i := 0
	for k := range c.params.Templates {
		t[i] = k
		i++
	}
	return t
}

//GetTemplate returns a template
func (c Config) GetTemplate(template string) (string, bool) {
	t, ok := c.params.Templates[template]
	return t, ok
}

//GetTemplateOptions return current configured options for template
func (c Config) GetTemplateOptions(template string) map[string]interface{} {
	return c.params.Options[template]
}
