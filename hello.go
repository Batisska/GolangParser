package main

import (
	"mylib/helpers"
	"fmt"
	"log"
	"github.com/PuerkitoBio/goquery"
	"database/sql"
	"strconv"
)

const (
	suplier = ".Product__info > .Product__code"
)

type scrape struct {
	doc *goquery.Document
	rowsId int
}
//send <- chan Results, results chan <- Results
func scrapeAttrs(doc_chan <-chan scrape, results chan<- map[string][]map[string]string) {
	for docs := range doc_chan {
		doc := docs.doc
		thisMap := make(map[string][]map[string]string)

		id_supliers := doc.Find(suplier)

		prod_id := id_supliers.Text()

		// Find the review items
		res := doc.Find(".ProductDetails").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			s.Find("li").Each(func(i int, g *goquery.Selection) {
				/*attr_name["attr"] = append(attr_name["attr"], g.Text())
				text := Attr{Name:}
				likes[prod_id] = append(likes[prod_id], text)*/
				aMap := map[string]string{
					"attr": g.Text()}
				thisMap[prod_id] = append(thisMap[prod_id], aMap)
			})
		})

		if len(res.Text()) == 0 {
			res = doc.Find(".Product__info > .Product__features > li").Each(func(i int, s *goquery.Selection) {
				aMap := map[string]string{
					"attr": s.Text(),
				}
				thisMap[prod_id] = append(thisMap[prod_id], aMap)
			})

		}
		s := map[string]string{
			"rows_id": strconv.Itoa(docs.rowsId),
		}
		thisMap[prod_id] = append(thisMap[prod_id], s)
		results <- thisMap
	}
}

func sqlattrs(c <-chan map[string][]map[string]string, stmt *sql.Stmt) {

	var attr []string
	var id string

	for s := range c {
		for t, c := range s {
			attr = append(attr, t)

			for _, b := range c {
				attr = append(attr, b["attr"])
				res, err := stmt.Query(t, b["attr"])
				if err != nil {
					log.Println("sqlattr -> ", err)
				}
				res.Close()
				id =  b["rows_id"]
			}
		}
		res := helpers.Query("UPDATE `parser`.`prod_links` SET `active`='1' WHERE `id`='"+id+"'") //where `active` = '1'
		res.Close()
	}


}

func main() {

	db := helpers.Db()
	defer db.Close()
	res, err := db.Query("SELECT * FROM `parser`.`prod_links` where `active` = '0'")
	if err != nil {
		log.Fatal("Query -> ", err)
	}
	defer res.Close()

	query := "INSERT INTO prod_attr(id_supliers,attr) VALUES (?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("Prepare -> ", err)
	}

	defer stmt.Close()

	var (
		id     int
		link   string
		active int
	)

	// получить страничку goquery.Document
	send := make(chan helpers.Results, 100)
	result := make(chan helpers.Results, 100)

	// Граббер
	result_scrape := make(chan map[string][]map[string]string, 100)
	send_scrape := make(chan scrape, 100)

	//Вставить в базу
	send_sqlattrs := make(chan map[string][]map[string]string, 100)

	for i := 0; i < 15; i++ {
		go helpers.DocGoquery(send, result)
	}

	for i := 0; i < 15; i++ {
		go scrapeAttrs(send_scrape, result_scrape)
	}
	defer close(send_scrape)


	for i := 0; i < 15; i++ {
		go sqlattrs(send_sqlattrs, stmt)
	}
	defer close(send_sqlattrs)


	go func () {
		i := 0
		for res.Next() {
			if err := res.Scan(&id, &link, &active); err != nil {
				log.Fatal(err)
			}
			x := helpers.Results{I: i, Url: link,Id:id}
			send <- x
			i++
		}
		close(send)
	}()

	a := 0
	for {
		select {
		case x := <-result:
			// id_code не скрабить если это не страница товара.
			id_code := x.Doc.Selection.Find(suplier)
			if x.Statuscode == 200 && id_code != nil {
				fmt.Println("for send ", x.I, "for result", a, x.Statuscode)
				f := scrape{x.Doc, x.Id}
				send_scrape <- f
				a++
			}
		case y := <-result_scrape:
			send_sqlattrs <- y


		}
	}
}
