# ms2my
a tool to transfer sql server database to mysql

## Install
```bash
go get github.com/hyperq/ms2my
```

## How to use
```bash
md2my -i 127.0.0.1 -d dbname -p 123456 -t user
```
you can use md2my -h check all flags

```bash
ms2my [flags]

Flags:
  -d, --dbname string      dbname
  -h, --help               help for ms2my
  -i, --ip string          ip (default "127.0.0.1")
  -p, --password string    mssql password
      --port int           mssql port (default 1433)
  -t, --tablename string   table names,you yan use , join mult
  -u, --username string    mssql username (default "sa")
```
