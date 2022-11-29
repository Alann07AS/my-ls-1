package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO LIST :
	// ls [OPTIONS] [FICHIERS]
	// Implement my-ls without flags or without path
	// Implement my-ls without flags with path
	// Implement my-ls with flags without path
	// Implement my-ls with flags & with path
	// Implements following flags:

	//-l long listing
	//type de fichier(premiere lettre)/autorisation/nombre de lien/Propriétaire/Groupe Utilisateur/Taille/Date et heure/Nom de fichier
	//-rw-r--r-- 1 root root 3.0k Jun 24 11:46 test.txt

	//-R
	//affiche les dossiers et les fichiers en récursive

	//-a showing hidden files
	//-r reverse
	//-t sorts files/directories list by time/date.

	args := os.Args
	firstArg := 1
	longListing := false

	if len(args) > 1 && args[firstArg] == "-l" {
		longListing = true
		firstArg++
	}
	files := args[firstArg:]
	if len(files) == 0 {
		files = []string{"."}
	}
	if longListing {
		showLongListing(files)
	} else {
		showShortListing(files)
	}
}

func showShortListing(files []string) {
	var noFileList []string
	var filesList []string
	var dirList []string

	for _, f := range files {
		fi, err := os.Stat(f)
		if err != nil {
			s := fmt.Sprintf("ls: %v: no file or directory", f)
			noFileList = append(noFileList, s)
			continue
		}
		if !fi.IsDir() {
			filesList = append(filesList, f)
			continue
		}

		dirList = addShortDirList(dirList, f)
	}

	for _, s := range noFileList {
		fmt.Println(s)
	}

	for _, s := range filesList {
		fmt.Println(s)
	}

	for _, s := range dirList {
		fmt.Println(s)
	}
}

func showLongListing(files []string) {
	var noFileList []string
	var filesList []string
	var dirList []string

	for _, f := range files {
		fi, err := os.Stat(f)
		if err != nil {
			s := fmt.Sprintf("ls: %v: no file or directory", f)
			noFileList = append(noFileList, s)
			continue
		}

		if !fi.IsDir() {
			size := calSize(fi.Size())
			perm := fi.Mode().Perm()
			s := fmt.Sprintf("%s %v %s", perm, size, f)
			filesList = append(filesList, s)
			continue
		}
		dirList = addLongDirList(dirList, f)
	}

	for _, s := range noFileList {
		fmt.Println(s)
	}

	for _, s := range filesList {
		fmt.Println(s)
	}

	for _, s := range dirList {
		fmt.Println(s)
	}
}

func addShortDirList(list []string, f string) []string {
	dir, err := os.Open(f)
	if err != nil {
		return list
	}
	filesNames, err := dir.Readdirnames(0)
	if err != nil {
		return list
	}
	list = append(list, "\n"+f+":")
	list = append(list, filesNames...)
	return list
}

func addLongDirList(list []string, f string) []string {
	dir, err := os.Open(f)
	if err != nil {
		return list
	}
	filesNames, err := dir.Readdir(0)
	if err != nil {
		return list
	}
	list = append(list, "\n"+f+":")
	s := ""
	for _, fi := range filesNames {
		size := calSize(fi.Size())
		perm := fi.Mode().Perm()
		s = fmt.Sprintf("%s %v %s", perm, size, fi.Name())
		list = append(list, s)
	}
	return list
}

func calSize(i int64) string {
	s := float64(i)
	unit := "B"

	if (s / 1024) > 1.0 {
		s = s / 1024
		unit = "K"
	}

	if (s / 1024) > 1.0 {
		s = s / 1024
		unit = "M"
	}

	return fmt.Sprintf("%6.2f%v", s, unit)
}
