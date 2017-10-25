package prompt

import (
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	os.Setenv("USER", "testing-user")
	c, err := newCache()
	if err != nil {
		t.Fatalf("unable to create new cache: %v", err)
	}
	c.data = map[string]interface{}{
		"data1": float64(10),
		"data2": "string",
		"data3": false,
	}
	err = c.save()
	if err != nil {
		t.Fatalf("unable to save cache: %v", err)
	}

	c2, err := newCache()

	for k, v := range c.data {
		if c2.data[k] != v {
			t.Errorf("expecting %v(%T) got %v(%T)", v, v, c2.data[k], c2.data[k])
		}
	}
}
