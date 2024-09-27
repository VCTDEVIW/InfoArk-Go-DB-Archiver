#!/bin/bash
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
## ---- * ---- * ---- * ----
function genConfigFile() {
    if [[ ! -f "$configFileName" ]]; then
        echo -e "\nConfiguration file \"$configFileName\" NOT found as ./$configFileName !\n\nThe file is now generated with default settings at ./$configFileName.";

        cat <<'EOF' >> ./$configFileName
#!/bin/bash
## ***********************************************************************************
## Caution! This configuration file is part of the Bash Shell script, which stands for
## it is formatted in Linux Bash as well, please be careful when modify this file.
##
## Created by Vastcom Technology Limited
## ***********************************************************************************
##
##
use_configFile=1
#### [mysql]: Configure the DSN(Data Source Name) of MySQL database server with 
## the connection string for logging bulk records.
## ↓ ↓
mysql_host="127.0.0.1"
mysql_port=3306
mysql_usr="root"
mysql_pwd="P@ssw0rd"
mysql_db_name="InfoCoreCheckScan"
mysql_table_name="ScanResult"
##
##
##
#### [infocore-log]: Define wherever the underlying log path for the cron-job.
## ↓ ↓
__DIR__=$(pwd)
soleLogFile="log_100.txt"
soleLogPath="$__DIR__/$soleLogFile"
#soleLogPath="./samples/vastcom/test1234_202407151749500.log"
##
##
##
#### [output-temp]: The script will create temporary files when processing, 
## it provides an option $keepSqlBatchFile to allow persistence of 
## the intermediate .sql (MySQL source format) files.
## Usage:   keepSqlBatchFile="yes" | "no"   # Set to yes|no for your choice.
## The $batchGenFile batch file is to be generated automatically,
## it requires to be placed under a directory, the default is ./sql/ for instance.
## ↓ ↓
keepSqlBatchFile="yes"
mission_timestamp=$(date +"%Y_%m_%d_%H%M_t%Ss")
batchGenFile="bulk_InfoCoreCheckScan_$mission_timestamp.sql"
batchGenFilePath="sql/$batchGenFile"
outputDir=""
outputDirPath="$outputDir"
scanByPassToken="exists"
samplingDir="./TestOutput"
samplingPath="$samplingDir"
global_workDir="$__DIR__/TestOutPut"
global_workPath="$global_workDir"
global_scanJobLogFolder="./"
#global_scanJobLogFolder="/mnt/hgfs/rhel_8.8_projects/Projects/InfoCore/samples/vastcom"
#global_scanJobLogFolder="./samples/vastcom"
EOF
    fi
}

configFileName="InfoCoreCheckScan_Config.txt"
genConfigFile
source $configFileName
if [[ "$use_configFile" == 1 ]]; then
    source $configFileName
else
    use_configFile=0
    mysql_host="127.0.0.1"
    mysql_port=3306
    mysql_usr="root"
    mysql_pwd="P@ssw0rd"
    mysql_db_name="InfoCoreCheckScan"
    mysql_table_name="ScanResult"
    __DIR__=$(pwd)
    soleLogFile="log_100.txt"
    soleLogPath="$__DIR__/$soleLogFile"
    keepSqlBatchFile="yes"
    mission_timestamp=$(date +"%Y_%m_%d_%H%M_t%Ss")
    batchGenFile="bulk_InfoCoreCheckScan_$mission_timestamp.sql"
    batchGenFilePath="sql/$batchGenFile"
    outputDir=""
    outputDirPath="$outputDir"
    scanByPassToken="exists"
    samplingDir="./TestOutput"
    samplingPath="$samplingDir"
    global_workDir="$__DIR__/TestOutPut"
    global_workPath="$global_workDir"
    global_scanJobLogFolder="./"
    #global_scanJobLogFolder="/mnt/hgfs/rhel_8.8_projects/Projects/InfoCore/samples/vastcom"
