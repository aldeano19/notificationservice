package main

import(
	// "log"
	"encoding/json"
	"time"
	// "reflect"
	"strings"
)

func getProjectOutput(ticketId string)(*ProjectOutputResponse) {
    client, err:= initClient() // client is of type *elastic.Client
    checkForError(err)

    getResult, err := client.Get().Index(TCC_INDEX).Type(TCC_TYPE_PROJECT_OUTPUT).Id(ticketId).Do()

    res := ProjectOutputDocument{}
    json.Unmarshal(*getResult.Source, &res)

    return &ProjectOutputResponse{res.RawOutput,
        res.SaveTime,
        ticketId,
        NO_ERROR}
}


func responseGetBrokenLinksProcessReport(ticketId string)(*GetBrokenLinksResponse){

	client, err:= initClient() // client is of type *elastic.Client
	checkForError(err)

	getResult, err := client.Get().Index(TCC_INDEX).Type(TCC_TYPE_BROKEN_LINKS).Id(ticketId).Do()

	res := GetBrokenLinksResponse{} // TODO: This might be a GetBrokenLinksDocument instead of GetBrokenLinksResponse
    json.Unmarshal(*getResult.Source, &res)

	if strings.EqualFold(res.StartTime, res.LastUpdateTime){
		res.LastUpdateTime = time.Now().Format("15:04:05")
	}
	// return nil
	return &GetBrokenLinksResponse{res.BrokenLinks,
		res.Domain,
		res.Status,
		res.StartTime,
		res.LastUpdateTime,
		NO_ERROR,
		ticketId,
		res.TotalInternalLinks}
}


// TODO: (bug) last update time did not update on startSiteSpellCheck()
func responseSiteSpellingCheckProcessReport(ticketId string)(*SiteSpellingResponse){

	client, err:= initClient() // client is of type *elastic.Client
	checkForError(err)

	getResult, err := client.Get().Index(TCC_INDEX).Type(TCC_TYPE_SITE_SPELLING).Id(ticketId).Do()

	res := SiteSpellingResponse{}
    json.Unmarshal(*getResult.Source, &res)

	if strings.EqualFold(res.StartTime, res.LastUpdateTime){
		// log.Println("yes")
		res.LastUpdateTime = time.Now().Format("15:04:05")
	}
	// return nil
	return &SiteSpellingResponse{res.Misspelled,
		res.Domain,
		res.Status,
		res.StartTime,
		res.LastUpdateTime,
		NO_ERROR,
		ticketId,
		res.TotalUniqueWords}
}

