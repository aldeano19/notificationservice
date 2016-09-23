package main

import (
	// "net/http"
	// "log"
	"gopkg.in/olivere/elastic.v2"
)


func initClient() (*elastic.Client, error){
	return elastic.NewClient(elastic.SetURL(ELASTIC_HOST))
  // return elastic.NewClient(elastic.SetURL("http://172.17.0.2:9200"))
  
}

func updateProjectOutputDocument(client *elastic.Client, outputDocument *ProjectOutputDocument) (string){

    // tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
  put1, err:= client.Index().
    Index(TCC_INDEX).
    Type(TCC_TYPE_PROJECT_OUTPUT).
    BodyJson(outputDocument).
    Do()
  
  checkForError(err)

  return put1.Id
}

func updateProcessDocument2(client *elastic.Client, siteSpellingDocument *SiteSpellingDocument) (string){

  // tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
  put1, err:= client.Index().
    Index(TCC_INDEX).
    Type(TCC_TYPE_SITE_SPELLING).
    BodyJson(siteSpellingDocument).
    Do()
  
  checkForError(err)

  return put1.Id
}

func updateProcessDocument(client *elastic.Client, brokenLinkDocument *BrokenLinkDocument) (string){

	// tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
  put1, err:= client.Index().
    Index(TCC_INDEX).
    Type(TCC_TYPE_BROKEN_LINKS).
    BodyJson(brokenLinkDocument).
    Do()
  
  checkForError(err)

  return put1.Id
}

func createIndexTestcasecentralIfNotExists(client *elastic.Client) (bool, error){
	// Use the IndexExists service to check if a specified index exists.
  exists, err := client.IndexExists(TCC_INDEX).Do()
  if err != nil{
  	return false, err
  }

  if !exists {
    // Create a new index.
    createIndex, err := client.CreateIndex(TCC_INDEX).Do()
    if err != nil || !createIndex.Acknowledged{
	    return false, err
    }
  }

  return true, nil
}

// func main() {
//   client, err := initClient()
//   if err != nil{
//     log.Println(err)
//   }
//   log.Println(client)
//   // createIndexTestcasecentralIfNotExists(client)
// }