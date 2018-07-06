package main

import (
	//"github.com/alezh/search/searchLogic"
	//"github.com/alezh/search/data"
	"github.com/huichen/sego"
	"github.com/alezh/search/searchLogic"
	"github.com/alezh/search/data"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"strconv"
)
var Segmenter sego.Segmenter

func init()  {
	Segmenter.LoadDictionary("C:/Users/user/go/src/github.com/huichen/sego/data/dictionary.txt")
}

func main()  {
	searchLogic.MPQHashTableInit()

	data.InsertData(Segmenter)

	router := httprouter.New()
	router.GET("/get/:name", Index)
	router.GET("/product/:key", Gets)

	log.Fatal(http.ListenAndServe("172.16.13.42:8080", router))


}


func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	text := []byte(ps.ByName("name"))
	segments := Segmenter.Segment(text)
	seg := sego.SegmentsToSlice(segments, true)
	for _,v :=range seg {
		if idss,ok := searchLogic.GetHashTableIsExist(v);ok{
			for _,i:=range idss {
				w.Write([]byte(strconv.Itoa(i)+"\n"))
			}
		}
	}
	w.Write([]byte("完成\n"))
}

func Gets(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	text := []byte(ps.ByName("key"))
	segments := Segmenter.Segment(text)
	seg := sego.SegmentsToSlice(segments, true)
	for _,v :=range seg {
		if idss,ok := searchLogic.GetHashTableIsExist(v);ok{
			for _,i:=range idss {
				goods := &data.Goods{}
				data.DataSourceName.Where("id=?",i).Get(goods)
				w.Write([]byte(strconv.Itoa(i)+"|"+goods.Title+"\n"))
			}
		}
	}
	w.Write([]byte("完成\n"))
}