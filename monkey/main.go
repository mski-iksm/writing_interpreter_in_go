package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	// userを返す
	// こんな値になるらしい https://blog.suganoo.net/entry/2018/09/11/185131
	// User.Name : user.name
	// User.Uid : 501
	// User.Gid : 503
	// User.Username : hoge
	// User.HomeDir : /home/hoge

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is Monley programming language!\n", user.Username)
	fmt.Print("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
