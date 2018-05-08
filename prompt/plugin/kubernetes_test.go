package plugin

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestKubernetes(t *testing.T) {
	testCases := []struct {
		kubeConfig        string
		expectedContext   string
		expectedNamespace string
		expectedPrompt    string
	}{
		{
			kubeConfig:        "../../testdata/config1",
			expectedContext:   "cluster1_context",
			expectedNamespace: "default",
			expectedPrompt:    "\\[\\033[0m\\]\\[\\033[94m\\]cluster1_context(default)\\[\\033[0m\\]",
		},
		{
			kubeConfig:        "../../testdata/config2",
			expectedContext:   "cluster1_context",
			expectedNamespace: "namespacex",
			expectedPrompt:    "\\[\\033[0m\\]\\[\\033[94m\\]cluster1_context(namespacex)\\[\\033[0m\\]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedContext, func(t *testing.T) {
			os.Setenv("KUBECONFIG", tc.kubeConfig)
			k := &Kubernetes{}
			k.Load(nil)

			if k.context != tc.expectedContext {
				t.Errorf("Expected context: %s, got %s", tc.expectedContext, k.context)
			}
			if k.namespace != tc.expectedNamespace {
				t.Errorf("Expected namespace: %s, got %s", tc.expectedNamespace, k.namespace)
			}

			output, _ := k.Get(termcolor.EscapedFormat)
			if output != tc.expectedPrompt {
				t.Errorf("Expected %s\nGot      %s", tc.expectedPrompt, output)
			}
		})
	}
}
