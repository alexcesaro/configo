package configo

import (
	"errors"
	"strings"
	"testing"
)

type testConfig struct {
	User     string
	Password string
	Age      int
}

var expectedConfig = testConfig{
	User:     "myself",
	Password: "my_password",
	Age:      28,
}

var testConfigFiles = map[string]string{
	"json": `{"user": "myself", "password": "my_password", "age": 28}`,
	"xml":  `<?xml version="1.0"?><config><User>myself</User><Password>my_password</Password><Age>28</Age></config>`,
	"yml":  "user: myself\npassword: my_password\nage: 28",
}

func TestJson(t *testing.T) {
	testLoad(t, "json")
}

func TestXml(t *testing.T) {
	testLoad(t, "xml")
}

func TestYml(t *testing.T) {
	testLoad(t, "yml")
}

func TestYmlNode(t *testing.T) {
	fileNotExist = func(filename string) bool {
		return false
	}
	readFile = func(filename string) ([]byte, error) {
		return []byte("node1:\n  user: not_me\n" +
			"\nnode2:\n  " +
			strings.Replace(testConfigFiles["yml"], "\n", "\n  ", -1) +
			"\n\nnode3:\n  user: not_me\n"), nil
	}

	var config testConfig
	if err := LoadNode("config.yml", "node2", &config); err != nil {
		t.Fatal(err)
	}

	if config != expectedConfig {
		t.Errorf("Invalid config, got: %#v, expected: %#v", config, expectedConfig)
	}
}

func TestNotFound(t *testing.T) {
	if IsNotFound(errors.New("Test error")) {
		t.Error("IsNotFound(errors.New()) should be false")
	}

	fileNotExist = func(filename string) bool {
		return true
	}

	var conf struct{}
	err := Load("/fake/file.json", conf)
	if !IsNotFound(err) {
		t.Errorf(`IsNotFound(Load("/fake/file")) should be true, err = %v`, err)
	}
}

func testLoad(t *testing.T, extension string) {
	fileNotExist = func(filename string) bool {
		return false
	}
	readFile = func(filename string) ([]byte, error) {
		return []byte(testConfigFiles[extension]), nil
	}

	var config testConfig
	if err := Load("config."+extension, &config); err != nil {
		t.Fatal(err)
	}

	if config != expectedConfig {
		t.Errorf("Invalid config, got: %#v, expected: %#v", config, expectedConfig)
	}
}