fi
## ---- * ---- * ---- * ----

function trimPathWithBlankspace() {
    local param="$@"
    #local param="2024-07-10 11:30:58.906 two5w-abcrevo test5w/upload - blank/new - FinancelCWP_2024-07-09_000008 - copy.jpg"
    local fetchURL=$(echo "$param"| cut -d ' ' -f3-)
    local composed=""

    # Split the input string using the forward slash as the delimiter
    IFS='/' read -ra var_parts <<< "$fetchURL"

    local trimFromSubFolder=("${var_parts[@]:1}")
    local findFirstSub=0
    local firstSubDir=""
    for part in "${trimFromSubFolder[@]}"; do
        if [[ $findFirstSub -eq 0 ]]; then
            firstSubDir=$part
        fi

        if [[ $part == *" "* ]]; then
            #composed+="\"$part\"/";
            composed+="'$part'/";
        else
            #composed+="$part/";
            composed+="$part/";
        fi

        findFirstSub=1
    done

    # Properize the relatively base directory imported in parameter
    local initString="${var_parts[@]:0}"
    local trimmed=$(echo "$initString"| awk -F "$firstSubDir" '{print $1}')
    trimmed=$(echo "$trimmed"| rev | sed 's/^[[:space:]]*//g' | rev)
    trimmed=$(echo "$trimmed"| sed 's/^ *//')
    local head=""
    if [[ $trimmed == *" "* ]]; then
        #head="\"$trimmed\"/";
        head="'$trimmed'/";
    else
        #head="$trimmed/";
        head="$trimmed/";
    fi

    # Remove the trailing forward slash
    composed=${composed%/}
    echo "$head""$composed"
}

function filter() {
    declare -A filterList
    local filterList['pattern_fileOnly']=*$outputDirPath'/'*'.'*
    local filterList['pattern_dirOnly']=*$outputDirPath'/'*'.'*

    declare -A argv_modeOpt
    local argv_modeOpt['file']=1
    local argv_modeOpt['dir']=1
    local argv_modeOpt['info']=1

    local modeOpt="$1"
    local fp="$2"

    if [[ ! -z "$fp" ]]; then
        local obj="$fp"
    else
        local obj="$soleLogPath"
    fi

    if [[ -z $modeOpt ]]; then { echo -e "\nEmpty modeOpt!"; exit 1; } fi
    if [[ "${argv_modeOpt[$modeOpt]:-}" == "" ]]; then { echo -e "\nInvalid modeOpt: \"$modeOpt\""; exit 1; } fi

    if [[ $modeOpt == 'file' ]]; then {
        while IFS= read -r line; do
        if [[ "$line" == ${filterList['pattern_fileOnly']} ]]; then
            if [[ "$line" == *'/'*' '*'/'* ]]; then
                if [[ ! "$line" == *"$scanByPassToken"* ]]; then
                    local str1=$(echo $line| cut -d ' ' -f1)
                    local str2=$(echo $line| cut -d ' ' -f2)
                    local resolve=$(trimPathWithBlankspace "$line")
                    echo -e "$str1 $str2 $resolve"
                fi
            else
                if [[ ! "$line" == *"$scanByPassToken"* ]]; then
                    local str1=$(echo $line| cut -d ' ' -f1)
                    local str2=$(echo $line| cut -d ' ' -f2)
                    local resolve=$(trimPathWithBlankspace "$line")
                    echo -e "$str1 $str2 $resolve"
                fi
            fi
        fi
        done < $obj
        exit 1
    }
    fi

    if [[ $modeOpt == 'dir' ]]; then {
        local ignoreToken=("Number of*" "Total transferred*" "*->*" "Total file*" "Literal data*" "Matched data*" "File list*" "Total bytes*" "sent*" "total size*")
        while IFS= read -r line; do
        if [[ "$line" != ${filterList['pattern_dirOnly']} ]] && [[ "$line" == *$outputDirPath* ]] && [[ "$line" != *'->'* ]]; then
            if [[ "$line" == *'/'*' '*'/'* ]]; then
                local okay=1
                local str1=$(echo $line| cut -d ' ' -f1)
                local str2=$(echo $line| cut -d ' ' -f2)
                for scan in "${ignoreToken[@]}"; do
                    if [[ "$line" == $scan ]]; then
                        okay=0
                        break
                    fi
                done

                if [[ "$okay" -eq 1 ]]; then
                    local resolve=$(trimPathWithBlankspace "$line")
                    echo -e "$str1 $str2 $resolve"
                fi
            else
                local okay=1
                local str1=$(echo "$line"| cut -d ' ' -f1)
                local str2=$(echo "$line"| cut -d ' ' -f2)
                local line=$(echo "$line"| cut -d ' ' -f4-)
                for scan in "${ignoreToken[@]}"; do
                    if [[ "$line" == $scan ]]; then
                        okay=0
                        break
                    fi
                done

                if [[ $okay -eq 1 ]]; then
                    echo -e "$str1 $str2 $line"
                fi
            fi
        fi
        done < $obj

        while IFS= read -r line; do
            local date=$(echo $line| cut -d ' ' -f1)
            local timestamp=$(echo $line| cut -d ' ' -f2)
            local get_dir=$(echo $line| cut -d ' ' -f3-| tr -d '\r')
            get_dir=$(echo "$get_dir"| sed "s/\/\//\//g")
            get_dir=$(dirname "$get_dir")
            echo "$date $timestamp $get_dir"
        done <<< $(filter file $obj)
        exit 1
    }
    fi

    if [[ $modeOpt == 'info' ]]; then {
        local findToken=("Number of*" "Total transferred*" "*->*" "Total file*" "Literal data*" "Matched data*" "File list*" "Total bytes*" "sent*" "total size*")
        while IFS= read -r line; do
            local sourceURL=$(echo "$line"| cut -d ' ' -f4-)
            local okay=0
            for scan in "${findToken[@]}"; do
                if [[ "$sourceURL" == $scan ]]; then
                    okay=1
                    break
                fi
            done

            if [[ $okay -eq 1 ]]; then
                echo $line
            fi
        done < $obj
        exit 1
    }
    fi
}

