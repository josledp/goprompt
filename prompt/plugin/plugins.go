package plugin

//Prompter is the interface which provides options/config to the plugin
type Prompter interface {
	GetOption(string) (interface{}, bool)
	GetCache(string) (interface{}, bool)
	Cache(string, interface{}) error
}

type mockPrompt struct {
	options map[string]interface{}
}

func (m mockPrompt) GetOption(key string) (interface{}, bool) {
	value, ok := m.options[key]
	return value, ok
}

func (m mockPrompt) GetCache(key string) (interface{}, bool) {
	return nil, false
}

func (m mockPrompt) Cache(key string, value interface{}) error {
	return nil
}
