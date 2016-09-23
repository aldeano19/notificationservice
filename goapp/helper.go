package main

import (
    "strings"
	"log"
// "fmt"
)


func checkForError(err error) {
	if err != nil{
        defer func() {
            if r := recover(); r != nil {
                log.Println("Bad Request.\n", err.Error())
            }
        }()
        
        panic(err)
    }
}

func filterOutDuplicates(links []string) ([]string){
    var outSlice []string
    for _, link := range links {
        if !stringInSlice(link, outSlice) {
            outSlice = append(outSlice, link)
        }
    }
    return outSlice
}

func filterOutSubPageReferences(links []string) ([]string) {
    var outSlice []string
    for _, link := range links{
        if containsValidDomain(link){
            if strings.Contains(link, "#") {
                link = strings.Split(link, "#")[0]
            }
            outSlice = append(outSlice, link)
        }
    }
    return outSlice
}



func filterOutExternalLinks(links []string, domain string) ([]string) {
    var outSlice []string
    for _, link := range links {
        if strings.Contains(link, domain){
            outSlice = append(outSlice, link)
        }
    }
    return outSlice
}

func stringInSlice(str string, list []string) bool {
 	for _, v := range list {
 		if v == str {
 			return true
 		}
 	}
 	return false
}

func generateTicketId() (int){
	return 0
}

// func main() {
// 	// Create a client and connect to http://192.168.2.10:9201
// 	client, err := elastic.NewClient(elastic.SetURL("http://172.17.0.2:9200/"))
// 	checkForError(err)

// 	exists, err := client.IndexExists("twitter").Do()
// 	checkForError(err)
// 	// if !exists {
// 	// 	// Index does not exist yet.
// 	// 	fmt.Println("No Exists")
// 	// }else{

// 	// }

// 	fmt.Println(exists)



// }