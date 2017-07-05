package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"github.com/djimenez/iconv-go"
	//"regexp"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
)

func ExampleScrape(url string) {

	res, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer res.Body.Close()

	// Convert the designated charset HTML to utf-8 encoded HTML.
	// `charset` being one of the charsets known by the iconv package.
	utfBody, err := iconv.NewReader(res.Body, "windows-1251", "utf-8")
	if err != nil {
		// handler error
	}
	// use utfBody using goquery
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		// handler error
	}
	//var digitsRegexp = regexp.MustCompile(`[^-]+`)

		id_supliers := doc.Find(".Product__info > .Product__code")
		// Find the review items
		doc.Find(".ProductDetails").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		s.Find("li").Each(func(i int, g *goquery.Selection) {
			fmt.Println(g.Text(),id_supliers.Text())
			insert(g.Text(), id_supliers.Text())
			//fmt.Println(digitsRegexp.FindStringSubmatch(g.Text()))
		})
	})
}

func main() {
	ExampleScrape("http://www.office-planet.ru/catalog/goods/mfu-lazernyje-monohromnyje2/352962/#js-tabDetailInfo")
}

func insert(rows string, id_supplier string)  {

	db, err := sql.Open("mysql", "root:root@/parser")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO office(string, id_supplier) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(rows,id_supplier)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

}
