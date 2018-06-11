package src

import (
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"log"
	"net/url"
	"golang.org/x/net/html"
	"strings"
	"gopkg.in/gomail.v2"
	"fmt"
)

type Response struct {
	Message string `json:"message"`
}

func Handler() (Response, error) {
	return Response{
		Message: "Go Serverless v1.0! Your function executed successfully!",
	}, nil
}

func main() {
	lambda.Start(Handler)
}

type Node struct {
	link string
	parent string
	linkText string
	isOutsideLink bool
	statusCode int
}

//kept depth as 4 for the depth-first-search
func crawl(root Node,depth int){
	Urls:=make(chan []Node,100) // buffer size 100
	resp, err := http.Get(string(root.link)) //get the html content
	if err != nil {
		log.Panic("error is ", err)
	}
	go getLinks(resp,root,Urls) //find links for root
	defer resp.Body.Close()
	counter :=1  //counter for links to be crawled in the Urls chan
	isCrawled[root.link] = true //isCrawled is
	for i:=0;i<depth;i++{
		for counter > 0{
			counter --
			next:= <- Urls  //inbound the found links
			//further logic for iterating over the next slice for broken links --> check step 4
		}
	}
	return
}

func getLinks(resp *http.Response, parent Node,Urls chan []Node) {
	wg.Add(1) //add link to waitGroup
	var links = make([]Node,0) //to get a set of links
	z := html.NewTokenizer(resp.Body) // a new tokenizer for the response of the html
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			Urls <- links
			wg.Done()
			return
		case tt == html.StartTagToken:
			t := z.Token()            // taken token
			isAnchor := t.Data == "a" // checking whether it is anchor

			if isAnchor{
				for _, a := range t.Attr { //going through all attributes of anchor i.e., t
					base,err := url.Parse(websiteURL1)
					if err != nil {
						log.Println("err is ",err)
					}
					z.Next()
					t1 := z.Token()

					if a.Key == "href"{
						a.Val = strings.TrimSpace(a.Val)

						if a.Val != parent.link {
							var newNode Node //create a node for the found link
							if !strings.Contains(a.Val,"tel:") && !strings.Contains(a.Val,"mailto:"){
								//parsing the url
								u, err := url.Parse(a.Val)
								if err != nil {
									log.Println(err)
								}else {
									uri := base.ResolveReference(u)
									a.Val = uri.String()
									//checking outsideLink
									if strings.Contains(a.Val, websiteURL1) {
										newNode = Node{a.Val, parent.link,linkText,false, 0}
									} else {
										newNode = Node{a.Val, parent.link,linkText,true, 0}
									}
								}

							}

							if !strings.Contains(newNode.link,"mailto:") && !isCrawled[newNode.link]{
								links = append(links, newNode)
							}
							break
						}
					}
				}
			}
		}
	}
	wg.Done()
	return
}

func sendMail(result string){
	m := gomail.NewMessage()
	m.SetHeader("From","*sender email *")
	m.SetHeader("To", "*receiver email*")
	m.SetHeader("Subject", "["+websiteURL1+"] Broken Links Detected")
	m.SetBody("text/html", result)
	d := gomail.NewDialer("smtp.gmail.com", 25, "*sender email*", "*sender password*")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		log.Panic(err)
	}
}

func doesSomething(){
	next:= <- Urls
	for _, url := range next {
		if _, done := isCrawled[url.link] ; !done {
			if url.link!=""{
				timeout := time.Duration(10 * time.Second)
				client := http.Client{Timeout: timeout}
				resp, err := client.Get(strings.TrimSpace(url.link))
				if resp != nil {
					if resp.StatusCode == 404 {
						brokenLinks = append(brokenLinks, url) //appended to broken links slice
					}else {
						if !url.isOutsideLink && url.linkText=="Link" {
							counter ++
							go getLinks(resp, url, Urls) //crawl the remaining links in chan
							isCrawled[url.link] = true
						}
					}
				}
			}
		}
	}


	for i := 0; i < len(brokenLinks); i++ {
		if brokenLinks[i].link != ""{
			fmt.Printf("Broken link : %s \n", brokenLinks[i].link)
			contentForMail = append (contentForMail,"Link Type: "+brokenLinks[i].linkText+"<br>Link URL: <a href='"+brokenLinks[i].link+"'>"+brokenLinks[i].link+"</a>" + "<br>Source: " + brokenLinks[i].parent+"<br><br>")
			results = append(results, "Link Type: "+brokenLinks[i].linkText+"<br>Link URL: <a href='"+brokenLinks[i].link+"'>"+brokenLinks[i].link+"</a>" + "<br>Source: <a href='"+brokenLinks[i].parent+"'>" + brokenLinks[i].parent+"</a><br><br>")
		}
	}
}