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
	Title, MDFile, HTMLFile, MDContent, HTMLContent string
	Date time.Time
}

// TODO: use path.join() where applicable

func main() {
	articleSlice := getArticles()
	makePages(articleSlice)
	makeList(articleSlice)
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
		mdContent, err := os.ReadFile("articles/" + mdFile)
		if err != nil {
			panic(err)
		}
		htmlContent := blackfriday.Run(mdContent)
		a := article{
			Title: title,
			Date: date,
			MDFile: mdFile,
			HTMLFile: htmlFile,
			MDContent: string(mdContent),
			HTMLContent: string(htmlContent),
		}
		articleSlice = append(articleSlice, a)
	}
	return
}

func makePages(articleSlice []article) {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
			"datef": func(fmt string, t time.Time) string {
				return t.Format(fmt)
			},
		}).
		ParseFiles("templates/base.html", "templates/page.html")
	if err != nil {
		panic(err)
	}
	for _, a := range articleSlice {
		file, err := os.OpenFile("htdocs/" + a.HTMLFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		data := struct{ 
			Stylesheets []string
			Main article
		}{ 
			Stylesheets: []string{"style.css", "page.css"},
			Main: a,
		}
		err = tmpl.ExecuteTemplate(file, "base.html", data)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
}

func makeList(articleSlice []article) {
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
		Stylesheets []string
	}{ 
		Stylesheets: []string{"style.css"},
		Main: articleSlice,
	}
	err = tmpl.ExecuteTemplate(file, "base.html", data)
	if err != nil {
		panic(err)
	}
}
