package main
import (
	"fmt"
	"oa/dbs"
)

func print(ch <-chan string) {
	for word := range ch {
		fmt.Printf(" %s /", word)
	}
	fmt.Println()
}

func Example() {
	fmt.Print("【全模式】：")
	print(dbs.Seg.CutAll("我来到北京清华大学"))

	fmt.Print("【精确模式】：")
	print(dbs.Seg.Cut("我来到北京清华大学", false))

	fmt.Print("【新词识别】：")
	print(dbs.Seg.Cut("他来到了网易杭研大厦", true))

	fmt.Print("【搜索引擎模式】：")
	print(dbs.Seg.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))
}
func main(){
	//Example()
	Db, _ := dbs.GetConnect()
	len,_ := dbs.GetLengthOfBlog("伟",Db)
	fmt.Println(len)
}
