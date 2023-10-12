package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	mysqlDBs map[string]*gorm.DB
)

type User struct {
	UserId      uint64 `json:"user_id" xorm:"user_id" gorm:"primaryKey;column:user_id"`
	Name        string `json:"user_name" xorm:"user_name" gorm:"column:user_name"`
	Age         int32  `json:"age" xorm:"age" gorm:"column:age"`
	CountryCode string `json:"country_code" xorm:"country_code" gorm:"column:country_code"`
	Number      string `json:"number" xorm:"number" gorm:"column:number"`
	Remark      string `json:"remark" xorm:"remark" gorm:"column:remark"`
}

// create table sql
// CREATE TABLE `test_user` (
// 	`user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
// 	`user_name` varchar(100) NOT NULL DEFAULT '',
// 	`age` int NOT NULL DEFAULT '0',
// 	`country_code` varchar(200) NOT NULL DEFAULT '',
// 	`number` varchar(200) NOT NULL DEFAULT '',
// 	`remark` varchar(200) NOT NULL DEFAULT '',
// 	PRIMARY KEY (`user_id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC

func (u User) TableName() string {
	return "test_user"
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	mysqlDBs = make(map[string]*gorm.DB)

	dialector := mysql.New(mysql.Config{
		// DSN:                       "root:123456@tcp(10.10.21.65:3307)/test?charset=utf8mb4&parseTime=True&loc=Local",
		//DSN:                       "root:gH5wRmTxAmAn@tcp(101.132.227.177:12345)/testdb",
		DSN:                       "root:4B0BlT0n5Ini@tcp(rm-uf635virzqqr2329u6o.mysql.rds.aliyuncs.com:3306)/test_db",
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Printf("MySQL连接异常:%s", err)
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxIdleTime(time.Minute * 20)

	mysqlDBs["test"] = db

	fmt.Println("Mysql连接成功")
	var result *gorm.DB

	// ------------- create -------------
	// user := User{
	// 	Name:        "Jinzhu",
	// 	Age:         18,
	// 	CountryCode: "86",
	// 	Number:      "18812345678",
	// 	Remark:      "remakr",
	// }

	// result = db.Create(&user) // 通过数据的指针来创建
	// if result.Error != nil {
	// 	fmt.Println("create failed")
	// } else {
	// 	fmt.Println("create succ")
	// 	fmt.Printf("user: %+v\n", user)
	// 	fmt.Printf("result: %+v\n", result)
	// }

	// user.UserId = 0
	// result = db.Create(&user) // 通过数据的指针来创建
	// user.UserId = 0
	// result = db.Create(&user) // 通过数据的指针来创建
	// user.UserId = 0
	// result = db.Create(&user) // 通过数据的指针来创建
	// user.UserId = 0
	// result = db.Create(&user) // 通过数据的指针来创建
	// user.UserId = 0
	// result = db.Create(&user) // 通过数据的指针来创建
	// return

	// ------------- find -------------
	firstUser := &User{}
	result = db.First(firstUser, "user_id = ?", "10002")
	if result.Error != nil {
		fmt.Println("First failed")
	} else {
		fmt.Printf("firstUser: %+v\n", firstUser)
	}

	// not recommend
	// firstUser2 := map[string]interface{}{}
	// db.Model(&User{}).First(&firstUser2, "user_id = ?", "10003")
	// fmt.Printf("firstUser2: %+v\n", firstUser2)

	// not recommend
	// firstUser3 := &User{UserId: 10005}
	// db.Model(&User{}).First(&firstUser3)
	// fmt.Printf("firstUser3: %+v\n", firstUser3)

	// select limit 1
	firstUser4 := &User{}
	db.Where("user_id = ?", 10006).First(&firstUser4)
	fmt.Printf("firstUser4: %+v\n", firstUser4)

	firstUser5 := &User{}
	result = db.Select("user_id,user_name").Where("user_id = ? and user_name = ?", 10004, "Jinzhu").First(&firstUser5)
	fmt.Printf("firstUser5: %+v, result:%+v\n", firstUser5, result)

	// multiply select
	users := make([]User, 0)
	db.Select("user_id,user_name").Find(&users)
	fmt.Printf("users: %+v\n", users)

	users2 := make([]User, 0)
	db.Select("user_id,user_name").Where("user_id = ?", 10003).Find(&users2)
	fmt.Printf("users2: %+v\n", users2)

	type CountResult struct {
		Name  string
		Total int
	}

	// group by
	cr := &CountResult{}
	db.Model(&User{}).Select("user_name as name, count(*) as total").Where("user_name = ?", "Jinzhu").Group("user_name").Find(&cr)
	fmt.Printf("count result: %+v\n", cr)

	type CountResult2 struct {
		Total int
	}
	cr2 := &CountResult2{Total: 0}
	db.Model(&User{}).Select("count(*) as total").Where("user_name = ?", "Jinzhu").Group("user_name").Find(&cr2)
	fmt.Printf("count result: %+v\n", cr2)

	// -------------- update -------------------
	result = db.Model(&User{}).Where("user_id = ?", 10001).Update("user_name", "wolf")
	fmt.Printf("update result: %+v\n", result)

	result = db.Model(&User{}).Where("user_id = ?", 10002).Updates(map[string]interface{}{"user_name": "xiaoming", "age": 20})
	fmt.Printf("updates result: %+v\n", result)

	// -------------- delete --------------------
	result = db.Where("user_id = ?", 10006).Delete(&User{})
	fmt.Printf("delete result: %+v\n", result)

	// not recommend
	// result = db.Delete(&User{UserId: 10005})
	// fmt.Printf("delete result: %+v\n", result)

	result = db.Where("user_id = ? and user_name = ?", 10003, "aaaa").Delete(&User{})
	fmt.Printf("delete result2: %+v\n", result)

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("user_id = ? and user_name = ?", 10003, "aaaa").Delete(&User{})
	})

	fmt.Printf("delete to sql: %s\n", sql)

	// SELECT * FROM `users` where user_id = 10002 FOR UPDATE
	findUser := &User{}
	sql = db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("user_id = ?", 10000).Clauses(clause.Locking{Strength: "UPDATE"}).Find(findUser)
	})

	result = db.Where("user_id = ?", 10000).Clauses(clause.Locking{Strength: "UPDATE"}).Find(findUser)
	fmt.Printf("select for update, result2: %+v, findUser:%+v\n", result, findUser)
	fmt.Printf("select for update to sql: %s\n", sql)
}
