/*
Package configo provides an easy way to read configuration files in JSON, XML or
YAML.

Example:

	package main

	import (
		"fmt"

		"github.com/alexcesaro/configo"
	)

	func main() {
		var conf struct {
			User     string // Field names must start with an uppercase letter
			Password string
		}
		err := configo.Load("config.json", &conf)
		if err != nil {
			panic(err)
		}
		fmt.Println(conf.User, conf.Password)
	}
*/
package configo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	yml "gopkg.in/yaml.v2"
)

// Load loads the file pointed by filename and write the data to config.
func Load(filename string, config interface{}) error {
	u, err := getUnmarshaler(path.Ext(filename))
	if err != nil {
		return err
	}

	content, err := getFileContent(filename)
	if err != nil {
		return err
	}

	if err = u(content, config); err != nil {
		return err
	}

	return nil
}

// LoadNode only loads the given node of a file and write the data to config.
// It currently only supports YAML config files.
func LoadNode(filename, node string, config interface{}) error {
	ext := path.Ext(filename)
	if path.Ext(filename) != ".yml" {
		return errors.New("configo: LoadNode only supports YAML format")
	}
	u, err := getUnmarshaler(ext)
	if err != nil {
		return err
	}

	content, err := getFileContent(filename)
	if err != nil {
		return err
	}
	content = getYAMLNode(content, node)

	if err = u(content, config); err != nil {
		return err
	}

	return nil
}

func getYAMLNode(content []byte, node string) []byte {
	delim := []byte(node + ":")
	i := bytes.Index(content, delim)
	if i == -1 {
		return nil
	}
	i += len(delim)

	j := i
	for {
		k := bytes.IndexByte(content[j:], '\n')
		if k == -1 {
			return content[i:]
		}
		j += k + 1
		if len(content[j:]) == 0 {
			return content[i:]
		}

		switch content[j] {
		case ' ', '\t', '\n':
			j++
			continue
		default:
			return content[i:j]
		}
	}
}

func getFileContent(filename string) ([]byte, error) {
	if fileNotExist(filename) {
		// If filename does not exist, look into the executable directory
		altFilename := filepath.Join(filepath.Dir(os.Args[0]), filename)
		if fileNotExist(altFilename) {
			return nil, &notFoundError{s: fmt.Sprintf(
				"configo: file not found: %q or %q do not exist",
				filename, altFilename,
			)}
		}

		filename = altFilename
	}
	return readFile(filename)
}

type unmarshaler func([]byte, interface{}) error

func getUnmarshaler(ext string) (unmarshaler, error) {
	switch ext {
	case ".json":
		return json.Unmarshal, nil
	case ".xml":
		return xml.Unmarshal, nil
	case ".yml":
		return yml.Unmarshal, nil
	default:
		return nil, fmt.Errorf("configo: unsupported extension %q,"+
			"supported extensions are json, xml and yml", ext)
	}
}

type notFoundError struct {
	s string
}

func (e *notFoundError) Error() string {
	return e.s
}

// IsNotFound returns whether the error means the config file was not found.
func IsNotFound(err error) bool {
	_, ok := err.(*notFoundError)
	return ok
}

var fileNotExist = func(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

var readFile = ioutil.ReadFile
