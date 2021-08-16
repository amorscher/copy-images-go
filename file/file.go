package file

import (
	"copy-images/model"
	"copy-images/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//CollectFilesConfig describes the configuration for the CollectFiles function
type CollectFilesConfig struct {
	ExcludedDirs        []string
	SupportedExtensions []string
}

//visit returns a function which collects all fileInfos having the correct file extension
func visit(files *[]model.FileInfo, collectFilesConfig CollectFilesConfig) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		// we skip the dir if it is included in the Excluded dirs
		if info.IsDir() {
			for _, excludedDir := range collectFilesConfig.ExcludedDirs {
				if strings.Contains(strings.ToLower(path), strings.ToLower(excludedDir)) {
					return filepath.SkipDir
				}
			}
		}
		// if the file does not match the  supported extensions or it is a dir we just return
		if !utils.ItemExists(collectFilesConfig.SupportedExtensions, strings.ToLower(filepath.Ext(path))) || info.IsDir() {
			return nil
		}
		var currentImage = model.FileInfo{Path: path, CreatedMonth: info.ModTime().Month(), CreateYear: info.ModTime().Year()}
		*files = append(*files, currentImage)
		return nil
	}
}

// CollectFiles collects all files according to the given collectFilesConfig in the provided files array
func CollectFiles(rootDir string, files *[]model.FileInfo, collectFilesConfig CollectFilesConfig) error {
	return filepath.Walk(rootDir, visit(files, collectFilesConfig))
}

// PrepareCopy creates a a json file according to model.FileCopyDescription describing all file copy actions
func PrepareCopy(targetDir string, filesToCopy []model.FileInfo, descFileName string) error {
	copiedFileNames := make(map[string]int)

	copyDescription := model.FileCopyDescription{Copies: make([]model.FileCopy, 0)}
	for _, fileToCopy := range filesToCopy {

		pathWithCreationDate := path.Join(strconv.Itoa(fileToCopy.CreateYear), fileToCopy.CreatedMonth.String())
		destinationPath := path.Join(targetDir, pathWithCreationDate)

		fileName := path.Base(fileToCopy.Path)
		//check if the file name has already been copied if yes create a new filename in order to not override elements in the destination
		if val, ok := copiedFileNames[path.Base(fileToCopy.Path)]; ok {
			copiedFileNames[path.Base(fileToCopy.Path)] = val + 1
			fileName = strings.Replace(fileName, path.Ext(fileName), "_"+strconv.Itoa(val+1)+path.Ext(fileName), 1)
		} else {
			copiedFileNames[path.Base(fileToCopy.Path)] = 0
		}
		absolutepath, _ := filepath.Abs(fileToCopy.Path)
		copyDescription.Copies = append(copyDescription.Copies, model.FileCopy{From: absolutepath, To: path.Join(destinationPath, fileName)})

	}
	//lets write the json
	// current_time := time.Now()
	// copy_desc_"+current_time.Format("2006-01-02-15:04:05")+".json")
	desc, _ := json.MarshalIndent(copyDescription, "", "     ")
	err := ioutil.WriteFile(path.Join(targetDir, descFileName), desc, 0644)
	fmt.Println(path.Join(targetDir, descFileName) + " written!")

	return err
}

//CopyFilesTo copies all filesToCopy to the targetDir
func CopyFilesTo(targetDir string, filesToCopy []model.FileInfo) error {
	copiedFileNames := make(map[string]int)

	numberOfImagesToCopy := len(filesToCopy)
	for index, fileToCopy := range filesToCopy {

		input, err := ioutil.ReadFile(fileToCopy.Path)
		if err != nil {
			return err
		}
		fmt.Printf("Copying %d/%d %s ... \n", (index + 1), numberOfImagesToCopy, fileToCopy.Path)

		pathWithCreationDate := path.Join(strconv.Itoa(fileToCopy.CreateYear), fileToCopy.CreatedMonth.String())
		destinationPath := path.Join(targetDir, pathWithCreationDate)

		//create the destination path
		err = os.MkdirAll(destinationPath, os.ModePerm)
		if err != nil {
			return err
		}
		fileName := path.Base(fileToCopy.Path)
		//check if the file name has already been copied if yes create a new filename in order to not override elements in the destination
		if val, ok := copiedFileNames[path.Base(fileToCopy.Path)]; ok {
			copiedFileNames[path.Base(fileToCopy.Path)] = val + 1
			fileName = strings.Replace(fileName, path.Ext(fileName), "_"+strconv.Itoa(val+1)+path.Ext(fileName), 1)
		} else {
			copiedFileNames[path.Base(fileToCopy.Path)] = 0
		}

		err = ioutil.WriteFile(path.Join(destinationPath, fileName), input, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteFiles removes all filesToDelete from the system
func DeleteFiles(filesToDelete []model.FileInfo) error {
	for _, fileToRemove := range filesToDelete {
		e := os.Remove(fileToRemove.Path)
		if e != nil {
			log.Fatal(e)
		}
	}
	return nil
}
