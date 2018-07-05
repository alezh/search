package data

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/huichen/sego"
	"github.com/alezh/search/searchLogic"
)

var DataSourceName *xorm.Engine

func InsertData(Segmenter sego.Segmenter)  {

	goods := make([]*Goods,0)

	DataSourceName = NewCgsOsSource()
	DataSourceName.Find(&goods)

	f := func(goods *Goods,Segmenter sego.Segmenter) {
		text := []byte(goods.Title)
		segments := Segmenter.Segment(text)
		seg := sego.SegmentsToSlice(segments, true)
		for _,v := range seg {
			if v != ""{
				searchLogic.InsertString(v,goods.Id)
			}
		}
	}

	for _,v := range goods{
		f(v,Segmenter)
	}
}



type Goods struct {
	Id          int
	Title       string
	CreatedAt   string  `xorm:"created TIMESTAMP"`
	UpdatedAt   string  `xorm:"updated TIMESTAMP"`
	DeletedAt   string  `xorm:"deleted TIMESTAMP"`
}



func NewCgsOsSource() *xorm.Engine {

	//dataSourceName := "root:123456@tcp(localhost:3306)/os?charset=utf8mb4"
	dataSourceName := "chenchunlai:qazwsx123@tcp(172.16.10.230:3306)/cgsos?charset=utf8mb4"

	engine, err := xorm.NewEngine("mysql", dataSourceName)

	if err!=nil{
		fmt.Println("orm failed to initialized")
		return nil
		//panic("orm failed to initialized")
		//return nil,errors.New("orm failed to initialized")
	}
	if errs := engine.Ping(); errs != nil{
		//fmt.Println("orm failed to initialized")
		return nil
	}
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "sc_")

	engine.SetTableMapper(tbMapper)
	//日志打印SQL
	engine.ShowSQL(false)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(1024)
	//设置最大打开连接数
	engine.SetMaxOpenConns(2048)

	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	return engine
}