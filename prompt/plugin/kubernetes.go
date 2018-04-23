package plugin

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/josledp/termcolor"
)

//Kubernetes is the plugin struct
type Kubernetes struct {
	context string
}

//Name returns the plugin name
func (Kubernetes) Name() string {
	return "k8s"
}

//Help returns help information about this plugin
func (Kubernetes) Help() (description string, options map[string]string) {
	description = "This plugins show the current context for kubernetes"
	return
}

//Load is the load function of the plugin
func (k *Kubernetes) Load(Prompter) error {
	file := os.Getenv("HOME") + string(os.PathSeparator) + ".kube/config"
	if _, err := os.Stat(file); err != nil {
		return nil
	}
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open ~/.kube/config")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), []byte("current-context: ")) {
			k.context = strings.Replace(string(scanner.Bytes()), "current-context: ", "", 1)
			break
		}
	}

	return nil
}

//Get returns the string to use in the prompt
func (k Kubernetes) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if k.context != "" {
		return format(fmt.Sprintf("%s", strings.Join(strings.Split(k.context, ".")[0:2], ".")), termcolor.Faint), []termcolor.Mode{termcolor.Faint}
	}
	return "", nil
}
