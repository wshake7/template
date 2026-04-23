package viperc

import "testing"

func TestReadConfig(t *testing.T) {
	m := make(map[string]any)
	_, err := ReadFile("D:\\code\\project\\backend-template\\common\\test.yaml", &m)
	if err != nil {
		t.Error(err)
	}
	for {

	}
}
