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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	yml "gopkg.in/yaml.v2"
)

// Load loads the file pointed by filename and write the data to config.
func Load(filename string, config interface{}) error {
	if fileNotExist(filename) {
		// If filename does not exist, look into the executable directory
		altFilename := filepath.Join(filepath.Dir(os.Args[0]), filename)
		if fileNotExist(altFilename) {
			return fmt.Errorf(
				"configo: file not found: %q or %q do not exist",
				filename,
				altFilename,
			)
		}

		filename = altFilename
	}
	content, err := readFile(filename)
	if err != nil {
		return err
	}

	ext := path.Ext(filename)

	switch ext {
	case ".json":
		err = json.Unmarshal(content, config)
	case ".xml":
		err = xml.Unmarshal(content, config)
	case ".yml":
		err = yml.Unmarshal(content, config)
	case "":
		return fmt.Errorf("configo: config file has no extension %q", filename)
	default:
		err = fmt.Errorf("configo: unknown extension %q, "+
			"known extensions are json, xml and yml", ext)
	}
	if err != nil {
		return err
	}

	return nil
}

var fileNotExist = func(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

var readFile = ioutil.ReadFile
