package main

import "ConcurrentLogProcessor/processors"

func main() {
	filePth := "log.txt"
	keywords := []string{"INFO", "error", "debug"}

	processors.ProcessLogFile(filePth, keywords)

}
