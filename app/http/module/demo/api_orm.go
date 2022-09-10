package demo

import (
	"database/sql"
	"time"

	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
	"github.com/26huitailang/yogo/framework/provider/orm"
)

// DemoOrm Orm的路由方法
func (api *DemoApi) DemoOrm(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)

	// init orm.DB
	gormService := c.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := gormService.GetDB(orm.WithConfigPath("database.default"))
	if err != nil {
		logger.Error(c, err.Error(), nil)
		c.AbortWithError(500, err)
		return
	}
	db.WithContext(c)

	err = db.AutoMigrate(&User{})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "migrate ok", nil)

	// insert one row
	email := "foo@bar.com"
	name := "foo"
	age := uint8(25)
	birthday := time.Date(2022, 9, 4, 11, 56, 0, 0, time.Local)
	user := &User{
		Name:         name,
		Email:        &email,
		Age:          age,
		Birthday:     &birthday,
		MemberNumber: sql.NullString{},
		ActivatedAt:  sql.NullTime{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = db.Create(user).Error
	logger.Info(c, "insert user", map[string]interface{}{
		"id":  user.ID,
		"err": err,
	})
	// 更新一条数据
	user.Name = "bar"
	err = db.Save(user).Error
	logger.Info(c, "update user", map[string]interface{}{"err": err, "id": user.ID})
	// 查询一条数据
	queryUser := &User{ID: user.ID}
	err = db.First(queryUser).Error
	logger.Info(c, "query user", map[string]interface{}{"err": err, "name": queryUser.Name})
	// 删除一条数据
	err = db.Delete(queryUser).Error
	logger.Info(c, "delete user", map[string]interface{}{"err": err, "id": user.ID})
	c.JSON(200, "ok")
}
