package main

import (
	"github.com/russross/blackfriday/v2"
	"os"
	"encoding/csv"
	"strings"
	"io"
	"text/template"
)

type article struct {
	Title, Date, MDFile, HTMLFile string
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
		date := record[1]
		mdFile := record[2]
		htmlFile := date + "-" + strings.Split(mdFile, ".")[0] + ".html"
		a := article{title, date, mdFile, htmlFile}
		articleSlice = append(articleSlice, a)
	}
	return
}

func makePages(articleSlice []article, webring string) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/page.html", "templates/footer.html")
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
		file, err := os.OpenFile("htdocs/" + article.HTMLFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		data := struct{ 
			Main string
			Footer string
		}{ 
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
	tmpl, err := template.ParseFiles("templates/base.html", "templates/list.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	// execute template
	file, err := os.OpenFile("htdocs/index.html", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data := struct{
		Main []article
		Footer string
	}{ 
		Main: articleSlice,
		Footer: webring,
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}
