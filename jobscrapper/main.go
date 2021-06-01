package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
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
	var jobs []extractedJob // 빈 배열 선언
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)           // 일자리 정보를 가져와
		jobs = append(jobs, extractedJobs...) // 배열에 저장
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "name", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.name, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}

}

// 각 페이지의 있는 일자리를 반환
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // 필요한 주소
	fmt.Println("Requesting:", pageURL)
	res, err := http.Get(pageURL) // 정보를 요청
	checkErr(err)                 // 에러 체크
	checkCode(res)                // status code 체크

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

	id, _ := card.Attr("data-jk") // 직업카드 id
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

// 웹페이지의 총 페이지수 가져오기
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
	if res.StatusCode != 200 { // 200이면 정상
		log.Fatalln("Request failed with status:", res.StatusCode)
		// - Status code 자릿수 -
		// 100 : 조건부 응답
		// 200 : 성공
		// 300 : 리다이렉션 완료(추가 동작 취해야함)
		// 400 : 상태 오류
		// 500 : 서버 오류
	}
}
