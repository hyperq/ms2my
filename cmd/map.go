package cmd

import "fmt"

type columnType struct {
	TransferType   string
	TransferInsert func(string) string
}

var ms2sqltype = map[string]columnType{
	"bigint": columnType{
		TransferType: "bigint",
	},
	"binary": columnType{
		TransferType: "binary",
	},
	//SQL SERVER的bit类型，对于零，识别为False，非零值识别为True。 MySQL中没有指定的bool类型，一般都使用tinyint来代替
	"bit": columnType{
		TransferType: "tinyint(1)",
	},
	"char": columnType{
		TransferType:   "char",
		TransferInsert: stringfunc,
	},
	"date": columnType{
		TransferType:   "date",
		TransferInsert: stringfunc,
	},
	"datetime": columnType{
		TransferType:   "datetime",
		TransferInsert: timefunc,
	},
	//mssql的保留到微秒(秒后小数点3位)，而mysql仅保留到秒
	"datetime2": columnType{
		TransferType:   "datetime",
		TransferInsert: timefunc,
	},
	//mssql的保留到微秒(秒后小数点7位)，而mysql仅保留到秒
	"datetimeoffset": columnType{
		TransferType:   "datetime",
		TransferInsert: timefunc,
	},
	//mssql的保留时区，这个需要程序自己转换 mssql的保留到微秒(秒后小数点7位)，而mysql仅保留到秒
	"decimal": columnType{
		TransferType: "decimal",
	},
	"float": columnType{
		TransferType: "float",
	},
	"int": columnType{
		TransferType: "int",
	},
	"money": columnType{
		TransferType: "float",
	},
	//默认转换为decimal(19,4)
	"nchar": columnType{
		TransferType:   "char",
		TransferInsert: stringfunc,
	},
	//SQL SERVER转MySQL按正常字节数转就可以
	"ntext": columnType{
		TransferType:   "text",
		TransferInsert: stringfunc,
	},
	"numeric": columnType{
		TransferType: "decimal",
	},
	"nvarchar": columnType{
		TransferType:   "varchar",
		TransferInsert: stringfunc,
	},
	"real": columnType{
		TransferType: "float",
	},
	"smalldatetime": columnType{
		TransferType:   "datetime",
		TransferInsert: timefunc,
	},
	"smallint": columnType{
		TransferType: "smallint",
	},
	"smallmoney": columnType{
		TransferType: "float",
	}, //默认转换为decimal(10,4)
	"text": columnType{
		TransferType:   "text",
		TransferInsert: stringfunc,
	},
	"time": columnType{
		TransferType:   "time",
		TransferInsert: timefunc,
	},
	//注意，mssql的保留到秒后小数点8位，而mysql仅保留到秒
	"timestamp": columnType{
		TransferType: "timestamp",
	},
	"tinyint": columnType{
		TransferType: "tinyint",
	},
	"uniqueidentifier": columnType{
		TransferType:   "char(36)",
		TransferInsert: stringfunc,
	},
	//对应mysql的UUID(),设置为文本类型即可。
	"varbinary": columnType{
		TransferType: "varbinary",
	},
	"varchar": columnType{
		TransferType:   "varchar",
		TransferInsert: stringfunc,
	},
	"xml": columnType{
		TransferType:   "text",
		TransferInsert: stringfunc,
	},
	//mysql不支持xml，修改为text
}

var stringfunc = func(v string) string {
	return fmt.Sprintf("'%s'", v)
}
var timefunc = func(v string) string {
	if len(v) > 18 {
		return fmt.Sprintf("'%s'", v[:19])
	}
	return "''"
}
