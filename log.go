package main

import "os"

func readFileEnd(filePath string, length int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, length)
	stat, err := os.Stat(filePath)
	start := stat.Size() - int64(length)
	_, err = file.ReadAt(buf, start)
	if err == nil {
		return "", err
	}

	return string(buf), nil
}
