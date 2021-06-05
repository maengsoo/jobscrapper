package scrapper

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
	location string
	salary   string
	summary  string
}

func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJob        // 빈 배열 선언
	c := make(chan []extractedJob) // []extractedJob의 정보를 받는 채널 생성
	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c) // 일자리 정보(카드묶음)를 채널에 요청함
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c                  // totalPages의 수만큼 채널에서 정보를 받음
		jobs = append(jobs, extractedJobs...) // 배열에 저장
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

// 각 페이지의 있는 일자리를 반환
func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(page*50) // 필요한 주소
	fmt.Println("Requesting:", pageURL)
	res, err := http.Get(pageURL) // 정보를 요청
	checkErr(err)                 // 에러 체크
	checkStatusCode(res)          // status code 체크

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	// html class, id 에서 정보가져오기
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	// searchCards의 길이 만큼 goroutines 실행
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

// id, title, name, location 정보 가져오기
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	const classJobID = ".data-jk"
	const classTitle = ".title>a"
	const classLocation = ".sjcl"
	const classSalary = ".salaryText"
	const classSummary = ".summary"

	id, _ := card.Attr(classJobID) // 직업카드 id
	title := CleanString(card.Find(classTitle).Text())
	location := CleanString(card.Find(classLocation).Text())
	salary := CleanString(card.Find(classSalary).Text())
	summary := CleanString(card.Find(classSummary).Text())

	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary}
}

// 가져온 id, title, name, location 의 공백을 없애고(TrimSpace), 배열을 만들고(Fields), 배열을 스페이스를 기준으로 한줄에 string(Join)으로 만듬
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 웹페이지의 총 페이지수 가져오기
func getPages(url string) int {
	pages := 0

	res, err := http.Get(url)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// 가져온 직업카드를 csv 파일로 저장
func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	c := make(chan []string)

	// jobs.csv에 형식에 맞게 job 정보 요청
	for _, job := range jobs {
		go writeInCsv(job, c)
	}

	// jobs 길이 만큼 채널에서 정보를 받아옴
	for i := 0; i < len(jobs); i++ {
		jobSlice := <-c
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

// 형식에 맞게 정보를 불러옴
func writeInCsv(job extractedJob, c chan<- []string) {
	c <- []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
}

// Check url error
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Check url statuscode
func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 { // 200이면 정상
		log.Fatalln("Request failed with status:", res.StatusCode)
		// - Status code 자릿수 -
		// 100 : 조건부 응답
		// 200 : 성공
		// 300 : 리다이렉션 완료(추가 동작 필요)
		// 400 : 상태 오류
		// 500 : 서버 오류
	}
}
