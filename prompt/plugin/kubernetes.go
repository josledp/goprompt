package plugin

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/josledp/termcolor"
	yaml "gopkg.in/yaml.v2"
)

//Kubernetes is the plugin struct
type Kubernetes struct {
	context   string
	namespace string
}

type k8sconfig struct {
	APIVersion     string            `yaml:"apiVersion"`
	Kind           string            `yaml:"kind"`
	Preferences    map[string]string `yaml:"preferences"`
	CurrentContext string            `yaml:"current-context"`
	Clusters       []struct {
		Name string `yaml:"name"`

		Cluster map[string]string `yaml:"cluster"`
	} `yaml:"clusters"`
	Contexts []struct {
		Name    string            `yaml:"name"`
		Context map[string]string `yaml:"context"`
	} `yaml:"contexts"`
	Users []struct {
		Name string                 `yaml:"name"`
		User map[string]interface{} `yaml:"user"`
	} `yaml:"users"`
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
	file := os.Getenv("KUBECONFIG")
	if file == "" {
		file = os.Getenv("HOME") + string(os.PathSeparator) + ".kube/config"
	}
	if _, err := os.Stat(file); err != nil {
		return nil
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open %s: %v", file, err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %v", file, err)
	}
	var out k8sconfig
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal file %s: %v", file, err)
	}

	k.context = out.CurrentContext
	for _, c := range out.Contexts {
		if c.Name == k.context {
			if _, ok := c.Context["namespace"]; ok {
				k.namespace = c.Context["namespace"]
			} else {
				k.namespace = "default"
			}
			break
		}
	}

	return nil
}

//Get returns the string to use in the prompt
func (k Kubernetes) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if k.context != "" {
		return format(fmt.Sprintf("%s(%s)", k.context, k.namespace), termcolor.FgHiBlue), []termcolor.Mode{termcolor.FgHiBlue}
	}
	return "", nil
}
