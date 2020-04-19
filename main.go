package main

import (
	"fmt"
	_ "html"
	_ "io/ioutil"
	_ "net/http"
	_ "regexp"

	parser "github.com/Petrus97/bif-parser/cmd"
	"github.com/Petrus97/bif-parser/morestrings"
)

func main() {
	fmt.Println("Hello")
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
	parser.ReadBIF_V2()
	// // USE A REGULAR EXPRESSION
	// re := regexp.MustCompile("h([[:alpha:]]+)king")
	// fmt.Println(re.FindStringSubmatchIndex("hacking hiking"))
	// fmt.Println(re.FindStringSubmatchIndex("licking"))
	// // REGEX REPLACING
	// re2 := regexp.MustCompile("ise")
	// s := "ize"
	// fmt.Println(re2.ReplaceAllString("realise", s))
	// fmt.Println(re2.ReplaceAllString("organise", s))
	// fmt.Println(re2.ReplaceAllString("analyse", s))
	// // HTML SCRAPING
	// resp, err := http.Get("https://petition.parliament.uk/petitions")
	// if err != nil {
	// 	fmt.Println("error http get")
	// }
	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// // fmt.Println(string(body))
	// if err != nil {
	// 	fmt.Println("http read error")
	// }

	// src := string(body)

	// r, _ := regexp.Compile("\\<h3\\>.*\\</h3\\>") // h3 compile
	// rHTML, _ := regexp.Compile("<[^>]*>")
	// titles := r.FindAllString(src, -1) // return all successive, to get the first put -1

	// for _, title := range titles {
	// 	cleanTitle := rHTML.ReplaceAllString(title, "")
	// 	fmt.Println(html.UnescapeString(cleanTitle))
	// }

}
