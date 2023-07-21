package cmd

import (
	"fmt"
	"strconv"
	"time"
)

var tabRet []string
var tableCount int
var failedCount int

type Database interface {
	// TableCreate (logDir string, tableMap map[string][]string) (result []string) 单线程
	TableCreate(logDir string, tblName string, ch chan struct{})
}

type Table struct {
	columnName             string
	dataType               string
	characterMaximumLength string
	isNullable             string
	columnDefault          string
	numericPrecision       string
	numericScale           string
	datetimePrecision      string
	columnKey              string
	columnComment          string
	ordinalPosition        int
	destType               string
	destNullable           string
	destDefault            string
	autoIncrement          int
	destSeqSql             string
	destDefaultSeq         string
	dropSeqSql             string
	destIdxSql             string
	viewSql                string
}

func (tb *Table) TableCreate(logDir string, tblName string, ch chan struct{}) {
	defer wg2.Done()
	tableCount += 1
	// 使用goroutine并发的创建多个表
	var createTblName string
	var CreateSql string
	// 查询MySQL表结构
	sql := fmt.Sprintf("show create table `%s`", tblName)
	//fmt.Println(sql)
	err := srcDb.QueryRow(sql).Scan(&createTblName, &CreateSql)
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("%v Table total %s create table %s", time.Now().Format("2006-01-02 15:04:05.000000"), strconv.Itoa(tableCount), tblName))
	txn, err := destDb.Begin()
	if err != nil {
		log.Error(err)
	}
	// 先禁用外键约束检查，避免表创建失败
	stmt, err := txn.Prepare("SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		log.Error("table ", tblName, " SET FOREIGN_KEY_CHECKS=0 failed  ", err)
		LogError(logDir, "tableCreateFailed", CreateSql, err)
		failedCount += 1
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Error(err)
	}
	// 创建前先删除目标表
	dropDestTbl := "drop table if exists " + fmt.Sprintf("`") + tblName + fmt.Sprintf("`") + " cascade"
	_, err = txn.Exec(dropDestTbl)
	if err != nil {
		log.Error(err)
	}
	// 以下创建表操作
	stmt, err = txn.Prepare(CreateSql)
	if err != nil {
		log.Error("table ", tblName, " create failed  ", err)
		LogError(logDir, "tableCreateFailed", CreateSql, err)
		failedCount += 1
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Error(err)
	}
	err = stmt.Close()
	if err != nil {
		log.Error(err)
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
	}
	<-ch
}
