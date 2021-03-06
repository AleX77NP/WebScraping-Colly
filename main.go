package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)


func main() {

	subjects := []string{}

	c := colly.NewCollector()

	err := c.Post("http://moodle.np.ac.rs/login/index.php", map[string]string {"username" : "username", "password": "password"})
	if err != nil {
		log.Fatal(err)
	}

	c.OnResponse(func(t *colly.Response){
		log.Println("response recieved", t.StatusCode)
	})
	c.OnHTML("h3", func(h *colly.HTMLElement) {
		if h.Attr("class") == "categoryname" {
			h.Request.Visit(h.ChildAttr("a", "href") + "&browse=courses&perpage=20&page=0")
		h.Request.Visit(h.ChildAttr("a", "href") + "&browse=courses&perpage=20&page=1")
		}
	})

	c.OnHTML("a", func(h *colly.HTMLElement) {
		//fmt.Printf("%s \n", h.Text)
		h.Request.Visit(h.ChildAttr("div", "class"))
	})

	c.OnHTML("div", func(h *colly.HTMLElement) {
		if h.Attr("class") == "coursename" {
			//fmt.Println(h.Text)
			subjects = append(subjects, h.Text)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://moodle.np.ac.rs/course/index.php?categoryid=28")

	f, err := os.Create("subjects.txt")

    if err != nil {
        log.Fatal(err)
    }

	for _,value := range subjects {
		_, err2 := f.WriteString(value + "\n")
		if err2 != nil {
			log.Fatal("Error while writing...")
		}
	}

	defer fmt.Println("Subjectes fetched from Moodle. Look them up in your text file.")
	defer f.Close()

}