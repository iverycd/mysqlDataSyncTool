# mysqlDataSyncTool

([CN](https://github.com/iverycd/mysqlDataSyncTool/blob/master/readme_cn.md))

## Features


Online migration of MySQL databases to target MySQL kernel databases, such as MySQL, PolarDB, Percona Server MySQL, MariaDB, OceanBase, TiDB, GaussDB

- Migrate the entire database table structure and table row data to the target database
- Multi threaded batch migration of table row data
- Data comparison between source and target databases


## Pre-requirement
The running client PC needs to be able to connect to both the source MySQL database and the target database simultaneously

run on Windows,Linux,macOS

## Installation

tar and run 

e.g.

`[root@localhost opt]# tar -zxvf mysqlDataSyncTool-linux64-0.0.1.tar.gz`

## How to use

The following is an example of a Windows platform, with the same command-line parameters as other operating systems

`Note`: Please run this tool in `CMD` on a `Windows` system, or in a directory with read and write permissions on `MacOS` or `Linux`

### 1 Edit yml configuration file

Edit the `example.cfg` file and input the source(src) and target(dest) database information separately

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

pageSize: Number of records per page for pagination query
```
e.g.
pageSize:100000
SELECT t.* FROM (SELECT id FROM test  ORDER BY id LIMIT 0, 100000) temp LEFT JOIN test t ON temp.id = t.id;
```
maxParallel: The maximum number of concurrency that can run goroutine simultaneously

tables: Customized migrated tables and customized query source tables, indented in yml format

exclude: Tables that do not migrate to target database, indented in yml format

batchRowSize: Number of rows in batch insert target table

### 2 Full database migration

Migrate entire database table structure, row data, views, index constraints, and self increasing columns to target database

mysqlDataSyncTool.exe  --config file.yml
```
e.g.
mysqlDataSyncTool.exe --config example.yml

on Linux and MacOS you can run
./mysqlDataSyncTool --config example.yml
```

### 3 View Migration Summary

After the entire database migration is completed, a migration summary will be generated to observe if there are any failed objects. By querying the migration log, the failed objects can be analyzed

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

### 4 Compare Source and Target database

After migration finish you can compare source table and target database table rows,displayed failed table only

`mysqlDataSyncTool.exe --config your_file.yml compareDb`

```
e.g.
mysqlDataSyncTool.exe --config example.yml compareDb

on Linux and MacOS you can run
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







## Other migration modes


#### 1 Full database migration

Migrate entire database table structure, row data, views, index constraints, and self increasing columns to target database

mysqlDataSyncTool.exe  --config file.yml

```
e.g.
mysqlDataSyncTool.exe --config example.yml
```

#### 2 Custom SQL Query Migration

only migrate some tables not entire database, and migrate the table structure and table data to the target database according to the custom query statement in file.yml

mysqlDataSyncTool.exe  --config file.yml -s

```
e.g.
mysqlDataSyncTool.exe  --config example.yml -s
```

#### 3 Migrate all table structures in the entire database

Create all table structure(only table metadata not row data) to  target database

mysqlDataSyncTool.exe  --config file.yml createTable -t

```
e.g.
mysqlDataSyncTool.exe  --config example.yml createTable -t
```

#### 4 Migrate the table structure of custom tables

Read custom tables from yml file and create target table 

mysqlDataSyncTool.exe  --config file.yml createTable -s -t

```
e.g.
mysqlDataSyncTool.exe  --config example.yml createTable -s -t
```