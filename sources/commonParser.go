package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Byron/core"
	"github.com/Byron/utils"
	"github.com/ttacon/chalk"
)

type Source struct {
	SourceName           string
	UrlREGEX             string
	IdREGEX              string
	SearchREGEX          string
	DownloadUrlREGEX     string
	TitleREGEX           string
	IsbnREGEX            string
	YearREGEX            string
	PublisherREGEX       string
	AuthorREGEX          string
	ExtensionREGEX       string
	PageREGEX            string
	LanguageREGEX        string
	SizeREGEX            string
	TimeREGEX            string
	CompletePageUrl      string
	IncompleteArticleUrl string
	AllUrls              []string
	Search               string
}

func (s *Source) GetArticles() {
	r, _ := regexp.Compile(s.UrlREGEX)
	processed := 0

	for i := 1; i < 2; i++ {
		time.Sleep(2 * time.Second)
		resp, err := http.Get(s.CompletePageUrl + strconv.Itoa(i))
		if err != nil {
			log.Println(err)
		}

		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		htmlFormat := string(html)

		if !core.ErrorsHandling(htmlFormat) {
			matches := r.FindAllStringSubmatch(htmlFormat, -1)
			fmt.Println(chalk.Green.Color("Processing page " + strconv.Itoa(i)))

			if len(matches) < 1 {
				break
			}

			for _, m := range matches {
				fmt.Println(chalk.Green.Color("Saving " + m[1]))
				s.AllUrls = append(s.AllUrls, s.IncompleteArticleUrl+m[1])
				processed++
			}

		} else {
			fmt.Println(chalk.Magenta.Color("Given 503. waiting to reconnect"))
			time.Sleep(10 * time.Second)
		}
		resp.Body.Close()
	}
	s.ProcessArticles()
}

func (s *Source) ProcessArticles() {
	fmt.Println(chalk.Green.Color("Start processing Articles.."))

	processed := 0

	//randomize urls processing
	//rand.Seed(time.Now().UnixNano())
	//rand.Shuffle(len(s.AllUrls), func(i, j int) { s.AllUrls[i], s.AllUrls[j] = s.AllUrls[j], s.AllUrls[i] })

	for _, u := range s.AllUrls {
		time.Sleep(2 * time.Second)
		resp, err := http.Get(u)
		if err != nil {
			log.Println(err)
		}

		articleHtml, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		articleHtmlFormat := string(articleHtml)

		if !core.ErrorsHandling(articleHtmlFormat) {

			newArticle := core.Article{
				SourceName: s.SourceName,
				Url:        u,
				Search:     s.Search,
			}

			log.Println("Article:", u)
			if RegexSet(s.TitleREGEX) {
				ArticleTitle, err := regexp.Compile(s.TitleREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Title = ArticleTitle.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.AuthorREGEX) {
				ArticleAuthors, err := regexp.Compile(s.AuthorREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Author = ArticleAuthors.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.PublisherREGEX) {
				ArticlePublisher, err := regexp.Compile(s.PublisherREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Publisher = ArticlePublisher.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.YearREGEX) {
				ArticleYear, err := regexp.Compile(s.YearREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Year = ArticleYear.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.LanguageREGEX) {
				ArticleLang, err := regexp.Compile(s.LanguageREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Language = ArticleLang.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.IsbnREGEX) {
				ArticleIsbn, err := regexp.Compile(s.IsbnREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Isbn = ArticleIsbn.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.TimeREGEX) {
				ArticleTime, err := regexp.Compile(s.TimeREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Time = ArticleTime.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.IdREGEX) {
				ArticleId, err := regexp.Compile(s.IdREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Id = ArticleId.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.SizeREGEX) {
				ArticleSize, err := regexp.Compile(s.SizeREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Size = ArticleSize.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.PageREGEX) {
				ArticlePages, err := regexp.Compile(s.PageREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Page = ArticlePages.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.ExtensionREGEX) {
				ArticleExtension, err := regexp.Compile(s.ExtensionREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.Extension = ArticleExtension.FindStringSubmatch(articleHtmlFormat)[1]
			}
			if RegexSet(s.DownloadUrlREGEX) {
				ArticleDownload, err := regexp.Compile(s.DownloadUrlREGEX)
				if err != nil {
					log.Println(err)
				}
				newArticle.DownloadUrl = ArticleDownload.FindStringSubmatch(articleHtmlFormat)[1]
			}

			/*
				Append and download because it's new
			*/

			AllArticles := s.ReadArticles(utils.GetMD5Hash(s.Search))
			newArticleFormatted := newArticle.FormatNewArticle()

			duplicated := 0
			for i := 0; i < len(AllArticles); i++ {

				if AllArticles[i].Url == newArticleFormatted.Url {
					duplicated = 1
					break
				}
			}

			if duplicated == 0 {

				AllArticlesUpdated := s.ReadArticles(utils.GetMD5Hash(s.Search))
				AllArticlesUpdated = append(AllArticlesUpdated, *newArticleFormatted)
				core.WriteInFile(utils.GetMD5Hash(s.Search), AllArticlesUpdated)

				/*
					Display relevant information about the new Document
				*/
				newArticle.DisplayInformation()
				fmt.Println(chalk.Green.Color("Added correctly: " + newArticle.Title))
				processed++
				fmt.Println(chalk.Magenta.Color("Processed: " + strconv.Itoa(processed)))

				/*
					After that, download the file and save it
				*/
				// core.FileDownload(
				// 	ArticleDownload.FindStringSubmatch(articleHtmlFormat)[1],  //Download url
				// 	ArticleId.FindStringSubmatch(articleHtmlFormat)[1],        //ID
				// 	ArticleExtension.FindStringSubmatch(articleHtmlFormat)[1], //extension ex. .pdf
				// )

			} else {
				fmt.Println(chalk.Red.Color("This article already exists, nothing to do here"))
			}

		} else {

			fmt.Println(chalk.Magenta.Color("Given 503. waiting to reconnect"))
			time.Sleep(5 * time.Second)
		}
		resp.Body.Close()
	}

	fmt.Println(chalk.Green.Color("All the documents were Downloaded :) "))
}

func (s *Source) ReadArticles(inventory string) []core.Article {
	var Articles []core.Article
	jsonFile, err := os.Open("Inventory/" + inventory + ".json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	fileValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(fileValue, &Articles)
	return Articles
}

func RegexSet(regex string) bool {
	return regex != ""
}
