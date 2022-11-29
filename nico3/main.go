package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		args = []string{"."}
	}
	if len(args) > 1 {
		args = args[1:]
	}
	// fmt.Println(args)

	for _, str := range args {
		fi, err := os.Stat(str)
		if err != nil {
			panic(err)
		}
		fmt.Println(fi.Name())
		fmt.Println(fi.IsDir())
		continue
	}
}
