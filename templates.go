package main

var templates map[string]string
var defaultOptions map[string]map[string]interface{}

func init() {
	templates = map[string]string{
		"Evermeet": "<(%python%) ><%aws%|><%user%@><%hostname%> <%lastcommand% ><%path%>< %git%><%userchar%> ",
		"Fedora":   "[ <(%python%) ><%aws%|><%user%@><%hostname%> <%lastcommand% ><%path%>< %git%> ]<%userchar%> ",
	}
	defaultOptions = map[string]map[string]interface{}{
		"Evermeet": map[string]interface{}{
			"path.fullpath": true,
		},
		"Fedora": map[string]interface{}{
			"path.fullpath": false,
		},
	}
}
