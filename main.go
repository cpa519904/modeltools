package main

import (
	"github.com/cpa519904/modeltools/conf"
	"github.com/cpa519904/modeltools/dbtools"
	"github.com/cpa519904/modeltools/generate"
)

func main() {
	// model保存路径
	//conf.ModelPath=""
	// 是否覆盖已存在model
	//conf.ModelReplace=""
	// 数据库驱动
	conf.MasterDbConfig.Host = "127.0.0.1"
	conf.MasterDbConfig.Port = "3306"
	conf.MasterDbConfig.User = "root"
	conf.MasterDbConfig.Pwd = "long"
	conf.MasterDbConfig.DbName = "mvideo"
	//初始化数据库
	dbtools.Init()
	//generate.Genertate() //生成所有表信息
	generate.Genertate("admin_info", "video_info") //生成指定表信息，可变参数可传入多个表名
}
