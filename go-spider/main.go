package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

type articleUrls struct {
	s    []string
	i    int
	lock sync.Mutex
}

func main() {
	urls := getArticleUrls()
	contents := getAllArticleContents(urls)
	writeTxt("./无聊小说.txt", contents)
}

func writeTxt(path string, contents []string) {
	fw, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fw.Close()
	for i, v := range contents {
		fmt.Printf("第%d章", i+1)
		_, err = fw.WriteString(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getAllArticleContents(urls []string) []string {
	articlesTotal := len(urls)
	groupSize := 50
	groupCount := articlesTotal / groupSize
	if articlesTotal%groupSize > 0 {
		groupCount++
	}
	contents := make([]string, articlesTotal)
	for i := 0; i < groupCount; i++ {
		l := groupSize
		if i == (groupCount - 1) {
			l = articlesTotal - i*groupSize
		}
		wait := &sync.WaitGroup{}
		wait.Add(l)
		for j := 0; j < l; j++ {
			index := i*groupSize + j
			url := urls[index]
			go getArticleContent(url, &contents[index], wait)
		}
		wait.Wait()
	}

	return contents
}

func getArticleContent(url string, contents *string, wait *sync.WaitGroup) {
	// 处理第一页
	url = "http://www.soruncg.com" + url
	doc := getDocument(url)
	title := doc.Find("h1.title").Text()
	//indexReg := regexp.MustCompile(`\d+`)
	//index, _ := strconv.Atoi(indexReg.FindString(title))
	content, _ := doc.Find("#content").Html()
	divReg, _ := regexp.Compile(`<div.+</div>`)
	content = strings.TrimSpace(divReg.ReplaceAllString(content, ""))
	brReg, _ := regexp.Compile(`(<br/>　*)+`)
	content = strings.TrimSpace(brReg.ReplaceAllString(content, "\n"))
	// 处理后续页
	urlReg, _ := regexp.Compile(`/(\d+)\.html`)
	nextPageUrl := urlReg.ReplaceAllString(url, `/${1}_2.html`)
	doc = getDocument(nextPageUrl)
	content1, _ := doc.Find("#content").Html()
	content1 = strings.TrimSpace(divReg.ReplaceAllString(content1, ""))
	content1 = strings.TrimSpace(brReg.ReplaceAllString(content1, "\n"))
	*contents = title + "\n" + content + "\n" + content1 + "\n"
	fmt.Println(title)
	wait.Done()
}

func getArticleUrls() []string {
	//latest := getIndexDoc(1).Find(".top p").Last().Children().Text()
	//r := regexp.MustCompile(`\d+`)
	//articlesTotal, _ := strconv.Atoi(r.FindString(latest))
	articlesTotal := 1745
	urls := &articleUrls{
		s: make([]string, articlesTotal),
	}
	wait := &sync.WaitGroup{}
	pageTotal := 88
	wait.Add(pageTotal)
	for page := 1; page <= pageTotal; page++ {
		//go getArticleUrl(page, urls, wait)
		getArticleUrl(page, urls, wait)
	}
	wait.Wait()
	return urls.s
}

func getArticleUrl(page int, urls *articleUrls, wait *sync.WaitGroup) {
	getIndexDoc(page).Find(".section-list").Last().Find("a").Each(func(i int, a *goquery.Selection) {
		urls.lock.Lock()
		fmt.Println(urls.i)
		urls.s[urls.i], _ = a.Attr("href")
		urls.i++
		urls.lock.Unlock()
	})
	wait.Done()
}

func getIndexDoc(page int) *goquery.Document {
	url := fmt.Sprintf("http://www.soruncg.com/ldks/63627/index_%d.html", page)
	return getDocument(url)
}

func getDocument(url string) *goquery.Document {
	var (
		resp *http.Response
		err  error
	)
	for retry := 2; retry >= 0; retry-- {
		resp, err = http.Get(url)
		if err != nil {
			if retry == 0 {
				log.Fatal(err)
			} else {
				continue
			}
		}
		break
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}
