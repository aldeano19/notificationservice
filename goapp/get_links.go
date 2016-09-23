package main

import (
	// "fmt"
	"os/exec"
	"strings"
	"regexp"
	"net/http"
	"time"
	// "log"
	// "reflect"
	"gopkg.in/olivere/elastic.v2"
)


type GetBrokenLinksResponse struct{
    BrokenLinks 	[]string
    Domain      	string
    Status      	string
    StartTime   	string
    LastUpdateTime	string
    Error       	string
    TicketId    	string
    TotalInternalLinks int
}

type BrokenLinkDocument struct{
	BrokenLinks 	[]string
    Domain      	string
    Status      	string
    StartTime   	string
    LastUpdateTime     string
    TotalInternalLinks int
}

func responseGetBrokenLinksStarted(domain string) (*GetBrokenLinksResponse){
	startTime := time.Now()
	var brokenLinks = make([]string, 0)
	status := STATUS_RUNNING
	lastUpdateTime := time.Now()
		
	totalInternalLinks := 0

	// Necesary to wait on db insert bc the document id will be used as the ticketId
	brokenLinkDoc := &BrokenLinkDocument{brokenLinks,
		domain,
		status,
		startTime.Format("15:04:05"),
		lastUpdateTime.Format("15:04:05"),
		totalInternalLinks}

	client, err:= initClient() // client is of type *elastic.Client
	checkForError(err)
	createIndexTestcasecentralIfNotExists(client)

	ticketId := updateProcessDocument(client, brokenLinkDoc)

	go startGetBrokenLinks(domain, ticketId, client)

	return &GetBrokenLinksResponse{brokenLinks,
		domain,
		status,
		startTime.Format("15:04:05"),
		lastUpdateTime.Format("15:04:05"),
		NO_ERROR,
		ticketId,
		totalInternalLinks}
}

func startGetBrokenLinks(sDomain, sTicketId string, sClient *elastic.Client) {
	uniqueLinks := getAllUniqueInternalLinks(sDomain)	
	brokenLinks := getBrokenLinks(uniqueLinks)

	lastUpdateTime := time.Now()

	if brokenLinks == nil{
		brokenLinks = make([]string, 0)
	}
	updateDoc := struct {
		BrokenLinks 		[]string
		TotalInternalLinks 	int
		Status 				string
		LastUpdateTime		string
	}{
		brokenLinks,
		len(uniqueLinks),
		STATUS_SUCCESS,
		lastUpdateTime.Format("15:04:05"),
	}

	updateRequest := elastic.NewBulkUpdateRequest().Index(TCC_INDEX).Type(TCC_TYPE_BROKEN_LINKS).Id(sTicketId).Doc(&updateDoc)
	bulkRequest := sClient.Bulk()
	bulkRequest = bulkRequest.Add(updateRequest)

	bulkResponse, err := bulkRequest.Do()
	if err != nil{
		panic(err)
	}

	if bulkResponse == nil{
		panic("No Update response.")
	}

}

func responseGetBrokenLinksFailedToStart(domain string, mError string) (*GetBrokenLinksResponse){
	brokenLinks := make([]string, 0)
	status := STATUS_NORUN
	startTime := time.Now()
	lastUpdateTime := time.Now()
	ticketId := "noid"

	return &GetBrokenLinksResponse{brokenLinks,
		domain,
		status,
		startTime.Format("15:04:05"),
		lastUpdateTime.Format("15:04:05"),
		mError,
		ticketId,
		0}
}

func containsValidDomain(link string) (bool){
    Re := regexp.MustCompile(`[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}`)
    return Re.MatchString(link)
}



func getUniqueInternalLinks(page string) ([]string){
	
	out, err := exec.Command("lynx", page, "-dump", "-listonly", "-nonumbers").CombinedOutput()
	if err != nil{
		return nil
	}

	result := string(out)

	lines := strings.Split(result,"\n")
	
	internalLinks := filterOutSubPageReferences(lines)

	return filterOutDuplicates(internalLinks)
}

func getNewLinks(newSlice, oldSlice []string) ([]string){
	for _, item := range newSlice{
		if !stringInSlice(item, oldSlice){
			oldSlice = append(oldSlice, item)
		}
	}
	return oldSlice
}

func getAllUniqueInternalLinks(page string) ([]string) {
	var collectedLinks []string
	homePageLinks := getUniqueInternalLinks(page)

	collectedLinks = getNewLinks(homePageLinks, collectedLinks)
	collectedLinks = filterOutExternalLinks(collectedLinks, page)
	// oldLength := 0
	// newLength := len(collectedLinks)

	return collectedLinks
}

func getLinkResponseStatus(url string) (int){
	resp, err := http.Get(url)
	if err != nil{
		return -1
	}
	return resp.StatusCode
}

func isBrokenLink(url string) (bool){
	return getLinkResponseStatus(url) != 200
}

func getBrokenLinks(urls []string) ([]string){
	var broken []string
	for _, url := range urls{
		if isBrokenLink(url){
			broken = append(broken, url)	
		}
	}
	return broken
}

// func main() {

// 	uniqueLinks := getAllUniqueInternalLinks("testcasecentral.com")
	
// 	brokenLinks := getBrokenLinks(uniqueLinks)
// 	// fmt.Println(uniqueLinks)
// 	for _, i := range brokenLinks{
// 		fmt.Println(i)
// 	}
// }
