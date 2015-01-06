# Configo

## Introduction

Configo provides an easy way to read configuration files in JSON, XML or YAML.


## Documentation

https://godoc.org/github.com/alexcesaro/configo


## Download

    go get github.com/alexcesaro/configo


## Example

```go
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
```
