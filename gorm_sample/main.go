package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
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

type FakeMsg struct {
	Id                 int64 `json:"id" gorm:"column:id"`
	FromUserid         int64 `json:"from_userid" gorm:"column:from_userid"`                   // 【首条消息发送方】，男搭讪，女推荐则是男用户；女消息，男推荐则是女用户
	ToUserid           int64 `json:"to_userid" gorm:"column:to_userid"`                       // 【首条消息接收方】，男搭讪，女推荐则是女用户；女消息，男推荐则是男用户
	MsgType            int32 `json:"msg_type" gorm:"column:msg_type"`                         // 会话创建类型
	FromReplyTime      int64 `json:"from_reply_time" gorm:"column:from_reply_time"`           // 【首条消息发送方】最近一条消息发送时间
	ToReplyTime        int64 `json:"to_reply_time" gorm:"column:to_reply_time"`               // 【首条消息接收方】最近一条消息发送时间
	FromMsgCount       int64 `json:"from_msg_count" gorm:"column:from_msg_count"`             // 【首条消息发送方】发送的消息总数，包括系统消息？ 如果是搭讪消息，但是搭讪消息有多条的时候, 这时候 >= 1 的判断就不准确了哇
	ToMsgCount         int64 `json:"to_msg_count" gorm:"column:to_msg_count"`                 // 【首条消息接收方】发送的消息总数，包括系统消息
	FromAccostDuration int64 `json:"from_accost_duration" gorm:"column:from_accost_duration"` // 【首条消息发送方】牵线回复间隔时长（秒）
	ToAccostDuration   int64 `json:"to_accost_duration" gorm:"column:to_accost_duration"`     // 【首条消息接收方】牵线回复间隔时长（秒）
	FromRead           bool  `json:"from_read" gorm:"column:from_read"`                       // 【首条消息发送方】最近一条消息是否已读
	ToRead             bool  `json:"to_read" gorm:"column:to_read"`                           // 【首条消息接收方】最近一条消息是否已读
	Ut                 int64 `json:"ut" gorm:"column:ut"`
	Ct                 int64 `json:"ct" gorm:"column:ct"`
}

func (FakeMsg) TableName() string {
	return "fake_msg"
}

type OpMonthScore struct {
	Mid           int64 `json:"mid"`
	Month         int32 `json:"month"` // 月份, 如201901， 使用整数加速索引
	Country       int32 `json:"country"`
	StartScore    int64 `json:"start_score"`    // 月初积分
	EndScore      int64 `json:"end_score"`      // 月末结余积分
	GainedScore   int64 `json:"gained_score"`   // 本月总获得积分，统一传正整数
	WithdrewScore int64 `json:"withdrew_score"` // 本月总提现积分，统一传正整数
	ClearedScore  int64 `json:"cleared_score"`  // 本月总清空积分，统一传正整数
	Ut            int64 `json:"ut"`
	Ct            int64 `json:"ct"`
}

func (p *OpMonthScore) TableName() string {
	return "op_month_score"
}
func (p OpMonthScore) SaveMonthScore(record OpMonthScore) (err error) {
	currentTs := time.Now().Unix()
	record.Ut = currentTs
	record.Ct = currentTs

	tx := gormDb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "mid"}, {Name: "month"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"gained_score":   gorm.Expr("gained_score + ?", record.GainedScore),
			"withdrew_score": gorm.Expr("withdrew_score + ?", record.WithdrewScore),
			"cleared_score":  gorm.Expr("cleared_score + ?", record.ClearedScore),
			"ut":             time.Now().Unix(),
		}),
	}).Create(&record)

	// 判断操作类型
	if tx.Error != nil {
		fmt.Println("Error occurred:", tx.Error)
	} else if tx.RowsAffected == 1 {
		fmt.Println("Insert operation")
	} else if tx.RowsAffected > 1 {
		fmt.Println("Update operation")
	}

	return
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

