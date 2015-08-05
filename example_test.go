package configo_test

import (
	"fmt"

	"github.com/alexcesaro/configo"
)

func main() {
	var conf struct {
		// Field names must start with an uppercase letter.
		User     string
		Password string
	}
	err := configo.Load("config.json", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.User, conf.Password)
}
