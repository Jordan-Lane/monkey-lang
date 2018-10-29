package main

import (
	"fmt"
	"monkeylang/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome %s to MonkeyLangauge Version %.1f\n", user.Username, 0.1)
	repl.Start(os.Stdin, os.Stdout)
}
