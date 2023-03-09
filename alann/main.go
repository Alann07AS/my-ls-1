package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
)

func main() {
	root := "./"
	fileSystem := os.DirFS(root)
	user, _ := user.Current()
	fmt.Println(user.HomeDir)
	ls, _ := fs.Glob(fileSystem, "./*")
	fmt.Println(ls)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		// fmt.Println(path)
		// fmt.Println(d.Type().Perm())
		return nil
	})
}
