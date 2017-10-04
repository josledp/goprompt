package plugins

//Prompter is the interface which provides options/config to the plugin
type Prompter interface {
	GetOption(string) (interface{}, bool)
	GetConfig(string) (interface{}, bool)
	SetConfig(string, interface{}) error
}

type mockPrompt struct {
	options map[string]interface{}
}

func (m mockPrompt) GetOption(key string) (interface{}, bool) {
	value, ok := m.options[key]
	return value, ok
}

func (m mockPrompt) GetConfig(key string) (interface{}, bool) {
	return nil, false
}

func (m mockPrompt) SetConfig(key string, value interface{}) error {
	return nil
}
