package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

var mssql *sql.DB

func generateCreate(table string) (creates string, err error) {
	rows, err := mssql.Query("select Top 1 * from " + table)
	if err != nil {
		return
	}
	columntype, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	create := "-- ----------------------------\n-- Table structure for %s\n-- ----------------------------\nDROP TABLE IF EXISTS `%s" +
		"`;\nCREATE TABLE `%s`(\n%s\n);\n"
	var columns []string
	for _, v := range columntype {
		name := v.Name()
		precision, scale, ok := v.DecimalSize()
		length, ok2 := v.Length()
		typename := strings.ToLower(v.DatabaseTypeName())
		ct, ok3 := ms2sqltype[typename]
		if !ok3 {
			err = errors.New("暂不支持" + typename)
			return
		}
		if ct.TransferType != "text" {
			if ok2 {
				ct.TransferType = fmt.Sprintf("%s(%d)", ct.TransferType, length)
			}
			if ok {
				ct.TransferType = fmt.Sprintf("%s(%d,%d)", ct.TransferType, precision, scale)
			}
		}
		col := fmt.Sprintf("\t`%s` %s,\n", name, ct.TransferType)
		columns = append(columns, col)
	}
	creates = fmt.Sprintf(create, table, table, table, strings.Trim(strings.Join(columns, ""), ",\n"))
	return
}
func reverse(b []uint8) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func generateInsert(table string) (inserts string, err error) {
	rows, err := mssql.Query("select * from " + table)
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
	for i := 0; i < columnlength; i++ {
		values[i] = new(interface{})
	}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {

			return
		}
		var key []string
		var val []string
		for k, v := range columntype {
			DatabaseTypeName := v.DatabaseTypeName()
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
				err = errors.New("暂不支持" + typename)
				return
			}
			if ct.TransferInsert != nil {
				vs = ct.TransferInsert(vs)
			}
			key = append(key, v.Name())
			val = append(val, vs)
		}
		is := fmt.Sprintf("INSERT INTO %s (%s) values (%s);", table, strings.Join(key, ","), strings.Join(val, ","))
		insertsqls = append(insertsqls, is)
	}
	inserts = "-- ----------------------------\n-- Records of bigbox \n-- ----------------------------\nBEGIN;\n" + strings.Join(
		insertsqls, "\n") + "\nCOMMIT;\n"
	return
}
