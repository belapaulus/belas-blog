package main

import (
	"github.com/russross/blackfriday/v2"
	"os"
	"encoding/csv"
	"strings"
	"io"
	"log"
	"fmt"
	"text/template"
)

type article struct {
	Title, Date, MDFile, HTMLFile string
}


func main() {
	file, err := os.Open("articles.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	var articleList []article
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
		title := record[0]
		date := record[1]
		mdFile := record[2]
		htmlFile := date + "-" + strings.Split(mdFile, ".")[0] + ".html"
		a := article{title, date, mdFile, htmlFile}
		articleList = append(articleList, a)
		makePage(mdFile, htmlFile)
	}
	makeList(articleList)
}

func makePage(mdFile, htmlFile string) {
	// load template
	tmpl, err := template.ParseFiles("templates/base.html", "templates/page.html")
	if err != nil {
		log.Fatal(err)
	}
	// load article
	article, err := os.ReadFile("articles/" + mdFile)
	if err != nil {
		log.Fatal(err)
	}
	// execute template
	file, err := os.OpenFile("htdocs/" + htmlFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data := struct{ Main string }{ string(blackfriday.Run(article)) }
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}
}

func makeList(articleList []article) {
	// load template
	tmpl, err := template.ParseFiles("templates/base.html", "templates/list.html")
	if err != nil {
		log.Fatal(err)
	}
	// execute template
	file, err := os.OpenFile("htdocs/index.html", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data := struct{ Main []article }{ articleList }
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}
}
