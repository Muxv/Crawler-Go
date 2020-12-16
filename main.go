package main

import (
	m "Crawler-go/models"
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	url       = "https://book.douban.com/top250"
	BookLists = make([]m.TopBook, 250, 250)
	numbers   = regexp.MustCompile("[0-9]+")
)

const pageCount = 10
const PageSize = 25
const agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/87.0.4280.88 Safari/537.36"

type Empty struct{}

func main() {
	finish := make(chan Empty)
	for i := 0; i < pageCount; i++ {
		go func(page int, ch chan Empty) {
			resp := CrawOnePage(page * PageSize)
			defer resp.Body.Close()
			ParseBookInfo(resp.Body, page)
			finish <- Empty{}
		}(i, finish)
	}

	for i := 0; i < pageCount; i++ {
		<-finish
	}

	log.Println("Now start to save books in mysql")
	for i := 0; i < PageSize*pageCount; i++ {
		BookLists[i].Insert()
	}
	log.Println("Insert Over")
}

func CrawOnePage(pageNum int) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Fail generate request. ", err)
	}
	req.Header.Add("User-Agent", agent)

	q := req.URL.Query()
	q.Add("start", strconv.Itoa(pageNum))

	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Fail to receive response. ", err)
	}
	return resp
}

func ParseBookInfo(r io.Reader, page int) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal("Fail to generate doc. ", err)
	}

	doc.Find(".article .indent table").Each(func(i int, selection *goquery.Selection) {
		cName := selection.Find("td .pl2 a").Text()
		cName = strings.Join(strings.Fields(cName), "")
		eName := selection.Find("td .pl2>span").Text()
		bInfo := selection.Find("td p.pl").Text()
		rStr := selection.Find("td span.rating_nums").Text()
		ra, _ := strconv.ParseFloat(rStr, 64)
		rnStr := selection.Find("td div.star.clearfix .pl").Text()
		rStr = numbers.FindString(rnStr)
		rn, _ := strconv.Atoi(rStr)
		com := selection.Find("td p.quote span").Text()

		book := m.TopBook{
			Topk:      page*PageSize + i,
			ChName:    cName,
			EnName:    eName,
			BasicInfo: bInfo,
			Rank:      ra,
			RankNum:   rn,
			Comment:   com,
		}
		BookLists[book.Topk] = book
	})
}

func saveRead(r io.Reader, filename string) {
	if filename == "" {
		filename = "./store.txt"
	}
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Fail to open %s, %s", filename, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.ReadFrom(r)
	if err != nil {
		log.Fatal("Fail to read from reader", err)
	}
}