function createDump() {
    local argv1="$1"
    local argv2="$2"
    local forceRun=0
    if [[ "$argv1" == '-f' ]] || [[ "$argv1" == '--force' ]]; then { forceRun=1; } fi
    if [[ "$argv2" == '-f' ]] || [[ "$argv2" == '--force' ]]; then { forceRun=1; } fi

    if [[ "$forceRun" == 0 ]]; then
    #### → → → →
    echo -e "\nDo you want to continue?\nThis process will create all empty-dump folders and files according to the specified \$soleLogFile\n$soleLogPath\nDespite they are all empty files, it will still pertain some exact meta file sizes to store due to OS file system merchanism.\nFor 10k+ scale empty files, it may consume up to 50 MiB or excessive when in earlier stage and it **may** seize to occupy after a while then drop to 0 MiB.\nIt consumes partial CPU usage and OS storage disk IO transaction, please be patience while settling the proceedings."
    echo -e "\n請注意！此過程將會完全按照所指定的日誌檔案中，所有的目錄及檔案逐一對照地建立，當然全部都是大小為 0 Bytes 的空文件。\n但由於 Linux OS 檔案系統的運作機制，包括諸如 inode+索引+文件Metadata 等技術，仍然會佔據一定大小的 OS File system 空間，例如 10萬個左右 的空文件可能需要佔到約估 50 MiB 的大小，之後有可能建立後 OS 緩衝處理完後回退到更少甚至 0 Bytes 的體積亦有可能不會。\n檔案建立過程會佔據一定的 CPU 使用量 及儲存碟的 IO 運算及請耐心等候！\n\n"
    echo "請輸入 (不分大小寫) Y 或 y 然後按 [Enter] 確認鍵以繼續："
    read -p "Type 'Y' or 'y' then press [Enter] to continue:    " confirmRun

    if [[ "$confirmRun" != 'Y' ]] && [[ "$confirmRun" != 'y' ]]; then
        echo -e "\nOperation aborted.\n"
        exit 1
    fi
    printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =

    echo -e "\nEnable automatic content insertion into dump files to diversify their hash(sha1sum)?\nType 'Y' or 'y' then press [Enter] to continue:\n"
    echo -e "默認下，由於所有檔案都是空文件因此它們的 hash(sha1sum) 值都是相同的。但您想啟動自動插入內容，讓它們能變成各自不同的 hash(sha1sum) checksum 值嗎?\n"
    read -p "請輸入 (不分大小寫) Y 或 y 然後按 [Enter] 確認鍵以繼續：   " confirmVary

    if [[ "$confirmVary" != 'Y' ]] && [[ "$confirmVary" != 'y' ]]; then
        echo -e "\nDiversified hash(sha1sum) file checksum disabled.\n"
        local varyHash=0
    else
        echo -e "\nDiversified hash(sha1sum) file checksum enabled.\n"
        local varyHash=1
    fi
    printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =

    echo -e "\nProcess begin in 3 seconds...\n"
    for n in {3..0}; do
        echo "Start in $n"
        sleep 1
    done
    #### → → → →
    fi

    if [ -d "$samplingPath" ]; then
        rm -rf "$samplingPath"
    fi

    # Retrieve all residing directories for creation in prior to settle sample files by touch
    while IFS= read -r line; do
        local get_dir=$(echo $line| cut -d ' ' -f3-| tr -d '\r')
        get_dir=$(echo "$get_dir"| sed "s/\/\//\//g")
        mkdir -p "$samplingPath/$get_dir"
        echo -e "\nHaving folder created at: $samplingPath/$get_dir"
    done <<< $(filter dir)

    # Retrieve all sample files to be touched according to the $soleLogFile
    while IFS= read -r line; do
        local get_path=$(echo $line| cut -d ' ' -f3-| tr -d '\r')
        local touch_fp=$(echo "$samplingPath/$get_path"| sed "s/\/\//\//g")
        touch $touch_fp

        if [[ $varyHash == 1 ]]; then
            echo "$get_path" > "$touch_fp"
        fi

        echo -e "\nFile created at: $touch_fp"
    done <<< $(filter file)

    echo -e "\n\nTotal files size:"
    du -sh "$samplingPath"
}

