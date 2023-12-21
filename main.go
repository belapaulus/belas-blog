package main

import (
	"github.com/russross/blackfriday/v2"
	"log"
	"os"
	"text/template"
)

const (
	articleDir  = "articles"
	templateDir = "templates"
	outDir      = "htdocs"
)

func main() {
	genArticlePage("test")
}

func genArticlePage(articleFileName string) {
	// load template
	tmpl, err := template.ParseFiles(templateDir + "/base.html")
	if err != nil {
		log.Fatal(err)
	}
	// load article
	article, err := os.ReadFile(articleDir + "/" + articleFileName + ".md")
	if err != nil {
		log.Fatal(err)
	}
	data := struct{ Content string }{string(blackfriday.Run(article))}
	// create output file
	f, err := os.OpenFile(outDir+"/"+articleFileName+".html", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}
}
