-- 建庫
create database [InfoShareLogScan]
go

-- 檢查表是否存在
IF EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'ScanResult' AND TABLE_SCHEMA = 'dbo') SELECT 1 ELSE SELECT 0

-- 建表
use [InfoShareLogScan];
CREATE TABLE ScanResult (
    rid INT IDENTITY(1,1) PRIMARY KEY,
    job_name NVARCHAR(255) COLLATE Chinese_PRC_CI_AS,
    job_timestamp DATETIME,
    log_timestamp DATETIME,
    filename NVARCHAR(255) COLLATE Chinese_PRC_CI_AS,
    hash_sha1sum NVARCHAR(96),
    file_size NVARCHAR(255),
    full_path NVARCHAR(500) COLLATE Chinese_PRC_CI_AS,
    inserted_at DATETIME DEFAULT GETDATE()
)
go
--
