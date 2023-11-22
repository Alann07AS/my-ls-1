package files

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"sort"
	"syscall"
	"time"
	"unicode"
)

type FS struct {
	Owner        string
	Group        string
	Size         int64
	ModTime      time.Time
	Name         string
	Inode        uint64
	Permissions  string
	IsDir        bool
	IsExe        bool
	FS           *FSarray
	TotalBlocks  int64
	SymbolicLink string //*FS
}

func ParseFS(path string, withRecursive bool) *FSarray {
	dir, errRoot := os.Open(path) // get dir

	parentDir, errparent := os.Open(path + "/..") // get parent dir

	FSS := FSarray{}
	files := []fs.FileInfo{}
	if errRoot == nil {
		rootFileInfo, _ := dir.Stat()
		rootFS := NewFS(rootFileInfo, false, &path)
		if rootFileInfo.IsDir() {
			rootFS.Name = "."
		}
		FSS = append(FSS, rootFS)
		defer dir.Close()
		// Read the directory contents
		files, _ = dir.Readdir(-1)
	}

	if errparent == nil {
		parentFileInfo, _ := parentDir.Stat()
		parentFS := NewFS(parentFileInfo, false, &path)
		parentFS.Name = ".."
		FSS = append(FSS, parentFS)
		defer parentDir.Close()
	}

	for _, file := range files {
		// Extract file permissions, owner, group, size, and modification time
		FSS = append(FSS, NewFS(file, withRecursive, &path))
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

type FStime FSarray

func (a FSarray) Len() int      { return len(a) }
func (a FSarray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FSarray) Less(i, j int) bool {
	iRunes := []rune(a[i].Name)
	jRunes := []rune(a[j].Name)

	if iRunes[0] == '.' {
		iRunes = iRunes[1:]
	}
	if jRunes[0] == '.' {
		jRunes = jRunes[1:]
	}

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

func (a FSarray) Get(index int) *FS { return a[index] }

func (a FStime) Len() int      { return len(a) }
func (a FStime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FStime) Less(i, j int) bool {
	return a[j].ModTime.Before(a[i].ModTime)
}

func NewFS(file fs.FileInfo, withRecursive bool, path *string) *FS {
	sl, _ := os.Readlink(*path + "/" + file.Name())

	return &FS{
		getOwnerName(file),
		getGroupName(file),
		file.Size(),
		file.ModTime(),
		file.Name(),
		file.Sys().(*syscall.Stat_t).Nlink,
		file.Mode().String(),
		file.IsDir(),
		file.Mode().Perm()&0o100 != 0,
		func() *FSarray {
			if withRecursive && file.IsDir() {
				return ParseFS(*path+"/"+file.Name(), true)
			} else {
				return &FSarray{}
			}
		}(),
		file.Sys().(*syscall.Stat_t).Blocks,
		sl,
	}
}
