package main

import "os"

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
