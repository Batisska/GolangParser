package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"github.com/djimenez/iconv-go"
	_"github.com/go-sql-driver/mysql"
	"regexp"
	"database/sql"
	"log"
	"strings"
	"github.com/go-sql-driver/mysql"
	"os"
	"time"
)

func Scrape(url string) {

	os.Setenv("HTTP_PROXY", "89.36.215.14:1189")

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	res, err := netClient.Get(url)
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

	var appends []string

	// Find the review items
	doc.Find(".Content a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		// (/catalog\/)
		re := regexp.MustCompile(`(/catalog/)`)
		//var digitsRegexp = regexp.MustCompile(`(/catalog/)`)
		//fmt.Println(digitsRegexp.FindStringSubmatch(s.Attr("href")))

		band, ok := s.Attr("href")
		if ok {
			if re.FindString(band) == "/catalog/" {
				//appends = append(appends, band)
				band := strings.TrimSpace(band)
				appends = append(appends, "http://www.office-planet.ru"+band)
			}

		}

	})
	insertCatlinks(appends)
}

func main() {
	Scrape("http://www.office-planet.ru/catalog/abc/")
}

func insertCatlinks(vals []string) {
	db, err := sql.Open("mysql", "root:root@/parser")
	if err != nil {
		log.Fatal("Open -> ", err)
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO catalog_links(link) VALUES (?)")
	if err != nil {
		log.Fatal("Prepare -> ", err)
	}

	for _, ok := range vals {
		res, err := stmt.Exec(ok)

		if err != nil {
			me, ok := err.(*mysql.MySQLError)
			if !ok {
				log.Fatal("Exec -> ", err)
			}
			if me.Number == 1062 {
				continue
			}
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal("LastInsertId -> ", err)
		}
		rowCnt, err := res.RowsAffected()
		if err != nil {
			log.Fatal("RowsAffected -> ", err)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	}
}