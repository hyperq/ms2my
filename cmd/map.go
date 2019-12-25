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

var keyWorks = map[string]bool{
	"ADD": true, "ALL": true, "ALTER": true,
	"ANALYZE": true, "AND": true, "AS": true,
	"ASC": true, "ASENSITIVE": true, "BEFORE": true,
	"BETWEEN": true, "BIGINT": true, "BINARY": true,
	"BLOB": true, "BOTH": true, "BY": true,
	"CALL": true, "CASCADE": true, "CASE": true,
	"CHANGE": true, "CHAR": true, "CHARACTER": true,
	"CHECK": true, "COLLATE": true, "COLUMN": true,
	"CONDITION": true, "CONNECTION": true, "CONSTRAINT": true,
	"CONTINUE": true, "CONVERT": true, "CREATE": true,
	"CROSS": true, "CURRENT_DATE": true, "CURRENT_TIME": true,
	"CURRENT_TIMESTAMP": true, "CURRENT_USER": true, "CURSOR": true,
	"DATABASE": true, "DATABASES": true, "DAY_HOUR": true,
	"DAY_MICROSECOND": true, "DAY_MINUTE": true, "DAY_SECOND": true,
	"DEC": true, "DECIMAL": true, "DECLARE": true,
	"DEFAULT": true, "DELAYED": true, "DELETE": true,
	"DESC": true, "DESCRIBE": true, "DETERMINISTIC": true,
	"DISTINCT": true, "DISTINCTROW": true, "DIV": true,
	"DOUBLE": true, "DROP": true, "DUAL": true,
	"EACH": true, "ELSE": true, "ELSEIF": true,
	"ENCLOSED": true, "ESCAPED": true, "EXISTS": true,
	"EXIT": true, "EXPLAIN": true, "FALSE": true,
	"FETCH": true, "FLOAT": true, "FLOAT4": true,
	"FLOAT8": true, "FOR": true, "FORCE": true,
	"FOREIGN": true, "FROM": true, "FULLTEXT": true,
	"GOTO": true, "GRANT": true, "GROUP": true,
	"HAVING": true, "HIGH_PRIORITY": true, "HOUR_MICROSECOND": true,
	"HOUR_MINUTE": true, "HOUR_SECOND": true, "IF": true,
	"IGNORE": true, "IN": true, "INDEX": true,
	"INFILE": true, "INNER": true, "INOUT": true,
	"INSENSITIVE": true, "INSERT": true, "INT": true,
	"INT1": true, "INT2": true, "INT3": true,
	"INT4": true, "INT8": true, "INTEGER": true,
	"INTERVAL": true, "INTO": true, "IS": true,
	"ITERATE": true, "JOIN": true, "KEY": true,
	"KEYS": true, "KILL": true, "LABEL": true,
	"LEADING": true, "LEAVE": true, "LEFT": true,
	"LIKE": true, "LIMIT": true, "LINEAR": true,
	"LINES": true, "LOAD": true, "LOCALTIME": true,
	"LOCALTIMESTAMP": true, "LOCK": true, "LONG": true,
	"LONGBLOB": true, "LONGTEXT": true, "LOOP": true,
	"LOW_PRIORITY": true, "MATCH": true, "MEDIUMBLOB": true,
	"MEDIUMINT": true, "MEDIUMTEXT": true, "MIDDLEINT": true,
	"MINUTE_MICROSECOND": true, "MINUTE_SECOND": true, "MOD": true,
	"MODIFIES": true, "NATURAL": true, "NOT": true,
	"NO_WRITE_TO_BINLOG": true, "NULL": true, "NUMERIC": true,
	"ON": true, "OPTIMIZE": true, "OPTION": true,
	"OPTIONALLY": true, "OR": true, "ORDER": true,
	"OUT": true, "OUTER": true, "OUTFILE": true,
	"PRECISION": true, "PRIMARY": true, "PROCEDURE": true,
	"PURGE": true, "RAID0": true, "RANGE": true,
	"READ": true, "READS": true, "REAL": true,
	"REFERENCES": true, "REGEXP": true, "RELEASE": true,
	"RENAME": true, "REPEAT": true, "REPLACE": true,
	"REQUIRE": true, "RESTRICT": true, "RETURN": true,
	"REVOKE": true, "RIGHT": true, "RLIKE": true,
	"SCHEMA": true, "SCHEMAS": true, "SECOND_MICROSECOND": true,
	"SELECT": true, "SENSITIVE": true, "SEPARATOR": true,
	"SET": true, "SHOW": true, "SMALLINT": true,
	"SPATIAL": true, "SPECIFIC": true, "SQL": true,
	"SQLEXCEPTION": true, "SQLSTATE": true, "SQLWARNING": true,
	"SQL_BIG_RESULT": true, "SQL_CALC_FOUND_ROWS": true, "SQL_SMALL_RESULT": true,
	"SSL": true, "STARTING": true, "STRAIGHT_JOIN": true,
	"TABLE": true, "TERMINATED": true, "THEN": true,
	"TINYBLOB": true, "TINYINT": true, "TINYTEXT": true,
	"TO": true, "TRAILING": true, "TRIGGER": true,
	"TRUE": true, "UNDO": true, "UNION": true,
	"UNIQUE": true, "UNLOCK": true, "UNSIGNED": true,
	"UPDATE": true, "USAGE": true, "USE": true,
	"USING": true, "UTC_DATE": true, "UTC_TIME": true,
	"UTC_TIMESTAMP": true, "VALUES": true, "VARBINARY": true,
	"VARCHAR": true, "VARCHARACTER": true, "VARYING": true,
	"WHEN": true, "WHERE": true, "WHILE": true,
	"WITH": true, "WRITE": true, "X509": true,
	"XOR": true, "YEAR_MONTH": true, "ZEROFILL": true,
}
