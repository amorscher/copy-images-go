package main

import (
	"copy-images/file"
	"copy-images/model"
	"fmt"
	"os"
	"time"
)

func main() {

	argsWithoutProg := os.Args[1:]
	mode := argsWithoutProg[0]
	source := argsWithoutProg[1]
	target := argsWithoutProg[2]

	var images []model.FileInfo
	var supportedFileEndings []string = []string{".png", ".jpeg", ".jpg", ".gif"}
	var excludedDirs []string = []string{"Android/Data", ".thumbnails", "WhatsApp/.Shared", "WhatsApp/Media/.Statuses", "WhatsApp/.Thumbs"}

	var collectFilesConfig file.CollectFilesConfig = file.CollectFilesConfig{ExcludedDirs: excludedDirs, SupportedExtensions: supportedFileEndings}

	err := file.CollectFiles(source, &images, collectFilesConfig)
	if err != nil {
		panic(err)
	}
	for _, file := range images {
		fmt.Println(file)
	}

	if mode == "--prepare" {
		fmt.Println("Writing copy description", len(images))
		currentTime := time.Now()
		err = file.PrepareCopy(target, images, "copy_desc_"+currentTime.Format("2006-01-02-15:04:05")+".json")
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Number of files to copy:", len(images))

	if mode == "--copy" {
		//copy the files to the target
		err = file.CopyFilesTo(target, images)

		if err != nil {
			panic(err)
		}

		fmt.Println("Copied all files:", len(images))
	}

}
