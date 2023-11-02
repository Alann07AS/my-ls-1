package files

import (
	"fmt"
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

	recu(FS, option, &path, &results)
	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Println()
}

func recu(FS *FSarray, option FSDisplayOption, path *string, results *[]string) {
	result := []string{}
	dirs := FSarray{}
	if option.WithRecursive {
		*results = append(*results, fmt.Sprintf("%s%s:\n", "\033[0m", *path))
	}
	for _, fs := range *FS {
		// ignore dot file
		if !option.WithDotfile && fs.Name[0:1] == "." {
			continue
		}

		colorName := func() string {
			if fs.IsDir {
				return "\033[34m"
			} else {
				return ""
			}
		}()

		if option.WidthDetails {
			result = append(result, fmt.Sprintf("%-*s %-*d %-*s %-*s %-*d %-*s %s%-s\n",
				10, fs.Permissions,
				2, fs.Inode,
				10, fs.Owner,
				10, fs.Group,
				6, fs.Size,
				10, fs.ModTime.Format("Jan 2 15:04"),
				colorName, fs.Name,
			))
		} else {
			result = append(result, fmt.Sprintf("%s%s ", colorName, fs.Name))
		}
		if fs.IsDir {
			dirs = append(dirs, fs)
		}
	}

	res := ""
	if option.ReverseResult {
		for i := len(result) - 1; i >= 0; i-- {
			res += fmt.Sprint("\033[0m", result[i])
		}
	} else {
		for i := 0; i < len(result); i++ {
			res += fmt.Sprint("\033[0m", result[i])
		}
	}
	*results = append(*results, res)

	if option.WithRecursive {
		for _, fs := range dirs {
			// ignore dot file
			if !option.WithDotfile && fs.Name[0:1] == "." {
				continue
			}
			if fs.IsDir {
				lastPath := *path
				*path += "/" + fs.Name
				fmt.Print("\n\n")
				recu(fs.FS, option, path, results)
				*path = lastPath
			}
		}
	}
}
