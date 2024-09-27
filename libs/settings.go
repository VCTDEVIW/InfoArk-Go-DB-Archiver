package vastcom_infoshare

import (
	. "fmt"
)

type (
	Filename string
)

type META_ConfigFile struct {
	Fd Filename
	Fd_Path string
	LogFileCount int
	LogReadChunkByKb int
}

// ConfigStruct defines the structure of the config file
type ConfigStruct struct {
	DB_Type         string `json:"DB_Type"`
	DB_Host         string    `json:"DB_Host"`
	DB_Port        	int `json:"DB_Port"`
	DB_User        	string `json:"DB_User"`
	DB_PW      		string `json:"DB_PW"`
	Database_Name  	string `json:"Database_Name"`
	Table_Name    	string `json:"Table_Name"`
	File_Store_Dir  string `json:"File_Store_Dir"`
	Log_Scan_Dir    string `json:"Log_Scan_Dir"`
	ChunkSizeKbPerLog int `json:"ChunkSizeKbPerLog"`
	LogString_TrimPos_Datetime int `json:"LogString_TrimPos_Datetime"`
	LogString_TrimPos_FilePath int `json:"LogString_TrimPos_FilePath"`
	MoveLogToScanned string `json:"MoveLogToScanned"`
}

var (
	Global_ConfigFilePath string = "./"
	Global_ConfigFilename Filename = "InfoShareLogScan.config.json"

	Default_DB_Type string = "msSQL"
	Default_DB_Host string = "localhost"
	Default_DB_Port int = 1433
	Default_DB_User string = "username"
	Default_DB_PW 	string = "P@ssw0rd!234"
	Default_DB_Name string = "InfoShareLogScan"
	Default_ST_Name string = "ScanResult"
	Default_ChunkLogForKilobyte int = 4

	Define_FileStoreDir string = "./"
	Define_LogScanDir 	string = "./"
	Define_MoveLogToScanned string = Define_LogScanDir + "scanned"
	Define_LogString_TrimPos_Datetime int = 23
	Define_LogString_TrimPos_FilePath int = 25

	SpawnConfigFilePem int = 0644

	LoadConfig ConfigStruct
	Load_DB_Driver = ""
	BulkSQLExecute []string
	BulkSQLBatch string
)

func init() {
	Print("")
}

func TriggerLoadEntrypoint() error {
	return nil
}

