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

// myuserテーブルと同じ構造にしておく。GoでDBとデータのやり取りをするときはこのStructを使用する。
// User structと中身が同じだが、今回のAPIの仕様上たまたまそうなっただけなのでDB用に別で定義している。
type Record struct {
	ID       int
	Bookname string
	URL      string
	Comment  string
	Time     string
}

func main() {
	// まずはデータベースに接続する。(パスワードとか違うと思うので随時修正して)
	dsn := "host=localhost user=postgres password=Hach8686 dbname=test port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// エラーでたらプロセス終了
		log.Fatalf("Some error occured. Err: %s", err)
	}

	/**
	 * ここでAPIサーバーの設定をする。
	 * rはrouterの略で何のAPIを用意するかを定義する。
	 * 例えば、r.GET("/postpage",　というのはGETメソッドでhttp://localhost:8080/postpageが叩かれたら処理をスタートするAPIを提供しているということ。
	 * 今回用意してるのは/postpage　GET、/showpage　GET、/user　POSTの三つのAPI(これを通常エンドポイントと呼びます)
	 */
	r := gin.Default()

	// ginに対して、使うHTMLのテンプレートがどこに置いてあるかを知らせる。
	r.LoadHTMLGlob("temp/*")

	// 用意していないエンドポイント以外を叩かれたら内部で/showpage　GETを叩いてデフォルトページを表示する様にする。
	r.NoRoute(func(c *gin.Context) {
		location := url.URL{Path: "/showpage"}
		c.Redirect(http.StatusFound, location.RequestURI())
	})

	// POST用のページ（post.html）を返す様にしている。
	// c.HTMLというのはこのAPIのレスポンスとしてHTMLファイルを返すよ、という意味
	r.GET("/postpage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post.html", gin.H{})
	})

	// 結果を表示するページを返す様にしている。
	r.GET("/showpage", func(c *gin.Context) {
		// このページにはDBに保存されているデータを全件表示する様にするため、HTMLファイルを返す前にDBからデータを全部とってきている。
		// Recordは上で定義したStruct。データが複数あるので、このStructのスライス型をrecordsに宣言する。
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
		// Userは上で定義したStruct。HTMLから送られてくるリクエストデータを割り付ける用のuserを宣言。
		var book Bookmark
		if err := c.ShouldBind(&book); err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument"})
			return
		}
		var record Record
		// 以下の様にしてInsert文を書いて、リクエストデータをDBに書くこむ。.Scan(&record)はDBに書き込む際に必要らしい。
		// recordはmyuserテーブルと構造を同じにしている。(Gormのお作法)
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

	// サーバーを立ち上げた瞬間は一旦ここまで実行されてListening状態となる。
	// r.POST( や　r.GET(　等の関数はAPIが呼ばれる度に実行される。
	r.Run()
}
