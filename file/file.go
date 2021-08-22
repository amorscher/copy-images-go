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
	"time"
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
		var currentImage = model.FileInfo{Path: path, CreationDate: info.ModTime()}
		*files = append(*files, currentImage)
		return nil
	}
}

// CollectFiles collects all files according to the given collectFilesConfig in the provided files array
func CollectFiles(rootDir string, files *[]model.FileInfo, collectFilesConfig CollectFilesConfig) error {
	return filepath.Walk(rootDir, visit(files, collectFilesConfig))
}

// PrepareCopy creates a a json file according to model.FileOperations
// describing all file file operations which would be performend by a real copy
func PrepareCopy(targetDir string, filesToCopy []model.FileInfo, descFileName string, cutoffDate time.Time) error {
	copiedFileNames := make(map[string]int)

	copyDescription := model.FileOperations{FileOperations: make([]model.FileOperation, 0)}
	for _, fileToCopy := range filesToCopy {
		var createionYear int = fileToCopy.CreationDate.Year()
		var creationMonth time.Month = fileToCopy.CreationDate.Month()
		pathWithCreationDate := path.Join(strconv.Itoa(createionYear), creationMonth.String())
		destinationPath := path.Join(targetDir, pathWithCreationDate)

		fileName := path.Base(fileToCopy.Path)
		//check if the file name has already been copied if yes create a new filename in order to not override elements in the destination
		if val, ok := copiedFileNames[path.Base(fileToCopy.Path)]; ok {
			copiedFileNames[path.Base(fileToCopy.Path)] = val + 1
			fileName = strings.Replace(fileName, path.Ext(fileName), "_"+strconv.Itoa(val+1)+path.Ext(fileName), 1)
		} else {
			copiedFileNames[path.Base(fileToCopy.Path)] = 0
		}
		absolutePath, _ := filepath.Abs(fileToCopy.Path)

		// determine the action type of the operation
		opType := operationType(fileToCopy, cutoffDate)

		copyDescription.FileOperations = append(
			copyDescription.FileOperations,
			model.FileOperation{
				From:   absolutePath,
				To:     path.Join(destinationPath, fileName),
				OpType: opType,
			})

	}
	//lets write the json
	// current_time := time.Now()
	// copy_desc_"+current_time.Format("2006-01-02-15:04:05")+".json")
	desc, _ := json.MarshalIndent(copyDescription, "", "     ")
	err := ioutil.WriteFile(path.Join(targetDir, descFileName), desc, 0644)
	fmt.Println(path.Join(targetDir, descFileName) + " written!")

	return err
}

//operationType returns the a valid model.ActionType according to the cutoffDate. All files created on and after the cutoffDate will be copied
func operationType(fileInfo model.FileInfo, cutoffDate time.Time) model.OpType {

	if fileInfo.CreationDate.Before(cutoffDate) {
		return model.MoveOp
	}
	return model.CopyOp
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
		var createionYear int = fileToCopy.CreationDate.Year()
		var creationMonth time.Month = fileToCopy.CreationDate.Month()

		pathWithCreationDate := path.Join(strconv.Itoa(createionYear), creationMonth.String())
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

//DeleteFiles removes all given files from the file-system
func DeleteFiles(files []model.FileInfo) error {
	numberOfFilesToDelete := len(files)
	for index, fileToRemove := range files {
		fmt.Printf("Removing %d/%d %s ... \n", (index + 1), numberOfFilesToDelete, fileToRemove.Path)
		e := os.Remove(fileToRemove.Path)
		if e != nil {
			//if we cannot delete just print a log
			log.Print(e)
		}
	}
	return nil
}

// DeleteFilesCreatedBefore removes all files form the filesystem which have a creation date smaller than provided cutoffDate
func DeleteFilesCreatedBefore(cutoffDate time.Time, files []model.FileInfo) []model.FileInfo {
	//filter the files matching the cutoffDate
	var filteredFiles []model.FileInfo

	for _, file := range files {
		if file.CreationDate.Before(cutoffDate) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	//ok now delete the files
	DeleteFiles(filteredFiles)

	return filteredFiles
}
