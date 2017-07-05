package main

import (
	"github.com/PuerkitoBio/goquery"
	_"github.com/go-sql-driver/mysql"
	"log"
	"mylib/helpers"
)

func main() {

	db := helpers.Db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM pag_links WHERE active = 0")
	helpers.CheckErr(err)

	for rows.Next() {
		var id int
		var link string
		var active int
		err = rows.Scan(&id, &link, &active)
		helpers.CheckErr(err)
		log.Println(
			"ID => ", id,
			"catalog_links => ", link,
		)
		err := ScrapeProdLinks(link, id)

		if err != nil {
			continue
		}
		helpers.UpdateID(id, "UPDATE `parser`.`pag_links` SET `active`='1' WHERE `id`=(?)")
	}
}

func ScrapeProdLinks(url string, id int) (err error) {

	var appends []string
	doc, statuscode, err := helpers.Get_urls_rus(url)

	if err != nil || statuscode != 200 {
		return err
	}

	//helpers.CheckStatus(statuscode,url)

	// Find the review items
	doc.Find(".Product__link").Each(func(i int, s *goquery.Selection) {

		// For each item found, get the band and title
		band, ok := s.Attr("href")
		if ok {
			appends = helpers.Unique(append(appends, "http://www.office-planet.ru"+band))
		}
	})
	helpers.Insertlinks(appends, "INSERT INTO prod_links(link) VALUES (?)")

	return err
}
