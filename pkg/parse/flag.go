package parse

import (
	"os"
	"strings"
)

var isParse = false // parse status

var flags = make(map[rune]struct{}) // flag collection after parse

var args = []string{} // args collection after parse

func parse() {
	if isParse {
		return
	} // do nothing if already parse

	for _, param := range os.Args[1:] { // loop on os args
		if param[0:2] == "--" || param[0:1] == "-" { // check if is flag
			for _, f := range strings.Split(param, "") {
				flags[rune(f[0])] = struct{}{} // save param
			}
		} else { // if is not a flag (arg)
			args = append(args, param)
		}
	}
	isParse = true
}

// return args without name file and flag
func GetArgs() []string {
	parse()
	return args
}

// check if almost one flag exist,
// case sensitive, "-" sensitive
//
// Example:
//
//	hasAll := CheckFlag("-A", "--all")
func CheckFlag(flag ...rune) bool {
	parse()
	for _, f := range flag {
		if _, ok := flags[f]; ok {
			return true
		}
	}
	return false
}
