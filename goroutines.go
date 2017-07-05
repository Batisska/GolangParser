package main

import (
	//"fmt"
	"mylib/helpers"
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
	"database/sql"
)
func scrapeAttr(c chan map[string][]map[string]string, doc *goquery.Document) <- chan map[string][]map[string]string {

	thisMap := make(map[string][]map[string]string)

	id_supliers := doc.Find(".Product__info > .Product__code")

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
	c <- thisMap
	return c
}


func sqlattr(c chan map[string][]map[string]string, stmt *sql.Stmt)  {


	var attr []string

	for s := range c {
		for t,c := range s{
			attr = append(attr,t)
			for _, b := range c {
				attr = append(attr,b["attr"])

				res, err := stmt.Query(t,b["attr"])
				if err != nil {
					log.Println("sqlattr -> ", err)
				}
				 res.Close()
				log.Println("add database")
			}
		}
	}
}

func main() {

	db := helpers.Db()
	defer db.Close()

	query := "INSERT INTO prod_attr(id_supliers,attr) VALUES (?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("Prepare -> ", err)
	}
	defer stmt.Close()


	res := helpers.Query("SELECT * FROM `parser`.`prod_links`") //where `active` = '1'
	defer res.Close()
	var (
		id     int
		link   string
		active int
	)

	result := make(chan helpers.Results, 15)
	c := make(chan map[string][]map[string]string, 10)
	b := make(chan map[string][]map[string]string, 2)
	go sqlattr(b, stmt)


	for res.Next() {
		if err := res.Scan(&id, &link, &active); err != nil {
			log.Fatal(err)
		}
		go helpers.DocGoquery(result)
	}
	for {
		select {
		case result := <-result:
			fmt.Println(result.Statuscode)
			go scrapeAttr(c, result.Doc)
		case ress := <- c:
			 b <- ress

		}
	}

}
