package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type cmdParams struct {
	fileNamePrefix string
	dirNamePrefix  string
	fileCount      int
	dirCount       int
	path           string
}

func validateCmdParams(params cmdParams) {

	validationFail := false
	if params.fileCount == 0 {
		fmt.Println("ERROR: File count (filecount) cannot be 0")
		validationFail = true
	}

	if params.dirCount > 0 && params.dirNamePrefix == "" {
		fmt.Println("ERROR: Directory name prefix (dirnameprefix) cannot be empty if dircount > 0")
		validationFail = true
	}

	if params.path == "" {
		fmt.Println("ERROR: Path (path) cannot be empty")
		validationFail = true
	}

	if params.fileNamePrefix == "" {
		fmt.Println("ERROR: File name prefix (filenameprefix) cannot be empty")
		validationFail = true
	}

	if validationFail {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func getCmdParams() cmdParams {

	fileNamePrefix := flag.String("filenameprefix", "", "The prefix will be a part of the file name created. (Required)")
	directoryNamePrefix := flag.String("dirnameprefix", "", "The prefix will be a part of the file name created. (Conditional - Needed if dircount > 0)")
	fileCount := flag.Int("filecount", 0, "Total number of files to be created. (Required)")
	directoryCount := flag.Int("dircount", 0, "Total number of directories to create. (Optional)")
	path := flag.String("path", "", "Path where the files/directories have to be created. (Required)")

	flag.Parse()

	cmdPrms := cmdParams{
		*fileNamePrefix, *directoryNamePrefix, *fileCount, *directoryCount, *path,
	}

	validateCmdParams(cmdPrms)

	return cmdPrms
}

func createFile(fpath string) {
	bytes := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		bytes[i] = byte(i)
	}

	file, _ := os.Create(fpath)
	file.Write(bytes)
	file.Close()
}

func createFilesInDirs(dirPath string, fileCount int, fileNamePrefix string, wg *sync.WaitGroup) {

	os.MkdirAll(dirPath, os.ModePerm)

	for j := 0; j < fileCount; j++ {
		fname := dirPath + "\\" + fileNamePrefix + "_" + strconv.Itoa(j) + ".txt"
		createFile(fname)
	}

	wg.Done()
}

func start(cmdPrms cmdParams) {
	if cmdPrms.dirCount == 0 {
		for j := 0; j < cmdPrms.fileCount; j++ {
			fname := cmdPrms.path + cmdPrms.fileNamePrefix + "_" + strconv.Itoa(j) + ".txt"
			createFile(fname)
		}
	} else {
		var wg sync.WaitGroup
		for i := 0; i < cmdPrms.dirCount; i++ {
			newpath := filepath.Join(cmdPrms.path, cmdPrms.dirNamePrefix+"_"+strconv.Itoa(i))
			wg.Add(1)
			go createFilesInDirs(newpath, cmdPrms.fileCount, cmdPrms.fileNamePrefix, &wg)
		}
		wg.Wait()
	}
}

func main() {
	cmdPrms := getCmdParams()

	start(cmdPrms)

}
