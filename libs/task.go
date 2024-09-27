package vastcom_infoshare

import (
	"runtime"
	. "fmt"
	_ "log"
	"encoding/json"
	"os"
	"io"
	"strings"
	"regexp"
	"crypto/sha1"
	"io/ioutil"
	"strconv"
	"path/filepath"
	_ "github.com/orcaman/concurrent-map"
	_ "time"
	_ "bufio"
	_ "errors"
	_ "sync"
)

func init() {
	Print("")
}

func UpdateConfigFile() {
	// Get the path to the config file
	var schema string
	schema = Global_ConfigFilePath + ToStr(Global_ConfigFilename)
	configPath, err := filepath.Abs(schema)
	if err != nil {
		Println(err)
		return
	}

	// Read the config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		Println(err)
		return
	}

	// Unmarshal the JSON data into a Config struct
	var config ConfigStruct
	err = json.Unmarshal(data, &config)
	if err != nil {
		Println(err)
		return
	}

	// Check if the MoveLogToScanned field exists
	if config.MoveLogToScanned != "" {
		Println("MoveLogToScanned exists")
		return
	} else {
		// If not present, append it to the config file
		config.MoveLogToScanned = Define_MoveLogToScanned
	}

	// Marshal the updated Config struct back to JSON
	updatedData, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		Println(err)
		return
	}

	// Write the updated JSON data back to the config file
	err = ioutil.WriteFile(configPath, updatedData, os.FileMode(SpawnConfigFilePem))
	if err != nil {
		Println(err)
		return
	}

	// Print the updated value
	Println("MoveLogToScanned:", config.MoveLogToScanned)
}

func InitConfigFile() *META_ConfigFile {
	fd := Global_ConfigFilename
	fd_path := Global_ConfigFilePath + ToStr(fd)

	inst := &META_ConfigFile{
		Fd: fd,
		Fd_Path: fd_path,
	}

	return inst
}

// Check if config file exists
func (load *META_ConfigFile) CheckConfigFileExist_sub_UseFileObject() {
	full_path := Global_ConfigFilePath + ToStr(Global_ConfigFilename)
	Println("Loading from *META_ConfigFile: " + full_path)

	_, err := os.Stat(full_path)
	if os.IsNotExist(err) {
		Println(err)
		load.CheckConfigFileExist_sub_CreateConfigFileByDefault()
	} else if err == nil {
		load.CheckConfigFileExist_sub_LoadConfigFile()
	} else {
		Println("Error checking config file:", err)
		return
	}
}

func (load *META_ConfigFile) CheckConfigFileExist_sub_CreateConfigFileByDefault() {
	full_path := Global_ConfigFilePath + ToStr(Global_ConfigFilename)

	config := ConfigStruct{
		DB_Type:          Default_DB_Type,
		DB_Host:          Default_DB_Host,
		DB_Port:        Default_DB_Port,
		DB_User:        Default_DB_User,
		DB_PW:      	Default_DB_PW,
		Database_Name:  Default_DB_Name,
		Table_Name:    	Default_ST_Name,
		File_Store_Dir:	Define_FileStoreDir,
		Log_Scan_Dir:	Define_LogScanDir,
		ChunkSizeKbPerLog: Default_ChunkLogForKilobyte,
		LogString_TrimPos_Datetime:	Define_LogString_TrimPos_Datetime,
		LogString_TrimPos_FilePath:	Define_LogString_TrimPos_FilePath,
		MoveLogToScanned: Define_MoveLogToScanned,
	}

	configData, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		Println("Error creating config file:", err)
		return
	}

	err = os.WriteFile(full_path, configData, os.FileMode(SpawnConfigFilePem))
	if err != nil {
		Println("Error writing config file:", err)
		return
	}

	return
}

func (load *META_ConfigFile) CheckConfigFileExist_sub_LoadConfigFile() {
	full_path := Global_ConfigFilePath + ToStr(Global_ConfigFilename)
	configData, err := os.ReadFile(full_path)
	if err != nil {
		Println("Error reading config file:", err)
		return
	}

	err = json.Unmarshal(configData, &LoadConfig)
	if err != nil {
		Println("Error parsing config file:", err)
		return
	}

	Printf("Config loaded: %+v\n", LoadConfig)
}

