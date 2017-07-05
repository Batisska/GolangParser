package main

import (
	"mylib/helpers"
	"github.com/PuerkitoBio/goquery"
)

func main(){

	db := helpers.Db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM catalog_links")
	helpers.CheckErr(err)

	for rows.Next() {
		helpers.CheckIp()
		var id int
		var link string
		err = rows.Scan(&id, &link)
		helpers.CheckErr(err)

		doc, statuscode := helpers.Get_urls_rus(link)

		helpers.CheckStatus(statuscode, link)

		var appends []string

		// Find the review items
		doc.Find(".Paginator a").Each(func(i int, s *goquery.Selection) {

			// For each item found, get the band and title
			band, ok := s.Attr("href")
			if ok {
				helpers.Write("paginatorlinks.log", band)
				appends = helpers.Unique(append(appends, "http://www.office-planet.ru"+band))
			}
		})
		appends = append(appends, link)
		helpers.Insertlinks(appends, "INSERT INTO pag_links(link) VALUES (?)")
	}
}