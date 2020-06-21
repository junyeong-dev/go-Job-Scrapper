package main

import (
	"log"
	"net/http"
	"std/fmt"
	"strconv"
	"strings"

	// goquery - jquery처럼 css selector를 통해 원하는 요소를 쉽게 찾을 수 있게 도와줌
	"github.com/PuerkitoBio/goquery"
)

// 취업 정보 struct
type extractedJod struct {
	id       string
	title    string
	location string
	salay    string
	summary  string
}

// 취업 정보를 가져올 base URL
var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	totalPages := getPages()
	fmt.Println(totalPages)
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}

}

// 페이지에 해당하는 URL을 return 하는 함수
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// 취업 정보를 가져옴; 취업 정보 카드
	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		extracteJod(card)
	})
}

// 취업 정보 카드에서 해당하는 데이터 추출
func extracteJod(card *goquery.Selection) {
	// Attr - 데이터와 존재여부를 return
	id, _ := card.Attr("data-jk")
	// Find - 원하는 속성을 가져옴
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	fmt.Println(id, title, location, salary, summary)
}

// 공백을 제거하고, 문자열의 배열로 만들어 준 후 다시 공백을 넣은 하나의 문자열로 만들어 return
// ex) "hello      golang      unbelievable" -> "hello", "golang", "unbelievable" -> "hello golang unbelievable"
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 전체 페이지를 return 하는 함수
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// 페이지를 가져오는데 Error가 발생하는지 체크
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// 페이지 코드를 체크
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