var gormDb *gorm.DB

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	mysqlDBs = make(map[string]*gorm.DB)

	dialector := mysql.New(mysql.Config{
		// DSN:                       "root:123456@tcp(10.10.21.65:3307)/test?charset=utf8mb4&parseTime=True&loc=Local",
		//DSN:                       "root:gH5wRmTxAmAn@tcp(101.132.227.177:12345)/testdb",
		DSN:                       "root1:Aafm6YG9Fja0@tcp(rm-wz95e87v3c5owq08u2o.mysql.rds.aliyuncs.com:3306)/me_socialize",
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Printf("MySQL连接异常:%s", err)
		return
	}

	gormDb = db

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxIdleTime(time.Minute * 20)

	mysqlDBs["test"] = db

	fmt.Println("Mysql连接成功")

	ms := OpMonthScore{
		Mid:           11111,
		Month:         202501,
		Country:       86,
		StartScore:    100,
		EndScore:      0,
		GainedScore:   8,
		WithdrewScore: 0,
		ClearedScore:  0,
		Ut:            0,
		Ct:            0,
	}

	ms.SaveMonthScore(ms)
	// mid := 117013582
	// mid := 117013747
	// st := 10000123
	// var list []FakeMsg
	// err = db.Where("ct >= ? and ( (to_userid = ? and to_accost_duration > 0 ) or (from_userid = ? and from_accost_duration > 0 ) )", st, mid, mid).
	// 	Order("ct desc").Limit(int(10)).Find(&list).Error
	// if err != nil {
	// 	fmt.Printf("find 1 err %v\n", err)
	// 	return
	// }
	// fmt.Printf("find 1 succ, result: %v\n len:%v", list, len(list))

	// var list2 []FakeMsg

	// sql1 := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Table(FakeMsg{}.TableName()).Where("ct >= ? and (to_userid = ? and to_accost_duration > 0)", st, mid).Find(&list2)
	// })

	// sql2 := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Table(FakeMsg{}.TableName()).Where("ct >= ? and (from_userid = ? and from_accost_duration > 0)", st, mid).Find(&list2)
	// })

	// unionSql := "(" + sql1 + ") union (" + sql2 + ") order by ct desc limit 10"

	// fmt.Printf("sql1:%v\n sql2:%v\n union sql:%v\n", sql1, sql2, unionSql)

	// db.Raw(`(?) union (?) order by ct desc limit ?`,
	// 	db.Table(FakeMsg{}.TableName()).Where("ct >= ? and (to_userid = ? and to_accost_duration > 0)", st, mid).Select("*"),
	// 	db.Table(FakeMsg{}.TableName()).Where("ct >= ? and (from_userid = ? and from_accost_duration > 0)", st, mid).Select("*"),
	// 	6,
	// ).Scan(&list2)
	// fmt.Printf("find 2 succ, result: %v\n len:%v\n", list2, len(list2))

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
	// firstUser := &User{}
	// result = db.First(firstUser, "user_id = ?", "10002")
	// if result.Error != nil {
	// 	fmt.Println("First failed")
	// } else {
	// 	fmt.Printf("firstUser: %+v\n", firstUser)
	// }

	// not recommend
	// firstUser2 := map[string]interface{}{}
	// db.Model(&User{}).First(&firstUser2, "user_id = ?", "10003")
	// fmt.Printf("firstUser2: %+v\n", firstUser2)

	// not recommend
	// firstUser3 := &User{UserId: 10005}
	// db.Model(&User{}).First(&firstUser3)
	// fmt.Printf("firstUser3: %+v\n", firstUser3)

	// select limit 1
	// firstUser4 := &User{}
	// db.Where("user_id = ?", 10006).First(&firstUser4)
	// fmt.Printf("firstUser4: %+v\n", firstUser4)

	// firstUser5 := &User{}
	// result = db.Select("user_id,user_name").Where("user_id = ? and user_name = ?", 10004, "Jinzhu").First(&firstUser5)
	// fmt.Printf("firstUser5: %+v, result:%+v\n", firstUser5, result)

	// // multiply select
	// users := make([]User, 0)
	// db.Select("user_id,user_name").Find(&users)
	// fmt.Printf("users: %+v\n", users)

	// users2 := make([]User, 0)
	// db.Select("user_id,user_name").Where("user_id = ?", 10003).Find(&users2)
	// fmt.Printf("users2: %+v\n", users2)

	// type CountResult struct {
	// 	Name  string
	// 	Total int
	// }

	// // group by
	// cr := &CountResult{}
	// db.Model(&User{}).Select("user_name as name, count(0) as total").Where("user_name = ?", "Jinzhu").Group("user_name").Find(&cr)
	// fmt.Printf("count result: %+v\n", cr)

	// type CountResult2 struct {
	// 	Total int
	// }
	// cr2 := &CountResult2{Total: 0}
	// db.Model(&User{}).Select("count(0) as total").Where("user_name = ?", "Jinzhu").Group("user_name").Find(&cr2)
	// fmt.Printf("count result: %+v\n", cr2)

	// // -------------- update -------------------
	// result = db.Model(&User{}).Where("user_id = ?", 10001).Update("user_name", "wolf")
	// fmt.Printf("update result: %+v\n", result)

	// result = db.Model(&User{}).Where("user_id = ?", 10002).Updates(map[string]interface{}{"user_name": "xiaoming", "age": 20})
	// fmt.Printf("updates result: %+v\n", result)

	// // -------------- delete --------------------
	// result = db.Where("user_id = ?", 10006).Delete(&User{})
	// fmt.Printf("delete result: %+v\n", result)

	// // not recommend
	// // result = db.Delete(&User{UserId: 10005})
	// // fmt.Printf("delete result: %+v\n", result)

	// result = db.Where("user_id = ? and user_name = ?", 10003, "aaaa").Delete(&User{})
	// fmt.Printf("delete result2: %+v\n", result)

	// sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Where("user_id = ? and user_name = ?", 10003, "aaaa").Delete(&User{})
	// })

	// fmt.Printf("delete to sql: %s\n", sql)

	// // SELECT * FROM `users` where user_id = 10002 FOR UPDATE
	// findUser := &User{}
	// sql = db.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Where("user_id = ?", 10000).Clauses(clause.Locking{Strength: "UPDATE"}).Find(findUser)
	// })

	// result = db.Where("user_id = ?", 10000).Clauses(clause.Locking{Strength: "UPDATE"}).Find(findUser)
	// fmt.Printf("select for update, result2: %+v, findUser:%+v\n", result, findUser)
	// fmt.Printf("select for update to sql: %s\n", sql)
}
