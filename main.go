package main

import (
	. "infoshare/libs"
	. "fmt"
	"os"
)

func pivotEncoding() {
	// Set the character encoding to UTF-8
    os.Setenv("LANG", "zh_CN.UTF-8")
    os.Setenv("LC_ALL", "zh_CN.UTF-8")
}

func main() {
	pivotEncoding()

	TriggerLoadEntrypoint()

	UpdateConfigFile()

	Println("Database connectivity: " + Load_DB_Driver + "\n")

	job := InitConfigFile()
	job.ProcessConfig()
	job.Init_DB_Task()
	job.ProcessTask()
}

