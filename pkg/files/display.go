package files

import (
	"fmt"
	"sort"
)

type FSDisplayOption struct {
	WithDotfile   bool //-a ✓
	WithRecursive bool //-R
	WidthDetails  bool //-l ✓
	ReverseResult bool //-r
	SortByTime    bool //-t
}

func FSDisplay(FS *FSarray, option FSDisplayOption) {
	path := "."
	results := []string{}

	if option.SortByTime {
		temp := FStime(*FS)
		sort.Sort(temp)
		hallo := FSarray(temp)
		FS = &hallo
	}

	if option.ReverseResult {
		reverseFS(FS)
	}

	recu(FS, option, &path, &results)
	// test := ""
	for _, r := range results {
		fmt.Print(r)
	}
	// fmt.Println()
}

func recu(FS *FSarray, option FSDisplayOption, path *string, results *[]string) {
	const RESET = "\033[0m"
	const BLACK = "\033[30m"
	const RED = "\033[31m"
	const GREEN = "\033[32m"
	const YELLOW = "\033[33m"
	const BLUE = "\033[34m"
	const PURPLE = "\033[35m"
	const CYAN = "\033[36m"
	const WHITE = "\033[37m"
	const ORANGE = "\033[38;5;208m"

	result := []string{}
	dirs := FSarray{}
	for _, fs := range *FS {
		// ignore dot file
		if !option.WithDotfile && fs.Name[0:1] == "." {
			continue
		}

		colorName := func() string {
			if fs.IsDir {
				return BLUE + fs.Name + RESET // blue
			} else if fs.IsExe {
				return GREEN + fs.Name + RESET
			} else {
				return fs.Name
			}
		}()
		name := func() string {
			if fs.SymbolicLink != "" {
				return CYAN + fs.SymbolicLink + RESET + " -> " + colorName // blue
			} else {
				return colorName
			}
		}()

		if option.WidthDetails {
			result = append(result, fmt.Sprintf("%-*s %-*d %-*s %-*s %-*d %-*s %s\n",
				10, fs.Permissions,
				2, fs.Inode,
				10, fs.Owner,
				10, fs.Group,
				6, fs.Size,
				10, fs.ModTime.Format("Jan 2 15:04"),
				name,
				// fs.TotalBlocks,
			))
		} else {
			result = append(result, fmt.Sprintf("%s  ", name))
		}
		if fs.IsDir {
			dirs = append(dirs, fs)
		}
	}

	// probleme reverse go run ../main.go -lrR
	if len(result) > 0 {
		// Add a newline character only if there are entries in the current directory.
		res := ""
		if option.ReverseResult {
			for i := len(result) - 1; i >= 0; i-- {
				res += fmt.Sprint(result[i])
			}
		} else {
			for i := 0; i < len(result); i++ {
				res += fmt.Sprint(result[i])
			}
		}
		res += "\n"

		if p := *path; option.WithRecursive && (p[len(p)-1] != '.' || p == "." || p == "..") {
			res = fmt.Sprintf("%s:\n", *path) + res
			res += "\n"
		}

		*results = append(*results, res)
	}

	if option.WithRecursive {
		for _, fs := range dirs {
			// ignore dot file
			if !option.WithDotfile && fs.Name[0:1] == "." {
				continue
			}
			p := *path + "/" + fs.Name
			recu(fs.FS, option, &p, results)
		}
	}
}

func reverseFS(FS *FSarray) {
	temp := make(FSarray, 0)

	for i := FS.Len() - 1; i >= 0; i-- {
		temp = append(temp, FS.Get(i))
	}

	*FS = temp
}
