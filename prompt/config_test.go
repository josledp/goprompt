package prompt

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	tt := []struct {
		name   string
		config Config
		expect string
	}{
		{
			name: "Full config",
			config: Config{
				params: parameters{
					CustomTemplate: "test",
					Options: map[string]interface{}{
						"option1": "value1",
						"option2": float64(10),
					},
				},
			},
			expect: "{\"custom_template\":\"test\",\"options\":{\"option1\":\"value1\",\"option2\":10}}",
		},
		{
			name: "Only Custom Template",
			config: Config{
				params: parameters{
					CustomTemplate: "test",
					Options:        nil,
				},
			},
			expect: "{\"custom_template\":\"test\",\"options\":null}",
		},
		{
			name: "only options",
			config: Config{
				params: parameters{
					Options: map[string]interface{}{
						"option1": "value1",
						"option2": float64(10),
					},
				},
			},
			expect: "{\"custom_template\":\"\",\"options\":{\"option1\":\"value1\",\"option2\":10}}",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			err := tc.config.save(w)
			if err != nil {
				t.Fatalf("error saving config: %v", err)
			}
			w.Flush()
			if b.String() != tc.expect {
				t.Errorf("expecting %s got %s", tc.expect, b.String())
			}
			c, err := NewConfig(bytes.NewReader(b.Bytes()))
			if err != nil {
				t.Fatalf("error loading config: %v", err)
			}
			if c.params.CustomTemplate != tc.config.params.CustomTemplate {
				t.Errorf("expecting custom template %s got %s", tc.config.params.CustomTemplate, c.params.CustomTemplate)
			}
			if !reflect.DeepEqual(c.params.Options, tc.config.params.Options) {
				t.Errorf("expection options %v got %v", c.params.Options, tc.config.params.Options)
			}
		})
	}
}
