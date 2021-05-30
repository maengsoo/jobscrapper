package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	name     string
	location string
	salary   string
	summary  string
}

var baseURL string = "https://kr.indeed.com/jobs?q=java&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
	fmt.Println(jobs)
}

// 한페이지 배열을 나열하는 메소드
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting:", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tapItem")

	// html class, id 에서 정보가져오기
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

// id, title, name, location 정보 가져오기
func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".jobTitle>span").Text())
	name := cleanString(card.Find(".companyName").Text())
	location := card.Find(".company_location .companyLocation").Text()
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".job-snippet").Text())

	return extractedJob{
		id:       id,
		title:    title,
		name:     name,
		location: location,
		salary:   salary,
		summary:  summary}
}

// 가져온 id, title, name, location 의 공백을 없애고(TrimSpace), 배열을 만들고(Fields), 배열을 스페이스를 기준으로 한줄에 string(Join)으로 만듬
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 페이징 수 가져오기
func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// Check url error
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Check url statuscode
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status:", res.StatusCode)
	}
}
