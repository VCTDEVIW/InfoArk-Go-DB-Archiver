package vastcom_infoshare

import (
	. "fmt"
	"database/sql"
    _ "github.com/denisenkom/go-mssqldb"
	_ "os"
    _"log"
)

func init() {
	Load_DB_Driver = "msSQL; Microsoft SQL Server"
	Println("")
}

func ConnectToMssql() bool {
	server := LoadConfig.DB_Host
    port := LoadConfig.DB_Port // Default SQL Server port is 1433
    user := LoadConfig.DB_User
    password := LoadConfig.DB_PW
    //database := LoadConfig.Database_Name

    // Connect to the database
    connString := Sprintf("server=%s,%d;user id=%s;password=%s;", server, port, user, password)
    db, err := sql.Open("mssql", connString)
    if err != nil {
        Println("Failed to connect to the MS SQL Server database server:\n", err)
        return false
    }
	defer func() {
		db.Close()
	}()

    // Ping the database to verify the connection
    err = db.Ping()
    if err != nil {
        Println("Failed to ping the target MS SQL Server database server:\n", err)
        return false
    }

    Println("Connected to the MS SQL Server successfully!\n")
	return true
}

func CheckDatabase() bool {
	server := LoadConfig.DB_Host
    port := LoadConfig.DB_Port // Default SQL Server port is 1433
    user := LoadConfig.DB_User
    password := LoadConfig.DB_PW
    database := LoadConfig.Database_Name

    // Connect to the database
    connString := Sprintf("server=%s,%d;user id=%s;password=%s;database=%s", server, port, user, password, database)
    db, err := sql.Open("mssql", connString)
    if err != nil {
        Println("Failed to open the database object:\n", err)
        return false
    }
	defer func() {
		db.Close()
	}()

    // Ping the database to verify the connection
    err = db.Ping()
    if err != nil {
        Println("Failed to validate the database object:\n", err)
        return false
    }

    Printf("Validate the database [%s] successfully!\n", LoadConfig.Database_Name)
	return true
}

func CheckTable() bool {
	server := LoadConfig.DB_Host
    port := LoadConfig.DB_Port // Default SQL Server port is 1433
    user := LoadConfig.DB_User
    password := LoadConfig.DB_PW
    database := LoadConfig.Database_Name
    table := LoadConfig.Table_Name

    // Connect to the database
    connString := Sprintf("server=%s,%d;user id=%s;password=%s;database=%s", server, port, user, password, database)
    db, err := sql.Open("mssql", connString)
    if err != nil {
        Println("Failed to open the database object:\n", err)
        return false
    }
	defer func() {
		db.Close()
	}()

    var count int
	var sqlCmd string
	sqlCmd = Sprintf("IF EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = '%s' AND TABLE_SCHEMA = 'dbo') SELECT 1 ELSE SELECT 0", table)
	err = db.QueryRow(sqlCmd).Scan(&count)
	if err != nil {
		//log.Fatal(err)
        Println(err)
	}

	if count == 0 {
        sqlCmd = Sprintf(`
		use [%s];
		CREATE TABLE %s (
			rid INT IDENTITY(1,1) PRIMARY KEY,
			job_name NVARCHAR(255) COLLATE Chinese_PRC_CI_AS,
			job_timestamp NVARCHAR(32),
			log_timestamp NVARCHAR(32),
			filename NVARCHAR(255) COLLATE Chinese_PRC_CI_AS,
			hash_sha1sum NVARCHAR(96),
			file_size NVARCHAR(255),
			full_path NVARCHAR(500) COLLATE Chinese_PRC_CI_AS,
			inserted_at DATETIME DEFAULT GETDATE()
		)
		`, database, table)
		// Create the table
		_, err = db.Exec(sqlCmd)
		if err != nil {
			//log.Fatal(err)
            Println(err)
		}
		Printf("Table [%s].%s created successfully.", database, table)
        return true
	} else {
		Printf("Table [%s].%s already exists.", database, table)
        return true
	}
}

func BulkSqlSubmit(_import string) {
	server := LoadConfig.DB_Host
    port := LoadConfig.DB_Port // Default SQL Server port is 1433
    user := LoadConfig.DB_User
    password := LoadConfig.DB_PW
    database := LoadConfig.Database_Name

    // Connect to the database
    connString := Sprintf("server=%s,%d;user id=%s;password=%s;database=%s", server, port, user, password, database)
    db, err := sql.Open("mssql", connString)
    if err != nil {
        Println("Failed to open the database object:\n", err)
    }
	defer func() {
		db.Close()
	}()

    _, err = db.Exec(_import)
    if err != nil {
        //log.Fatal(err)
        Println(err)
    }
    Println("Multiple records inserted successfully.")
}

