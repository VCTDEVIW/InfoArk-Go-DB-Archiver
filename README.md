# InfoArk-Go-DB-Archiver 日誌歸檔入庫工具
- ✅ Go(Golang) 版本至少為 1.18 或以上
- ✅ 數據庫：Microsoft SQL Server (不限版本或 2012 以上)
- ✅ Windows OS 必須為 64位 及 Windows 7 或以上
- ✅ Linux OS 必須使用 64位 及內核版本為 3.10 或以上一切發行版 (CentOS 7+, RHEL 7+, Debian 9+ 等等)


## 到 Google 官方下載 Go(Golang) 環境：
https://down.golang.org/dl


## 環境準備好後，請下載項目依賴庫：
```
go get -d github.com/denisenkom/go-mssqldb

go get -d github.com/orcaman/concurrent-map
```

## 【正式使用項目】：
- 初始化項目：
```
go mod init infoshare
```
- 開發測試：
```
go run main.go
```

- 編譯：
```
go build main.go
```
