package files

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"syscall"
	"time"
	"unicode"
)

type FS struct {
	Owner       string
	Group       string
	Size        int64
	ModTime     time.Time
	Name        string
	Inode       uint64
	Permissions string
	IsDir       bool
	FS          *FSarray
}

func ParseFS(path string, withRecursive bool) *FSarray {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer dir.Close()

	// Read the directory contents
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	FSS := FSarray{}

	for _, file := range files {
		// Extract file permissions, owner, group, size, and modification time
		FSS = append(FSS, &FS{
			getOwnerName(file),
			getGroupName(file),
			file.Size(),
			file.ModTime(),
			file.Name(),
			file.Sys().(*syscall.Stat_t).Nlink,
			file.Mode().String(),
			file.IsDir(),
			func() *FSarray {
				if withRecursive && file.IsDir() {
					return ParseFS(path+"/"+file.Name(), true)
				} else {
					return &FSarray{}
				}
			}(),
		})
		// maxTimeWidth, modTime.Format("Jan 2 15:04"),
	}

	sort.Sort(FSS)
	return &FSS
}

// Helper function to get the owner name from the UID
func getOwnerName(fileInfo os.FileInfo) string {
	uid := fileInfo.Sys().(*syscall.Stat_t).Uid
	u, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		return ""
	}
	return u.Username
}

// Helper function to get the group name from the GID
func getGroupName(fileInfo os.FileInfo) string {
	gid := fileInfo.Sys().(*syscall.Stat_t).Gid
	g, err := user.LookupGroupId(fmt.Sprint(gid))
	if err != nil {
		return ""
	}
	return g.Name
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

type FSarray []*FS

func (a FSarray) Len() int      { return len(a) }
func (a FSarray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FSarray) Less(i, j int) bool {
	iRunes := []rune(a[i].Name)
	jRunes := []rune(a[j].Name)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}

	// If the strings are the same up to the length of the shortest string,
	// the shorter string comes first
	return len(iRunes) < len(jRunes)
}
