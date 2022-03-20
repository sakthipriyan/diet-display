package main

import (
	"fmt"
	"os"
)

func IsDatabaseCreated(filePath string) bool {
	_, err := os.Stat(filePath)
	fmt.Println(err)
	return err == nil
}