func (load *META_ConfigFile) ProcessConfig() {
	load.CheckConfigFileExist_sub_UseFileObject()
	load.UpdateLogFileCount()
	load.UpdateReadLogChunkSizeKb()
}

func (load *META_ConfigFile) Init_DB_Task() {
	Println("\n")
	db_Connect := ConnectToMssql()
	if (db_Connect == false) {
		os.Exit(1)
	}

	db_CheckDB := CheckDatabase()
	if (db_CheckDB == false) {
		os.Exit(1)
	}

	Println("\n")
	CheckTable()
	Println("\n")
}

func (load *META_ConfigFile) ProcessTask() {
	load.ScanLogDir()
}

func (load *META_ConfigFile) UpdateLogFileCount() {
	load.LogFileCount = DirListFilesCount(LoadConfig.Log_Scan_Dir)
}

func (load *META_ConfigFile) UpdateReadLogChunkSizeKb() {
	load.LogReadChunkByKb = ToInt(LoadConfig.ChunkSizeKbPerLog)
}

func (load *META_ConfigFile) ScanLogDir() {
	workDir := LoadConfig.Log_Scan_Dir
	Println("ScanLogDir at: " + workDir)

    files, err := DirListFiles(workDir)
    if err != nil {
        Println("Error:", err)
        return
    }

    for _, file := range files {
        Println("Log instance found: " + file)
		schema := Sprintf("%s", file)
		load.ReadLog(schema)

		curFileSchema := GetFilename(schema)
		Printf("\n\nCurrent log %s will be rotated/moved out to Scanned folder at %s%s...\n", schema, LoadConfig.MoveLogToScanned, CheckOSPathSlash()+GetFilename(schema))
		MoveFile(schema, LoadConfig.MoveLogToScanned, curFileSchema)
    }
}

func (load *META_ConfigFile) ReadLog(filePath string) {
	Println("File object schema imported as: " + filePath)
	var var_ChunkSizeKbPerLog int = LoadConfig.ChunkSizeKbPerLog
	var_ChunkSizeKbPerLog = 10
	readLogChunkSizeKb := var_ChunkSizeKbPerLog * 1024
	Println("Preset chunk size for " + ToStr(readLogChunkSizeKb) + " bytes per each log instance read.\n")

	file, err := os.Open(filePath)
    if err != nil {
        Println("Error:", err)
        return
    }
    defer file.Close()

    // Specify the chunk size (e.g., X Kilobytes)
    chunkSize := readLogChunkSizeKb

    // Create a buffer to hold the chunk data
    buffer := make([]byte, chunkSize)

	// Read the file in chunks
    for {
        // Read a chunk of data from the file
        n, err := file.Read(buffer)
        if err != nil {
            if err == io.EOF {
                // We've reached the end of the file
                break
            } else {
                Println("Error reading file: ", err)
                return
            }
        }

        // Read a chunk of LoadConfig.ChunkSizeKbPerLog (*1024) Bytes|KiB
		log_chunk := string(buffer[:n])

		// Split the chunk into lines
		lines := strings.Split(string(log_chunk), "\n")
		
		BulkSQLExecute = append(BulkSQLExecute)
		// Process each line in the chunk
		for _, line := range lines {
			if (line != "") {
				// Process the log line
				HandleLogLine(line, filePath)
			}
		}

		decor_hr := "========================================================================================================================================"

		RecordCount := len(BulkSQLExecute)
		if (RecordCount < 1) {
			// 跳到下一個 chunk block 進行處理以防意外地出現讀取錯誤的問題
			continue
		}
		Printf("\n%s\nTotal record(s) available(eligible*): %d row(s).\n\n", decor_hr, RecordCount)

		lastRow_init := BulkSQLExecute[(RecordCount -1):]
		lastRow := replaceStringInSlice(lastRow_init, "),", ")")
		Print("\n\n( ( ( < Last row > ) ) ):\n", lastRow, "\n\n", decor_hr)

		BulkSQLBatch = Sprintf(`
INSERT INTO %s (job_name, job_timestamp, log_timestamp, filename, hash_sha1sum, file_size, full_path)
VALUES
%s
		`, LoadConfig.Table_Name, BulkSQLExecute)

		BulkSQLBatch = strings.ReplaceAll(BulkSQLBatch, "[", "")
		BulkSQLBatch = strings.ReplaceAll(BulkSQLBatch, ",\n]", ";")

		Println(BulkSQLBatch)

		BulkSqlSubmit(BulkSQLBatch)

		// Clear up BulkSQLExecute
		BulkSQLExecute = make([]string, 0)
		//BulkSQLExecute = BulkSQLExecute[0:0]

		Printf("\n\n%s\nInspect clean-up:\n%s\n\n", decor_hr, BulkSQLExecute)
    }
}

