package file_test

import (
	"copy-images/file"
	"copy-images/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var basicExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

var basicCollectConfig file.CollectFilesConfig = file.CollectFilesConfig{ExcludedDirs: []string{}, SupportedExtensions: basicExtensions}

var basicTestDir string = "../test-data"

func TestThatFilesCanBeFoundInFlatDir(t *testing.T) {

	//GIVEN
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo

	//WHEN
	var result = file.CollectFiles(path.Join(basicTestDir, "subdir", "subsubdir"), &filesToCopy, basicCollectConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 4, len(filesToCopy), "4 Files Must be found")

}

func TestThatFilesCanBeFoundInNestedDir(t *testing.T) {

	//GIVEN
	var filesToCopy []model.FileInfo

	//WHEN
	var result = file.CollectFiles(basicTestDir, &filesToCopy, basicCollectConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 12, len(filesToCopy), "12 Files Must be found")

}

func TestThatExcludedDirsAreIgnored(t *testing.T) {

	//GIVEN
	var filesToCopy []model.FileInfo

	var collectFilesConfig file.CollectFilesConfig = file.CollectFilesConfig{ExcludedDirs: []string{"subdir"}, SupportedExtensions: basicExtensions}

	//WHEN
	var result = file.CollectFiles(basicTestDir, &filesToCopy, collectFilesConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 4, len(filesToCopy), "Only 4 Files Must be found as subdir should be ignored!")

}

func TestThatExcludedDirsCanIncludePaths(t *testing.T) {

	//GIVEN
	var filesToCopy []model.FileInfo

	var collectFilesConfig file.CollectFilesConfig = file.CollectFilesConfig{ExcludedDirs: []string{"subdir/subsubdir"}, SupportedExtensions: basicExtensions}

	//WHEN
	var result = file.CollectFiles(basicTestDir, &filesToCopy, collectFilesConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 8, len(filesToCopy), "Only 8 Files Must be found as subdir/subsubdir should be ignored!")

}

func TestThatExcludedDirsMatchesInCaseSensitive(t *testing.T) {

	//GIVEN
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	caseSensitiveDir := "Subdir/subSubDir"
	var collectFilesConfig file.CollectFilesConfig = file.CollectFilesConfig{ExcludedDirs: []string{caseSensitiveDir}, SupportedExtensions: basicExtensions}

	//WHEN
	var result = file.CollectFiles(basicTestDir, &filesToCopy, collectFilesConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 8, len(filesToCopy), "Only 8 Files Must be found as Subdir/subSubDir should be ignored!")

}

func TestThatCorrectPathIsAttached(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo

	//WHEN
	var result = file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, path.Join(testDir, "test.gif"), filesToCopy[0].Path)
	assert.Equal(t, path.Join(testDir, "test.jpeg"), filesToCopy[1].Path)
	assert.Equal(t, path.Join(testDir, "test.jpg"), filesToCopy[2].Path)
	assert.Equal(t, path.Join(testDir, "test.png"), filesToCopy[3].Path)

}

func TestThatCorrectDateIsAttached(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo

	//WHEN
	var result = file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 2021, filesToCopy[0].CreationDate.Year())
	assert.Equal(t, time.August, filesToCopy[0].CreationDate.Month())

}

func TestThatFilesAreCopiedToTargetDir(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	var filesToCopy []model.FileInfo
	file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)
	tempDir := t.TempDir()

	//WHEN
	var result = file.CopyFilesTo(tempDir, filesToCopy)

	//THEN
	var copiedFiles []model.FileInfo
	assert.Nil(t, result, "No error must be thrown")
	file.CollectFiles(tempDir, &copiedFiles, basicCollectConfig)
	assert.Equal(t, 4, len(copiedFiles), "All 4 files must be copied")

}

func TestThatFilesAreCopiedAndSortedAccordingToCreationDate(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	var filesToCopy []model.FileInfo
	file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)
	tempDir := t.TempDir()

	//WHEN
	var result = file.CopyFilesTo(tempDir, filesToCopy)

	//THEN
	var copiedFiles []model.FileInfo
	assert.Nil(t, result, "No error must be thrown")
	file.CollectFiles(tempDir, &copiedFiles, basicCollectConfig)
	assert.Equal(t, 4, len(copiedFiles), "All 4 files must be copied")
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.gif"), copiedFiles[0].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpeg"), copiedFiles[1].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpg"), copiedFiles[2].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.png"), copiedFiles[3].Path)

}

