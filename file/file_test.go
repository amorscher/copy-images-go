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

func TestThatFilesCanBeFoundInFlatDir(t *testing.T) {

	//GIVEN
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}

	//WHEN
	var result = file.CollectFiles("../test-data/subdir/subsubdir", &filesToCopy, &testExtensions)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 4, len(filesToCopy), "4 Files Must be found")

}

func TestThatFilesCanBeFoundInNestedDir(t *testing.T) {

	//GIVEN
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}

	//WHEN
	var result = file.CollectFiles("../test-data", &filesToCopy, &testExtensions)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 12, len(filesToCopy), "12 Files Must be found")

}

func TestThatCorrectPathIsAttached(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data/subdir/subsubdir"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}

	//WHEN
	var result = file.CollectFiles(testDir, &filesToCopy, &testExtensions)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, path.Join(testDir, "test.gif"), filesToCopy[0].Path)
	assert.Equal(t, path.Join(testDir, "test.jpeg"), filesToCopy[1].Path)
	assert.Equal(t, path.Join(testDir, "test.jpg"), filesToCopy[2].Path)
	assert.Equal(t, path.Join(testDir, "test.png"), filesToCopy[3].Path)

}

func TestThatCorrectDateIsAttached(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data/subdir/subsubdir"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}

	//WHEN
	var result = file.CollectFiles(testDir, &filesToCopy, &testExtensions)

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.Equal(t, 2021, filesToCopy[0].CreateYear)
	assert.Equal(t, time.August, filesToCopy[0].CreatedMonth)

}

func TestThatFilesAreCopiedToTargetDir(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data/subdir/subsubdir"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}
	file.CollectFiles(testDir, &filesToCopy, &testExtensions)
	tempDir := t.TempDir()

	//WHEN
	var result = file.CopyFilesTo(tempDir, filesToCopy)

	//THEN
	var copiedFiles []model.FileInfo
	assert.Nil(t, result, "No error must be thrown")
	file.CollectFiles(tempDir, &copiedFiles, &testExtensions)
	assert.Equal(t, 4, len(copiedFiles), "All 4 files must be copied")

}

func TestThatFilesAreCopiedAndSortedAccordingToCreationDate(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data/subdir/subsubdir"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}
	file.CollectFiles(testDir, &filesToCopy, &testExtensions)
	tempDir := t.TempDir()

	//WHEN
	var result = file.CopyFilesTo(tempDir, filesToCopy)

	//THEN
	var copiedFiles []model.FileInfo
	assert.Nil(t, result, "No error must be thrown")
	file.CollectFiles(tempDir, &copiedFiles, &testExtensions)
	assert.Equal(t, 4, len(copiedFiles), "All 4 files must be copied")
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.gif"), copiedFiles[0].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpeg"), copiedFiles[1].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpg"), copiedFiles[2].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.png"), copiedFiles[3].Path)

}

func TestCopyFilesWithEqualFileNames(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}
	file.CollectFiles(testDir, &filesToCopy, &testExtensions)
	tempDir := t.TempDir()

	//WHEN
	var result = file.CopyFilesTo(tempDir, filesToCopy)

	//THEN
	var copiedFiles []model.FileInfo
	assert.Nil(t, result, "No error must be thrown")
	file.CollectFiles(tempDir, &copiedFiles, &testExtensions)
	assert.Equal(t, 12, len(copiedFiles), "All 12 files must be copied")
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.gif"), copiedFiles[0].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpeg"), copiedFiles[1].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpg"), copiedFiles[2].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.png"), copiedFiles[3].Path)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test_1.gif"), copiedFiles[4].Path)

}

func TestDeleteFilesRemovesFilesFormFileSystem(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}
	file.CollectFiles(testDir, &filesToCopy, &testExtensions)
	tempDir := t.TempDir()
	//copy to temp dir
	file.CopyFilesTo(tempDir, filesToCopy)
	var copiedFiles []model.FileInfo
	//find all the files
	file.CollectFiles(tempDir, &copiedFiles, &testExtensions)

	//WHEN
	//remove them
	file.DeleteFiles(copiedFiles)

	//THEN
	var emptyFiles []model.FileInfo
	file.CollectFiles(tempDir, &emptyFiles, &testExtensions)

	assert.Equal(t, 0, len(emptyFiles), "All files must be deleted")

}

func TestPerpareCopyCreatesDescJson(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data/subdir/subsubdir"
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}
	file.CollectFiles(testDir, &filesToCopy, &testExtensions)
	tempDir := t.TempDir()

	//WHEN
	var result = file.PrepareCopy(tempDir, filesToCopy, "test_desc.json")

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	assert.True(t, fileExists(path.Join(tempDir, "test_desc.json")), "Desc file must exist")

}

func TestPerpareCopyCreatesDescJsonWithCorrectCopyDescs(t *testing.T) {

	//GIVEN
	const testDir string = "../test-data/subdir/subsubdir"
	anbsoluteTestDir, _ := filepath.Abs(testDir)
	fmt.Println(os.Getwd())
	var filesToCopy []model.FileInfo
	var testExtensions []string = []string{".png", ".jpg", ".jpeg", ".gif"}
	file.CollectFiles(testDir, &filesToCopy, &testExtensions)
	tempDir := t.TempDir()

	//WHEN
	var result = file.PrepareCopy(tempDir, filesToCopy, "test_desc.json")

	//THEN
	assert.Nil(t, result, "No error must be thrown")
	descPath := path.Join(tempDir, "test_desc.json")
	assert.True(t, fileExists(descPath), "Desc file must exist")
	input, _ := ioutil.ReadFile(descPath)
	copyDesc := model.FileCopyDescription{Copies: make([]model.FileCopy, 0)}
	json.Unmarshal(input, &copyDesc)
	fmt.Println(copyDesc)
	assert.Equal(t, 4, len(copyDesc.Copies), "4 copy descriptions must be created")
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.gif"), copyDesc.Copies[0].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.gif"), copyDesc.Copies[0].To)
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.jpeg"), copyDesc.Copies[1].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpeg"), copyDesc.Copies[1].To)
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.jpg"), copyDesc.Copies[2].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.jpg"), copyDesc.Copies[2].To)
	assert.Equal(t, path.Join(anbsoluteTestDir, "test.png"), copyDesc.Copies[3].From)
	assert.Equal(t, path.Join(tempDir, "2021", "August", "test.png"), copyDesc.Copies[3].To)

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
