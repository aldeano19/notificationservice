package main

import(
    "strings"
    "log"
    "io/ioutil"
    "os/exec"
    "regexp"
    "time"
    "gopkg.in/olivere/elastic.v2"
)

type SiteSpellingDocument struct{
    Misspelled          []string
    Domain              string
    Status              string
    StartTime           string
    LastUpdateTime      string
    TotalUniqueWords    int
}

type SiteSpellingResponse struct{
    Misspelled          []string
    Domain              string
    Status              string
    StartTime           string
    LastUpdateTime      string
    Error               string
    TicketId            string
    TotalUniqueWords    int
}

func startSiteSpellCheck(sDomain string, sTicketId string, sClient *elastic.Client) {

    lastUpdateTime := time.Now()
    
    words_to_check := fetch_page_alpha(sDomain)
    words_to_check = filterOutDuplicates(words_to_check)
    incorrect_words := make([]string, 0)
    low_word := ""
    for _, word := range words_to_check{
        low_word = strings.ToLower(word)
        if incorrectly_spelled(low_word){
            incorrect_words = append(incorrect_words, word)
        }
    }

    updateDoc := struct {
        Misspelled          []string
        TotalUniqueWords    int
        Status              string
        LastUpdateTime      string
    }{
        incorrect_words,
        len(words_to_check),
        STATUS_SUCCESS,
        lastUpdateTime.Format("15:04:05"),
    }

    updateRequest := elastic.NewBulkUpdateRequest().Index(TCC_INDEX).Type(TCC_TYPE_SITE_SPELLING).Id(sTicketId).Doc(&updateDoc)

    bulkRequest := sClient.Bulk()

    bulkRequest = bulkRequest.Add(updateRequest)

    bulkResponse, err := bulkRequest.Do()
    if err != nil{
        log.Println(err.Error())
    }

    if bulkResponse == nil{
        panic("No Update response.")
    }
}

func responseSiteSpellingCheckStarted(domain string) (*SiteSpellingResponse){
    startTime := time.Now()
    var misspelled = make([]string, 0)
    status := STATUS_RUNNING
    lastUpdateTime := time.Now()

    totalUniqueWords := 0

    siteSpellingDoc := &SiteSpellingDocument{misspelled,
        domain,
        status,
        startTime.Format("15:04:05"),
        lastUpdateTime.Format("15:04:05"),
        totalUniqueWords}

    client, err := initClient() // client is of type *elastic.Client
    checkForError(err)
    createIndexTestcasecentralIfNotExists(client)

    ticketId := updateProcessDocument2(client, siteSpellingDoc)

    log.Println("ticketId : " + ticketId)

    go startSiteSpellCheck(domain, ticketId, client)

    return &SiteSpellingResponse{misspelled,
        domain,
        status,
        startTime.Format("15:04:05"),
        lastUpdateTime.Format("15:04:05"),
        NO_ERROR,
        ticketId,
        totalUniqueWords}
}

func responseSiteSpellingCheckFailedToStart(domain string, mError string) (*SiteSpellingResponse){
    misspelled := make([]string, 0)
    status := STATUS_NORUN
    startTime := time.Now()
    lastUpdateTime := time.Now()
    ticketId := "noid"

    return &SiteSpellingResponse{misspelled,
        domain,
        status,
        startTime.Format("15:04:05"),
        lastUpdateTime.Format("15:04:05"),
        mError,
        ticketId,
        0}
}

// func startSiteSpellingCheck() {
    
// }

func get_lowercase_collection_of_words() []string{
    byte_contents, err:= ioutil.ReadFile(WORD_COLLECTION)
    if err != nil{
        panic(err)
    }

    words := strings.Split(
        strings.ToLower(
            strings.Replace(
                string(byte_contents), "\n", " ", -1)), " ")


    return words
}

func is_valid_word(word string) bool {
    Re := regexp.MustCompile(`^[A-Za-z\s]+$`)
    return Re.MatchString(word)
}


func fetch_page_alpha(url string) []string{
    lynx_cmd := "lynx -dump -nonumbers " + url + " | grep -v 'http'"

    out, err := exec.Command("/bin/sh", "-c", lynx_cmd).CombinedOutput()
    if err != nil{panic(err)}

    words := strings.Split(
        strings.Join(
            strings.Fields(string(out))," "), " ")

    valid_words := make([]string, 0)

    special_chars := []string{".", ",", ":", ";"}
    // special_chars := []string{""}

    for _,word := range words{
        for _,char := range special_chars{
            word = strings.Replace(word, char, "", 1)
        }

        if is_valid_word(word){
            valid_words = append(valid_words, word)
        }
    }

    return valid_words
}

// binary search of a word in slice
func search_word(word string, slice []string) int{
    lo := 0
    hi := len(slice)-1
    mid := 0

    for lo <= hi{
        mid = lo + (hi-lo)/2
        // log.Println(lo, mid, hi)

        

        if word == slice[mid]{
            return mid
        }else if word < slice[mid]{
            hi = mid-1
        }else{
            lo = mid+1
        }
    }

    return -1
}

// replaces 'repacement' for 'termination' in 'word' and searches again 
func incorrect_with_termination(word string, termination string, replacement string) bool{
    word_collection := get_lowercase_collection_of_words()
    if word[len(word)-len(termination):] == termination{
        new_word := strings.Replace(word, termination, replacement, 1)
        return search_word(new_word, word_collection) == -1
    }else{
        return true
    }
}

func incorrectly_spelled(word string) bool{
    word_collection := get_lowercase_collection_of_words()
    if search_word(word, word_collection) == -1{
        // Add check for different terminations as needed.
        check := incorrect_with_termination(word, "ies", "y")
        if check{
            check = incorrect_with_termination(word, "s", "")
        }
        return check
    }
    return false
}
