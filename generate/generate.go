package generate

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"io"
	"github.com/cpa519904/modeltools/conf"
	"github.com/cpa519904/modeltools/dbtools"
	"os"
	"strings"
)

// 表名称 与model相对应 多个,隔开
var TableGenerateNames = ""

// 表 对应 model 名称 与表名相对应 多个,隔开
var ModelGenerateNames = ""

func Genertate(tableNames ...string) {
	tableNamesStr := ""
	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}
	if len(tableNamesStr) == 0 {
		if len(TableGenerateNames) > 0 {
			tableNames := strings.Split(TableGenerateNames, ",")
			for _, name := range tableNames {
				if tableNamesStr != "" {
					tableNamesStr += ","
				}
				tableNamesStr += "'" + name + "'"
			}
		}

	}
	tables := getTables(tableNamesStr) //生成所有表信息
	//tables := getTables("admin_info","video_info") //生成指定表信息，可变参数可传入过个表名
	tableMdel := strings.Split(ModelGenerateNames, ",")
	for i, table := range tables {
		fields := getFields(table.Name)
		if len(tableMdel) > 0 {
			//index:=(i+1)
			tableMdelBase := tableMdel[i]
			table.ModelName = tableMdelBase
		}
		generateModel(table, fields)
	}
}

//获取表信息
func getTables(tableNames string) []Table {
	db := dbtools.GetMysqlDb()
	var tables []Table
	if tableNames == "" {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + conf.MasterDbConfig.DbName + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + tableNames + ") AND table_schema='" + conf.MasterDbConfig.DbName + "';").Find(&tables)
	}
	return tables
}

//获取所有字段信息
func getFields(tableName string) []Field {
	db := dbtools.GetMysqlDb()
	var fields []Field
	db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

//生成Model
func generateModel(table Table, fields []Field) {
	content := "package models\n\n"
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + table.Comment + "\n"
	}
	if len(table.ModelName) == 0 {
		content += "type " + generator.CamelCase(table.Name) + " struct {\n"
	} else {
		content += "type " + table.ModelName + " struct {\n"
	}

	//生成字段
	for _, field := range fields {
		fieldName := generator.CamelCase(field.Field)
		fieldJson := getFieldJson(field)
		fieldType := getFiledType(field)
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " `" + fieldJson + "` " + fieldComment + "\n"
	}
	content += "}"

	filename := conf.ModelPath + generator.CamelCase(table.Name) + ".go"
	if len(table.ModelName) > 0 {
		filename = conf.ModelPath + table.ModelName + ".go"
	}
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		if !conf.ModelReplace {
			fmt.Println(generator.CamelCase(table.Name) + " 已存在，需删除才能重新生成...")
			return
		}
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		if len(table.ModelName) > 0 {
			fmt.Println(table.ModelName + " 已生成...")
		} else {
			fmt.Println(generator.CamelCase(table.Name) + " 已生成...")
		}

	}
}

//获取字段类型
func getFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")

	switch typeArr[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

//获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + field.Field + `"`
}

//获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