function createMysqlImport_SqlFile() {
cat <<EOF > $__DIR__/"$batchGenFilePath"
-- Check if the 'Default: InfoCoreCheckScan' database exists, and create it if it doesn't
-- utf8mb4_0900_ai_ci may NOT be available for legacy MySQL version lower than 8.4
-- CREATE DATABASE IF NOT EXISTS InfoCoreCheckScan CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
CREATE DATABASE IF NOT EXISTS $mysql_db_name CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the 'Default: InfoCoreCheckScan' database
USE $mysql_db_name;

-- Create a new table called '$mysql_table_name'
CREATE TABLE IF NOT EXISTS $mysql_table_name (
    rid INT AUTO_INCREMENT PRIMARY KEY,
    job_name VARCHAR(255),
    job_timestamp VARCHAR(16),
	log_timestamp TIMESTAMP,
    filename VARCHAR(255),
	hash_sha1sum VARCHAR(128),
    file_size VARCHAR(128),
	full_path VARCHAR(500),
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci
;

-- Create 2 indices: idx_filename(filename), idx_hash_sha1sum(hash_sha1sum)
-- CREATE INDEX idx_filename ON $mysql_table_name (filename);
-- CREATE INDEX idx_hash_sha1sum ON $mysql_table_name (hash_sha1sum);

-- Bulk sql import
EOF
}

function createMysqlImport_TestPrefix() {
cat <<EOF >> $__DIR__/"$batchGenFilePath"
INSERT INTO $mysql_table_name (job_name, job_timestamp, log_timestamp, filename, hash_sha1sum, file_size, full_path)
VALUES
EOF
}

function createMysqlImport_TestAppendRow() {
cat <<EOF >> $__DIR__/"$batchGenFilePath"
  ("test1234", "202407151749500", "2024-07-10 11:30:58.816", "file.xlsx", "9e6c66fcd1f923d083284f5937104a0d", "512B", "'Output - blank'/file.xlsx"),
EOF
}

function createMysqlImport_TestAppendRowEnding() {
cat <<EOF >> $__DIR__/"$batchGenFilePath"
  ("test1234", "202407151749500", "2024-07-10 11:30:58.816", "file.docx", "f0aaf54f8d18a4e4230867455b40d69c", "512B", "Output/files/file.docx")
;
EOF
}

function createMysqlImport_Prefix() {
cat <<EOF >> $__DIR__/"$batchGenFilePath"
INSERT INTO $mysql_table_name (job_name, job_timestamp, log_timestamp, filename, hash_sha1sum, file_size, full_path)
VALUES
EOF
}

function executeMainTask() {
    local importLogFile="$1"
    if [[ -z "$importLogFile" ]]; then
        importLogFile="$soleLogPath"
    fi

    local sqlDir=$(echo "$batchGenFilePath"| rev| cut -d '/' -f2| rev)
    mkdir -p "$sqlDir"

    createMysqlImport_SqlFile
    createMysqlImport_Prefix

    local this_dir=$(pwd)
    while IFS= read -r line; do
        local get_path=$(echo $line| cut -d ' ' -f3-| tr -d '\r')
        local filename=$(echo "$get_path"| rev | cut -d '/' -f1 | rev)
        local sha1sum=$(sha1sum "$global_workPath/$get_path")
        local sha1sum_trim=$(echo "$sha1sum"| cut -d ' ' -f1)
        local timestamp=$(echo "$line"| awk -F ' ' '{print $1 " " $2}')
        local file_size=$(du -sh "$global_workPath/$get_path"| awk '{print $1"B"}')
        local job_name=$(basename "$importLogFile"| awk -F '_' '{print $1}')

        local job_timestamp=$(basename "$importLogFile"| sed 's/.*_//')
        if [[ "$job_timestamp" == *'.'* ]]; then
            job_timestamp=$(echo "$job_timestamp"| awk -F '.' '{print $1}')
        fi

        echo $sha1sum

        ## INSERT INTO $mysql_table_name (job_name, log_timestamp, filename, hash_sha1sum, file_size, full_path)
        cat <<EOF >> $__DIR__/"$batchGenFilePath"
  ("$job_name", "$job_timestamp", "$timestamp", "$filename", "$sha1sum_trim", "$file_size", "$get_path"),
EOF
    done <<< $(filter file $importLogFile)

    cat <<'EOF' >> $__DIR__/"$batchGenFilePath"
  ("", "", "2024-12-31 00:00:00.123", "", "", "", "")
  ;

use InfoCoreCheckScan;
DELETE FROM ScanResult WHERE `job_name` IS NULL OR `job_name` = '';
EOF

    local resolveSqlSuffix=$(tmp=$(echo "A"))

    batchOpt=$(echo "$keepSqlBatchFile"| awk '{print tolower($0)}')
    if [[ "$batchOpt" == "no" ]]; then
        rm -rf "$__DIR__/$batchGenFilePath"
    fi
}

function createMysqlImport_Test() {
    createMysqlImport_SqlFile
    createMysqlImport_TestPrefix
    createMysqlImport_TestAppendRow
    createMysqlImport_TestAppendRowEnding
    exit 1
}

function checkwithMysqlClient() {
    cmd=$(mysql --version)
    if [[ "$cmd" ]]; then
        echo -e "\nMySQL Client installation check:     [ OK ]\nVersion: $cmd"
        printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
    else
        echo -e "\nFail to perform this script!\nMySQL Client is NOT yet installed on the system, the process has been terminated."
        printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
        echo -e "Operation aborted.\n"
        exit 1
    fi
}

function connectWithMysql() {
    cmd=$(mysql -u"$mysql_usr" -h"$mysql_host" -P$mysql_port -p"$mysql_pwd" -e"select current_user();")
    if [[ "$cmd" ]]; then
        echo -e "\nMySQL Server connection check:       [ OK ]\nBoth connection and credentials for $mysql_usr@$mysql_host:$mysql_port is ready."
        printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
    else
        echo -e "\nFail to connect with specified MySQL Server!\nThe process has been terminated."
        printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
        echo -e "Operation aborted.\n"
        exit 1
    fi
}

function checkLogFileExists() {
    local argv1="$1"
    if [[ ! -z "$argv1" ]]; then
        if [[ -f "$argv1" ]]; then
            echo -e "\nInfoCore specified log file check:       [ OK ]\nLog file path: $argv1\nInfoCore /NFS directory path: $global_workPath"
            printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
        else
            echo -e "\nFail to scan log file at $argv1!\nThe process has been terminated."
            printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
            echo -e "Operation aborted.\n"
        fi
        exit 1
    fi

    if [[ -f "$soleLogPath" ]]; then
        echo -e "\nInfoCore specified log file check:       [ OK ]\nLog file path: $soleLogPath\nInfoCore /NFS directory path: $global_workPath"
        printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
    else
        echo -e "\nFail to scan log file at $soleLogPath!\nThe process has been terminated."
        printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
        echo -e "Operation aborted.\n"
        exit 1
    fi
}

function importSqlToMysql() {
    printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
    cmd=$(mysql -u"$mysql_usr" -h"$mysql_host" -P$mysql_port -p"$mysql_pwd" --default-character-set=utf8 < $__DIR__/"$batchGenFilePath")
}

function bulkImportToMysql_once() {
    checkwithMysqlClient
    connectWithMysql
    checkLogFileExists
    executeMainTask "$soleLogPath"
    importSqlToMysql
}

function bulkImportToMysql_bulk() {
    local targetLogFile="$1"
    checkwithMysqlClient
    connectWithMysql
    checkLogFileExists "$targetLogFile"
    executeMainTask "$targetLogFile"
    importSqlToMysql
}

function lsdir() {
    for logFile in $(ls -F "$global_scanJobLogFolder" | grep -v '/$'); do
        local cmd=$(echo $logFile| tr -d '*')
        echo $global_scanJobLogFolder/$cmd
        nl "$global_scanJobLogFolder/$cmd"
        echo -e "\n\n"
    done
}

function run() {
    for logFile in $(ls -F "$global_scanJobLogFolder" | grep -v '/$'); do
        local fp=$(echo $logFile| tr -d '*')
        if [[ ! "$logFile" == " " ]] || [[ ! "$logFile" == "" ]]; then
            bulkImportToMysql_bulk "$global_scanJobLogFolder/$fp"
        fi
    done
}

main() {
    argv1="$1"

    if [[ -z "$argv1" ]]; then { run; exit 1; } fi
    if [[ "$argv1" == 'run' ]]; then { run; exit 1; } fi
    if [[ "$argv1" == 'test' ]]; then { createDump --force; bulkImportToMysql_once; exit 1; } fi
    if [[ "$argv1" == 'lsdir' ]]; then { lsdir; exit 1; } fi
    if [[ "$argv1" == 'checkDbCli' ]]; then { checkwithMysqlClient; exit 1; } fi
    if [[ "$argv1" == 'checkMysql' ]]; then { connectWithMysql; exit 1; } fi
    if [[ "$argv1" == 'checkLog' ]]; then { checkLogFileExists; exit 1; } fi

    $@
}

main $@