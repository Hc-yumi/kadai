package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// HTMLからリクエスト来た時にGo内でそのデータが受け取れるようにこのStructを用意する。
type Bookmark struct {
	Name    string `form:"bookName"`
	URL     string `form:"bookUrl"`
	Comment string `form:"bookcomment"`
}

// booklistテーブルと同じ構造。
type Record struct {
	ID       int
	Bookname string
	URL      string
	Comment  string
	Time     string
}

func main() {
	// まずはデータベースに接続する。(パスワードは各々異なる)
	dsn := "host=localhost user=postgres password=Hach8686 dbname=test port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// エラーでたらプロセス終了
		log.Fatalf("Some error occured. Err: %s", err)
	}

	/*
	 * APIサーバーの設定をする。
	 * rはrouterの略で何のAPIを用意するかを定義する。
	 * postpage　GET、/showpage　GET、/user　POST
	 */
	
	r := gin.Default()

	// ginに対して、使うHTMLのテンプレートがどこに置いてあるかを知らせる。
	r.LoadHTMLGlob("temp/*")

	// 用意していないエンドポイント以外を叩かれたら内部で/showpage　GETを叩いてデフォルトページを表示する様にする。
	r.NoRoute(func(c *gin.Context) {
		location := url.URL{Path: "/showpage"}
		c.Redirect(http.StatusFound, location.RequestURI())
	})

	// POST用のページ（post.html）を返す。
	// c.HTMLというのはこのAPIのレスポンスとしてHTMLファイルを返すよ、という意味
	r.GET("/postpage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post.html", gin.H{})
	})

	// 結果を表示するページを返す。
	r.GET("/showpage", func(c *gin.Context) {
		var records []Record
		// &recordsをDBに渡して、取得したデータを割り付ける。
		dbc := conn.Raw("SELECT id, bookname,url,comment,to_char(time,'YYYY-MM-DD HH24:MI:SS') AS time FROM booklist").Scan(&records)
		if dbc.Error != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// レスポンスとして、show.htmlを返すが、一緒にrecordsも返している。これにより、HTML内でデータをマッピング表示することができる。
		c.HTML(http.StatusOK, "show.html", gin.H{
			"Books": records,
		})
	})

	// データを登録するAPI。POST用のページ（post.html）の内部で送信ボタンを押すと呼ばれるAPI。
	r.POST("/book", func(c *gin.Context) {
		
		var book Bookmark
		if err := c.ShouldBind(&book); err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument"})
			return
		}
		var record Record
		// 以下の様にしてInsert文を書いて、リクエストデータをDBに書くこむ。.Scan(&record)はDBに書き込む際に必要らしい。
		// recordはbooklistテーブルと構造を同じにしている。(Gormのお作法)
		dbc := conn.Raw(
			"insert into booklist(bookname, url, comment) values(?, ?, ?)",
			book.Name, book.URL, book.Comment).Scan(&record)
		if dbc.Error != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// DBへの保存が成功したら結果を表示するページに戻るために/showpageのAPIを内部で読んでそちらでページの表示を行う。
		location := url.URL{Path: "/showpage"}
		c.Redirect(http.StatusFound, location.RequestURI())
	})


	// データの削除
	r.DELETE("/book/:id", func(c *gin.Context) {
		id:=c.Param("id")
		fmt.Println("id is ", id)
		var records []Record
		dbc := conn.Raw("DELETE FROM booklist where id=?",id).Scan(&records)
		
		if dbc.Error != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		location := url.URL{Path: "/showpage"}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())
	})




	// showpageで書籍名をinputしてボタン押したら→入力した書籍名とおなじ列をdeleteする
	r.DELETE("/book/:bookname", func(c *gin.Context) {
		bookname:=c.Param("bookname")
		fmt.Println("bookname is ", bookname)
		var records []Record
		dbc := conn.Raw("DELETE FROM booklist where bookname=?",bookname).Scan(&records)
		
		if dbc.Error != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		location := url.URL{Path: "/showpage"}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())

		// showpageで書籍名をinputしてボタン押したら→入力した書籍名とおなじ列をdeleteする

	})



	// サーバーを立ち上げた瞬間は一旦ここまで実行されてListening状態となる。
	// r.POST( や　r.GET(　等の関数はAPIが呼ばれる度に実行される。
		r.Run()

}
