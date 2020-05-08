package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/didi/gendry/scanner"
)

var mssql *sql.DB

type columns struct {
	TableName           string `ddb:"tablename"`
	TableComment        string `ddb:"tablecomment"`
	ColumnName          string `ddb:"columnname"`
	PK                  int    `ddb:"pk"`
	ColumnType          string `ddb:"columntype"`
	ColumnLength        int    `ddb:"columnlength"`
	ColumnDecimalLength int    `ddb:"columndecimallength"`
	Null                int    `ddb:"nulls"`
	ColumnDefault       string `ddb:"columndefault"`
	ColumnComment       string `ddb:"columncomment"`
}

func generateCreate(table string) (creates string, err error) {
	rows, err := mssql.Query(fmt.Sprintf(`
	SELECT
      tablename=d.name,
      tablecomment=isnull(f.value,''),
      columnname=a.name,
      pk=case when exists(SELECT 1 FROM sysobjects where xtype='PK' and name in (
         SELECT name FROM sysindexes WHERE indid in(
			SELECT indid FROM sysindexkeys WHERE id=a.id AND colid=a.colid
         ))) then 1 else 0 end,
      columntype=b.name,
      columnlength=COLUMNPROPERTY(a.id,a.name,'PRECISION'),
  	  columndecimallength=isnull(COLUMNPROPERTY(a.id,a.name,'Scale'),0),
  	  nulls=case when a.isnullable=1 then 1 else 0 end,
  	  columndefault=isnull(e.text,''),
      columncomment=isnull(g.[value],'')
 	FROM syscolumns a
      left join systypes b on a.xusertype=b.xusertype
	  inner join sysobjects d on a.id=d.id and d.xtype='U' and d.name<>'dtproperties'
      left join syscomments e on a.cdefault=e.id
      left join sys.extended_properties g on a.id=g.major_id and a.colid=g.minor_id
      left join sys.extended_properties f on d.id=f.major_id and f.minor_id=0
    where d.name='%s' 
    order by a.id,a.colorder
	`, table))
	if err != nil {
		return
	}
	var cs []columns
	err = scanner.ScanClose(rows, &cs)
	if err != nil {
		return
	}
	create := "-- ----------------------------\n-- Table structure for %s\n-- ----------------------------\nDROP TABLE IF EXISTS `%s" +
		"`;\nCREATE TABLE `%s`(\n%s\n%s\n);\n"
	var columns []string
	var pks string
	for _, v := range cs {
		if v.PK == 1 && pks == "" {
			pks = fmt.Sprintf("\tprimary key(%s)", v.ColumnName)
		}
		ct, ok3 := ms2sqltype[v.ColumnType]
		if !ok3 {
			//err = errors.New("暂不支持" + v.ColumnType)
			//return
			continue
		}
		ts := ct.TransferType
		if v.ColumnType == "timestamp" {
			ts = "timestamp(6)"
		} else {
			if ct.TransferType != "text" && ct.TransferType != "datetime" && v.ColumnType != "uniqueidentifier" && v.ColumnType != "bit" {
				if v.ColumnLength > 0 {
					ts = fmt.Sprintf("%s(%d)", ct.TransferType, v.ColumnLength)
				}
				if v.ColumnDecimalLength > 0 {
					ts = fmt.Sprintf("%s(%d,%d)", ct.TransferType, v.ColumnLength, v.ColumnDecimalLength)
				}
			}
		}
		col := fmt.Sprintf("\t`%s` %s", checkik(v.ColumnName), ts)
		if v.Null == 1 {
			col = col + " NOT NULL"
		}
		if v.ColumnDefault != "" && ct.TransferType != "text" {
			col = col + " DEFAULT " + strings.Trim(strings.Trim(v.ColumnDefault, "("), ")")
		}
		if v.ColumnComment != "" {
			col = col + " COMMENT '" + v.ColumnComment + "'"
		}
		col += ",\n"
		columns = append(columns, col)
	}
	creates = fmt.Sprintf(create, strings.ToLower(table), strings.ToLower(table), strings.ToLower(table), strings.Join(columns, ""), pks)
	return
}
func reverse(b []uint8) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

type columntypes struct {
	Name     string
	Typename string
}

func generateInsert(table string) (inserts string, err error) {
	rows, err := mssql.Query(fmt.Sprintf("select * from [%s]", table))
	if err != nil {

		return
	}
	columntype, err := rows.ColumnTypes()
	if err != nil {

		return
	}
	columnlength := len(columntype)
	var insertsqls []string
	values := make([]interface{}, columnlength)
	cts := make([]columntypes, columnlength)
	for i := 0; i < columnlength; i++ {
		values[i] = new(interface{})
		cts[i].Name = checkik(columntype[i].Name())
		cts[i].Typename = columntype[i].DatabaseTypeName()
	}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {

			return
		}
		var key []string
		var val []string
		for k, v := range cts {
			DatabaseTypeName := v.Typename
			if DatabaseTypeName == "BINARY" || DatabaseTypeName == "VARBINARY" {
				continue
			}
			value := *(values[k].(*interface{}))
			rv, ok := value.([]uint8)
			var vs = ""
			if !ok {
				if value == nil {
					vs = "null"
				} else {
					vs = fmt.Sprint(value)
				}
			} else {
				if DatabaseTypeName == "UNIQUEIDENTIFIER" {
					reverse(rv[0:4])
					reverse(rv[4:6])
					reverse(rv[6:8])
					vs = fmt.Sprintf("%X-%X-%X-%X-%X", rv[0:4], rv[4:6], rv[6:8], rv[8:10], rv[10:])
				} else {
					vs = string(rv)
				}
			}
			typename := strings.ToLower(DatabaseTypeName)
			ct, ok3 := ms2sqltype[typename]
			if !ok3 {
				//err = errors.New("暂不支持" + typename)
				continue
			}
			if ct.TransferInsert != nil {
				vs = ct.TransferInsert(vs)
			}
			key = append(key, v.Name)
			val = append(val, vs)
		}
		is := fmt.Sprintf("INSERT INTO %s (%s) values (%s);", strings.ToLower(table), strings.Join(key, ","), strings.Join(val, ","))
		insertsqls = append(insertsqls, is)
	}
	inserts = "-- ----------------------------\n-- Records of bigbox \n-- ----------------------------\nBEGIN;\n" + strings.Join(
		insertsqls, "\n") + "\nCOMMIT;\n"
	return
}

func checkik(key string) string {
	lkey := strings.ToUpper(key)
	if _, ok := keyWorks[lkey]; ok {
		return key + "s"
	}
	return key
}