func replaceStringInSlice(slice []string, old, new string) []string {
    result := make([]string, len(slice))
    copy(result, slice)

    for i, str := range result {
        result[i] = strings.ReplaceAll(str, old, new)
    }

    return result
}

func HandleLogLine(log_line string, logFile string) {
	keywords := []string{" -> ", "exists", "Number of", "Total file", "Total transferred", "Literal data", "Matched data", "File list", "Total bytes", "bytes/sec", "total size"}
    for _, keyword := range keywords {
        if strings.Contains(log_line, keyword) {
			return
		}
    }

	// Check if the line matches the expected format, ignore other non-relating log strings
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})\s+(.+)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(log_line)

	if len(matches) != 3 {
		return
	}

	pretrim_log_line := strings.ReplaceAll(log_line, "\r", "")
	if strings.HasSuffix(pretrim_log_line, "/") {
		return
	}

	ProcessLogFilter(log_line, logFile)
}

func ProcessLogFilter(log_line string, logFile string) {
	log_timestamp := log_line[:LoadConfig.LogString_TrimPos_Datetime]

	init_file_path := log_line[LoadConfig.LogString_TrimPos_FilePath:]
	file_path := strings.ReplaceAll(init_file_path, "\r", "")

	file_name := RetrieveFilename(file_path)
	
	file_sha1sum := RetrieveFileSha1Checksum(file_path)

	file_size_MiB := RetrieveFileSize_MiB(file_path)

	job_meta_info := RetrieveJobMetaInfo(logFile)

	job_name_init := Explode(job_meta_info[0], CheckOSPathSlash())
	job_name := job_name_init[1]

	job_timestamp_init := Explode(job_meta_info[1], ".")
	job_timestamp := job_timestamp_init[0]

	//ret := Sprintf("%s, %s, %s, '%s', %s, %s, %s", job_name, job_timestamp, log_timestamp, file_path, file_name, file_sha1sum, file_size_MiB)
	//Println(ret)

	sqlBatchQuery := Sprintf("	('%s', '%s', '%s', '%s', '%s', '%s', '%s'),", job_name, job_timestamp, log_timestamp, file_name, file_sha1sum, file_size_MiB, file_path)

	BulkSQLExecute = append(BulkSQLExecute, sqlBatchQuery + "\n")
}

func RetrieveFileSha1Checksum(file_path string) string {
	schema := LoadConfig.File_Store_Dir + CheckOSPathSlash() + file_path
	fileBytes, err := ioutil.ReadFile(schema)
	if err != nil {
		return "null"
	}

	sha1Hash := sha1.Sum(fileBytes)

	// Return the SHA1 checksum as a hexadecimal string
	return Sprintf("%x", sha1Hash)
}

func RetrieveFileSize_MiB(file_path string) string {
	schema := LoadConfig.File_Store_Dir + CheckOSPathSlash() + file_path
	fileInfo, err := os.Stat(schema)
	if err != nil {
		return "0"
	}

	// Convert file size from bytes to megabytes
	fileSizeInBytes := fileInfo.Size()
	fileSizeInMB := float64(fileSizeInBytes)
	fileSizeStr := strconv.FormatFloat(fileSizeInMB, 'f', 4, 32)

	return fileSizeStr + " Bytes"
}

func RetrieveFilename(file_path string) string {
	fileName := filepath.Base(file_path)

	ret := Sprintf("%s", fileName)
	return ret
}

func RetrieveJobMetaInfo(info string) []string {
	meta := Explode(info, "_")
	return meta
}

func CheckOSPathSlash() string {
	switch runtime.GOOS {
	case "windows":
		return "\\"
	case "linux":
		return "/"
	default:
		return "/"
	}
}