func TestCopyFilesWithEqualFileNames(t *testing.T) {

	//GIVEN
	var filesToCopy []model.FileInfo
	file.CollectFiles(basicTestDir, &filesToCopy, basicCollectConfig)
	tempDir := t.TempDir()

	//WHEN
	var result = file.CopyFilesTo(tempDir, filesToCopy)

	//THEN
	var copiedFiles []model.FileInfo
	assert.Nil(t, result, "No error must be thrown")
	file.CollectFiles(tempDir, &copiedFiles, basicCollectConfig)
	assert.Equal(t, 12, len(copiedFiles), "All 12 files must be copied")
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.gif"), copiedFiles[0].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpeg"), copiedFiles[1].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpg"), copiedFiles[2].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.png"), copiedFiles[3].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test_1.gif"), copiedFiles[4].Path)

}

func TestDeleteFilesRemovesFilesFromFileSystem(t *testing.T) {

	//GIVEN
	tempDir := t.TempDir()
	//copy to temp dir
	var copiedFiles []model.FileInfo = copyFilesToTemp(basicTestDir, tempDir)

	//WHEN
	//remove them
	file.DeleteFiles(copiedFiles)

	//THEN
	var emptyFiles []model.FileInfo
	file.CollectFiles(tempDir, &emptyFiles, basicCollectConfig)

	assert.Equal(t, 0, len(emptyFiles), "All files must be deleted")

}

func TestDeleteFilesRemovesFilesFromFileSystemCreatedBefore(t *testing.T) {

	//GIVEN
	tempDir := t.TempDir()
	//copy to temp dir
	var copiedFiles []model.FileInfo = copyFilesToTemp(basicTestDir, tempDir)
	//the date
	cutoffDate, _ := time.Parse("2006-01-02", "2021-03-03")
	//manipulate some dates all in the same month
	copiedFiles[0].CreationDate, _ = time.Parse("2006-01-02", "2021-03-02")
	copiedFiles[1].CreationDate, _ = time.Parse("2006-01-02", "2021-03-01")

	//WHEN
	//remove them
	var deletedFiles []model.FileInfo = file.DeleteFilesCreatedBefore(cutoffDate, copiedFiles)

	//THEN
	var keptFiles []model.FileInfo
	file.CollectFiles(tempDir, &keptFiles, basicCollectConfig)

	assert.Equal(t, 2, len(deletedFiles), "2 Files should be deleted")
	assert.Equal(t, 10, len(keptFiles), "10 Files should be kept")

	//data of all kept files has to be after delteBeforeDate
	for _, keptFile := range keptFiles {
		assert.True(t, keptFile.CreationDate.After(cutoffDate))
	}
}

func TestDeleteFilesRemovesFilesFromFileSystemCreatedBefore_DoesNotRemoveExactCutoffDate(t *testing.T) {

	//GIVEN
	tempDir := t.TempDir()
	//copy to temp dir
	var copiedFiles []model.FileInfo = copyFilesToTemp(basicTestDir, tempDir)
	//the date
	cutoffDate, _ := time.Parse("2006-01-02", "2021-03-03")
	//manipulate some dates all in the same month
	//put it on the exact date
	copiedFiles[0].CreationDate = cutoffDate
	copiedFiles[1].CreationDate = cutoffDate
	copiedFiles[2].CreationDate = cutoffDate

	//WHEN
	//remove them
	var deletedFiles []model.FileInfo = file.DeleteFilesCreatedBefore(cutoffDate, copiedFiles)

	//THEN
	var keptFiles []model.FileInfo
	file.CollectFiles(tempDir, &keptFiles, basicCollectConfig)

	assert.Equal(t, 0, len(deletedFiles), "0 Files should be deleted")

}

func TestDeleteFilesRemovesFilesFromFileSystemCreatedBefore_differentMonth(t *testing.T) {

	//GIVEN
	tempDir := t.TempDir()
	//copy to temp dir
	var copiedFiles []model.FileInfo = copyFilesToTemp(basicTestDir, tempDir)
	//the date
	cutoffDate, _ := time.Parse("2006-01-02", "2021-03-03")
	//manipulate some dates all in the same month
	copiedFiles[0].CreationDate, _ = time.Parse("2006-01-02", "2021-03-02")
	// different month before
	copiedFiles[1].CreationDate, _ = time.Parse("2006-01-02", "2021-02-28")

	//WHEN
	//remove them
	var deletedFiles []model.FileInfo = file.DeleteFilesCreatedBefore(cutoffDate, copiedFiles)

	//THEN
	var keptFiles []model.FileInfo
	file.CollectFiles(tempDir, &keptFiles, basicCollectConfig)

	assert.Equal(t, 2, len(deletedFiles), "2 Files should be deleted")
	assert.Equal(t, 10, len(keptFiles), "10 Files should be kept")

	//data of all kept files has to be after delteBeforeDate
	for _, keptFile := range keptFiles {
		assert.True(t, keptFile.CreationDate.After(cutoffDate))
	}
}

