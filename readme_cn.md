# mysqlDataSyncTool

## 一、工具特性以及环境要求
### 1.1 功能特性

在线迁移MySQL数据库到目标MySQL内核的数据库，如MySQL,PolarDB,Percona Server MySQL,MariaDB,OceanBase,TiDB,GaussDB

- 迁移全库表结构以及表行数据到目标数据库
- 多线程批量方式迁移表行数据
- 数据比对源库以及目标库

### 1.2 环境要求
在运行的客户端PC需要同时能连通源端MySQL数据库以及目标数据库

支持Windows、Linux、MacOS

### 1.3 如何安装

解压之后即可运行此工具

若在Linux环境下请使用tar解压，例如：


`[root@localhost opt]# tar -zxvf mysqlDataSyncTool-linux64-0.0.1.tar.gz`

## 二、使用方法

以下为Windows平台示例，其余操作系统命令行参数一样

`注意:`在`Windows`系统请在`CMD`运行此工具，如果是在`MacOS`或者`Linux`系统，请在有读写权限的目录运行

### 2.1 编辑yml配置文件

编辑`example.cfg`文件，分别输入源库跟目标数据库信息

```yaml
src:
  host: 192.168.1.3
  port: 3306
  database: test
  username: root
  password: 11111
dest:
  host: 192.168.1.200
  port: 5432
  database: test
  username: test
  password: 11111
pageSize: 100000
maxParallel: 30
batchRowSize: 1000
tables:
  test1:
    - select * from test1
  test2:
    - select * from test2
exclude:
  operalog1
  operalog2
  operalog3
```

pageSize: 分页查询每页的记录数
```
e.g.
pageSize:100000
SELECT t.* FROM (SELECT id FROM test  ORDER BY id LIMIT 0, 100000) temp LEFT JOIN test t ON temp.id = t.id;
```
maxParallel: 最大能同时运行goroutine的并发数

tables: 自定义迁移的表以及自定义查询源表，按yml格式缩进

exclude: 不需要迁移的表，按yml格式缩进

batchRowSize: 批量insert目标表的行数

### 2.2 全库迁移

迁移全库表结构、行数据，索引约束、自增列等对象

mysqlDataSyncTool.exe  --config 配置文件
```
示例
mysqlDataSyncTool.exe --config example.yml

如果是Linux或者macOS请在终端运行
./mysqlDataSyncTool --config example.yml
```

### 2.3 查看迁移摘要

全库迁移完成之后会生成迁移摘要，观察下是否有失败的对象，通过查询迁移日志可对迁移失败的对象进行分析

```bash
+-------------------------+---------------------+-------------+----------+
|        SourceDb         |       DestDb        | MaxParallel | PageSize |
+-------------------------+---------------------+-------------+----------+
| 192.168.149.37-sourcedb | 192.168.149.33-test |     30      |  100000  |
+-------------------------+---------------------+-------------+----------+

+-----------+----------------------------+----------------------------+-------------+--------------+
|Object     |         BeginTime          |          EndTime           |FailedTotal  |ElapsedTime   |
+-----------+----------------------------+----------------------------+-------------+--------------+
|Table      | 2023-07-21 17:12:51.680525 | 2023-07-21 17:12:52.477100 |0            |796.579837ms  |
|TableData  | 2023-07-21 17:12:52.477166 | 2023-07-21 17:12:59.704021 |0            |7.226889553s  |
+-----------+----------------------------+----------------------------+-------------+--------------+


Table Create finish elapsed time  5.0256021s

```

### 2.4 比对数据库

迁移完之后比对源库和目标库，查看是否有迁移数据失败的表

`windows使用:mysqlDataSyncTool.exe --config your_file.yml compareDb`

```
e.g.
mysqlDataSyncTool.exe --config example.yml compareDb

在Linux，MacOS使用示例如下
./mysqlDataSyncTool --config example.yml compareDb
```

```bash
Table Compare Result (Only Not Ok Displayed)
+-----------------------+------------+----------+-------------+------+
|Table                  |SourceRows  |DestRows  |DestIsExist  |isOk  |
+-----------------------+------------+----------+-------------+------+
|abc_testinfo           |7458        |0         |YES          |NO    |
|log1_qweharddiskweqaz  |0           |0         |NO           |NO    |
|abcdef_jkiu_button     |4           |0         |YES          |NO    |
|abcdrf_yuio            |5           |0         |YES          |NO    |
|zzz_ss_idcard          |56639       |0         |YES          |NO    |
|asdxz_uiop             |290497      |190497    |YES          |NO    |
|abcd_info              |1052258     |700000    |YES          |NO    |
+-----------------------+------------+----------+-------------+------+ 
INFO[0040] Table Compare finish elapsed time 11.307881434s 
```




## 三、其他迁移模式

### 1 全库迁移

迁移全库表结构、行数据，视图、索引约束、自增列等对象

mysqlDataSyncTool.exe  --config 配置文件

```
示例
mysqlDataSyncTool.exe --config example.yml
```

### 2 自定义SQL查询迁移

不迁移全库数据，只迁移部分表，根据配置文件中自定义查询语句迁移表结构和表数据到目标库

mysqlDataSyncTool.exe  --config 配置文件 -s

```
示例
mysqlDataSyncTool.exe  --config example.yml -s
```

### 3 迁移全库所有表结构

仅在目标库创建所有表的表结构

mysqlDataSyncTool.exe  --config 配置文件 createTable -t

```
示例
mysqlDataSyncTool.exe  --config example.yml createTable -t
```

### 4 迁移自定义表的表结构

仅在目标库创建自定义的表

mysqlDataSyncTool.exe  --config 配置文件 createTable -s -t

```
示例
mysqlDataSyncTool.exe  --config example.yml createTable -s -t
```
