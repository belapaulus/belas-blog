package main

import (
	"github.com/russross/blackfriday/v2"
	"os"
	"encoding/csv"
	"strings"
	"io"
	"time"
	"text/template"
)

type article struct {
	Title, MDFile, HTMLFile string
	Date time.Time
}

func main() {
	articleSlice := getArticles()
	webring := openring()
	makePages(articleSlice, webring)
	makeList(articleSlice, webring)
}

func getArticles() (articleSlice []article) {
	file, err := os.Open("articles.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		title := record[0]
		date, err := time.Parse("2006-01-02", record[1])
		if err != nil {
			panic(err)
		}
		mdFile := record[2]
		htmlFile := date.Format("2006-01-02") + "-" + strings.Split(mdFile, ".")[0] + ".html"
		a := article{
			Title: title,
			Date: date,
			MDFile: mdFile,
			HTMLFile: htmlFile,
		}
		articleSlice = append(articleSlice, a)
	}
	return
}

func makePages(articleSlice []article, webring string) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/page.html")
	if err != nil {
		panic(err)
	}
	for _, article := range articleSlice {
		// TODO: use path.join()
		content, err := os.ReadFile("articles/" + article.MDFile)
		if err != nil {
			// TODO: replace log.fatal with panic
			panic(err)
		}
		file, err := os.OpenFile("htdocs/" + article.HTMLFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		data := struct{ 
			Stylesheets []string
			Title string
			Main string
			Footer string
		}{ 
			Stylesheets: []string{"style.css", "page.css"},
			Title: article.Title,
			Main: string(blackfriday.Run(content)),
			Footer: webring,
		}
		err = tmpl.Execute(file, data)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
}

func makeList(articleSlice []article, webring string) {
	// load template
	tmpl, err := template.New("template").Funcs(template.FuncMap{
			"datef": func(fmt string, t time.Time) string {
				return t.Format(fmt)
			},
		}).
		ParseFiles("templates/base.html", "templates/list.html")
	if err != nil {
		panic(err)
	}
	// execute template
	file, err := os.OpenFile("htdocs/index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data := struct{
		Main []article
		Footer string
		Stylesheets []string
	}{ 
		Stylesheets: []string{"style.css"},
		Main: articleSlice,
		Footer: webring,
	}
	err = tmpl.ExecuteTemplate(file, "base.html", data)
	if err != nil {
		panic(err)
	}
}