func TestDeleteFilesRemovesFilesFromFileSystemCreatedBefore_differentYear(t *testing.T) {

	//GIVEN
	tempDir := t.TempDir()
	//copy to temp dir
	var copiedFiles []model.FileInfo = copyFilesToTemp(basicTestDir, tempDir)
	//the date
	cutoffDate, _ := time.Parse("2006-01-02", "2021-03-03")
	//manipulate some dates
	copiedFiles[0].CreationDate, _ = time.Parse("2006-01-02", "2021-03-02")
	// after form the month and day but before according to year
	copiedFiles[1].CreationDate, _ = time.Parse("2006-01-02", "2020-08-08")

	//WHEN
	//remove them
	var deletedFiles []model.FileInfo = file.DeleteFilesCreatedBefore(cutoffDate, copiedFiles)

	//THEN
	var keptFiles []model.FileInfo
	file.CollectFiles(tempDir, &keptFiles, basicCollectConfig)

	assert.Equal(t, 2, len(deletedFiles), "2 Files should be deleted")
	assert.Equal(t, 10, len(keptFiles), "10 Files should be kept")

	//data of all kept files has to be after delteBeforeDate
	for _, keptFile := range keptFiles {
		assert.True(t, keptFile.CreationDate.After(cutoffDate))
	}
}

func TestPerpareCopyCreatesJson(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)
	tempDir := t.TempDir()

	//WHEN
	var result = file.PrepareCopy(tempDir, filesToCopy, "test_desc.json", time.Now())

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.True(t, fileExists(path.Join(tempDir, "test_desc.json")), "Desc file must exist")

}

func TestPerpareCopyCreatesJsonWithCorrectFileOperations(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	anbsoluteTestDir, _ := filepath.Abs(testDir)
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)
	tempDir := t.TempDir()

	//WHEN
	var result = file.PrepareCopy(tempDir, filesToCopy, "test_desc.json", time.Now())

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	descPath := path.Join(tempDir, "test_desc.json")
	assert.True(t, fileExists(descPath), "Desc file must exist")
	input, _ := ioutil.ReadFile(descPath)
	copyDesc := model.FileOperations{FileOperations: make([]model.FileOperation, 0)}
	json.Unmarshal(input, &copyDesc)
	fmt.Println(copyDesc)
	assert.Equal(t, 4, len(copyDesc.FileOperations), "4 copy descriptions must be created")
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.gif"), copyDesc.FileOperations[0].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.gif"), copyDesc.FileOperations[0].To)
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.jpeg"), copyDesc.FileOperations[1].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpeg"), copyDesc.FileOperations[1].To)
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.jpg"), copyDesc.FileOperations[2].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpg"), copyDesc.FileOperations[2].To)
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.png"), copyDesc.FileOperations[3].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.png"), copyDesc.FileOperations[3].To)

}

func TestPerpareCopyCreatesDescJsonWithCorrectOperationType(t *testing.T) {

	//GIVEN
	var testDir string = path.Join(basicTestDir, "subdir", "subsubdir")
	var filesToCopy []model.FileInfo
	file.CollectFiles(testDir, &filesToCopy, basicCollectConfig)
	tempDir := t.TempDir()
	//the date
	cutoffDate, _ := time.Parse("2006-01-02", "2021-03-03")
	//manipulate some dates all in the same month
	//put it on the exact date
	filesToCopy[0].CreationDate, _ = time.Parse("2006-01-02", "2021-03-02")
	filesToCopy[1].CreationDate, _ = time.Parse("2006-01-02", "2021-03-01")
	filesToCopy[2].CreationDate, _ = time.Parse("2006-01-02", "2021-03-03")

	//WHEN
	var result = file.PrepareCopy(tempDir, filesToCopy, "test_desc.json", cutoffDate)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	descPath := path.Join(tempDir, "test_desc.json")
	assert.True(t, fileExists(descPath), "Desc file must exist")
	input, _ := ioutil.ReadFile(descPath)
	fileOps := model.FileOperations{FileOperations: make([]model.FileOperation, 0)}
	json.Unmarshal(input, &fileOps)
	fmt.Println(fileOps)
	//the ones before the cutoffDate should be moved others should be just copied
	assert.Equal(t, model.MoveOp, fileOps.FileOperations[0].OpType)
	assert.Equal(t, model.MoveOp, fileOps.FileOperations[1].OpType)
	assert.Equal(t, model.CopyOp, fileOps.FileOperations[2].OpType)
	assert.Equal(t, model.CopyOp, fileOps.FileOperations[3].OpType)

}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//copyFilesToTemp copies all files form the sourceDir to the tempDir
func copyFilesToTemp(sourceDir string, tempDir string) []model.FileInfo {
	var filesToCopy []model.FileInfo
	file.CollectFiles(basicTestDir, &filesToCopy, basicCollectConfig)
	//copy to temp dir
	file.CopyFilesTo(tempDir, filesToCopy)
	var copiedFiles []model.FileInfo
	//find all the copied files
	file.CollectFiles(tempDir, &copiedFiles, basicCollectConfig)

	return copiedFiles
}
