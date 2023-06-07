package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Info struct {
	CreateTime time.Time `form:"create_time" binding:"required,timing" time_format:"2006-01-02"`
	UpdateTime time.Time `form:"update_time" binding:"required,timing" time_format:"2006-01-02"`
}

// 自定义验证规则断言
func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

type User struct {
	Username string `validate:"min=6,max=10"`
	Age      uint8  `validate:"gte=1,lte=10"`
	Sex      string `validate:"oneof=female male"`
}

func main() {
	validate := validator.New()
	user1 := User{Username: "asong", Age: 11, Sex: "null"}
	err := validate.Struct(user1)
	if err != nil {
		fmt.Println(err)
	}

	user2 := User{Username: "asong111", Age: 8, Sex: "male"}
	err = validate.Struct(user2)
	if err != nil {
		fmt.Println(err)
	}

	route := gin.Default()
	// 注册验证
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("timing", timing)
		if err != nil {
			fmt.Println("success")
		}
	}

	route.GET("/time", getTime)
	route.Run(":8080")
}

func getTime(c *gin.Context) {
	var b Info
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "time are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
