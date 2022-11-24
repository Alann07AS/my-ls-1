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

	//-l
	//type de fichier(premiere lettre)/autorisation/nombre de lien/Propriétaire/Groupe Utilisateur/Taille/Date et heure/Nom de fichier
	//-rw-r--r-- 1 root root 3.0k Jun 24 11:46 test.txt

	//-R
	//affiche les dossiers et les fichiers en récursive

	//-a
	//-r
	//-t
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
			perm := permString(fi.Mode())
			s := fmt.Sprintf("%v %v %s", perm, size, f)
			filesList = append(filesList, s)
			continue
		}
		dirList = addDirList(dirList, f)

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
		dirList = addDirList(dirList, f)

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

func addDirList(list []string, f string) []string {
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

func calSize(s int64) string {
	unit := "B"
	return fmt.Sprintf("%v%v", s, unit)
}

func permString(m os.FileMode) string {
	p := "-rw-r--r--"
	return p
}
