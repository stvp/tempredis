package tempredis

import "testing"

type urlTestcase struct {
	config Config
	url    string
}

func TestConfig_URL(t *testing.T) {
	tests := []urlTestcase{
		{
			Config{},
			"redis://0.0.0.0:6379",
		},
		{
			Config{"bind": "127.0.0.1"},
			"redis://127.0.0.1:6379",
		},
		{
			Config{"bind": "127.0.0.1", "port": "1111"},
			"redis://127.0.0.1:1111",
		},
		{
			Config{"bind": "127.0.0.1", "port": "1111", "requirepass": "letmein"},
			"redis://:letmein@127.0.0.1:1111",
		},
	}

	for i, test := range tests {
		got := test.config.URL().String()
		if got != test.url {
			t.Errorf("tests[%d]: got %#v, expected: %#v", i, got, test.url)
		}
	}
}
