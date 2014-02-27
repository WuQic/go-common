// The messages package looks for i18n folders within the current
// directory and GOPATH and loads them into the system
package localize

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/goinggo/tracelog"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	// T is the translate function for the specified user
	// locale and default locale specified during the load
	T i18n.TranslateFunc
)

// Load looks for i18n folders inside the current directory and the GOPATH
// to find translation files to load
func Load(userLocale string, defaultLocal string) error {
	gopath := os.Getenv("GOPATH")
	pwd, err := os.Getwd()
	if err != nil {
		tracelog.COMPLETED_ERROR(err, "messages", "Load")
		return err
	}

	tracelog.INFO("messages", "Load", "PWD[%s] GOPATH[%s]", pwd, gopath)

	// Load any translation files we can find
	searchDirectory(pwd)
	if gopath != "" {
		searchDirectory(gopath)
	}

	// Create a translation function for use
	T, err = i18n.Tfunc(userLocale, defaultLocal)
	if err != nil {
		return err
	}

	return nil
}

// searchDirectory recurses through the specified directory looking
// for i18n folders. If found it will load the translations files
func searchDirectory(directory string) {
	// Read the directory
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		tracelog.COMPLETED_ERROR(err, "messages", "searchDirectory")
		return
	}

	// Look for i18n folders
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() == true {
			// Is this an i18n folder
			if fileInfo.Name() == "i18n" {
				loadTranslationFiles(fmt.Sprintf("%s/%s", directory, fileInfo.Name()))
				continue
			}

			// Look for more sub-directories
			searchDirectory(fmt.Sprintf("%s/%s", directory, fileInfo.Name()))
			continue
		}
	}
}

// loadTranslationFiles loads the found translation files into the i18n
// messaging system for use by the application
func loadTranslationFiles(directory string) {
	// Read the directory
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		tracelog.COMPLETED_ERROR(err, "messages", "loadTranslationFiles")
		return
	}

	// Look for JSON files
	for _, fileInfo := range fileInfos {
		if path.Ext(fileInfo.Name()) != ".json" {
			continue
		}

		fileName := fmt.Sprintf("%s/%s", directory, fileInfo.Name())

		tracelog.INFO("messages", "loadTranslationFiles", "Loading %s", fileName)
		i18n.MustLoadTranslationFile(fileName)
	}
}
