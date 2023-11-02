package main

import (
	"my_ls/pkg/files"
	"my_ls/pkg/parse"
)

func main() {
	args := parse.GetArgs()
	p := "./"
	if len(args) != 0 {
		p = args[0]
	}

	withRecursive := parse.CheckFlag("-R")
	FS := files.ParseFS(p, withRecursive) // parse FS

	files.FSDisplay(FS, files.FSDisplayOption{
		WithDotfile:   parse.CheckFlag("-a"),
		WithRecursive: withRecursive,
		WidthDetails:  parse.CheckFlag("-l"),
		ReverseResult: parse.CheckFlag("-r"),
		SortByTime:    parse.CheckFlag("-t"),
	})
}
