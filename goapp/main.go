package main

import (
    "fmt"
    "net/http"
    "net/smtp"
    "io/ioutil"
    // "io"
    "log"
    // "strings"
    // "os"
    // "bytes"
    "strconv"
    "encoding/json"
    // "net/url"
    "strings"
    "regexp"
    // "time"
    // "gopkg.in/olivere/elastic.v2"
    // "reflect"
)

type EmailUser struct {
    Username    string
    Password    string
    EmailServer string
    Port        int
}

type Response struct{
    Sent        bool
    Recipient   string
    Message     string
    Error       string
}

func checkForErrorAndReportToWeb(err error, w http.ResponseWriter){
    if err != nil{
        defer func() {
            if r := recover(); r != nil {
                fmt.Fprint(w, "Bad Request.\n", err.Error(), "Ignore mesage from here on.====")

                messageToEriel := "Error at testcasecentral Notification REST API: " + err.Error() 
                emailSubject := "Error at testcasecentral Notification REST API"
                sendSms("7864288315", messageToEriel)
                sendEmail("eriel@testcasecentral.com", emailSubject, messageToEriel)
            }
        }()
        
        panic(err)
    }
}

func sendEmail(recipient string, subject string, emailBody string) (error){
    emailUser := &EmailUser{"testcasecentral@gmail.com", "s82.ad8aj3,", "smtp.gmail.com", 587}

    msg := []byte("To: " + recipient+ "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" + emailBody + "\r\n")

    auth := smtp.PlainAuth("", emailUser.Username, emailUser.Password, emailUser.EmailServer)

    err := smtp.SendMail(
        emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
        auth,
        emailUser.Username,
        []string{recipient},
        msg)
    return err
}

func sendSms(recipient string, smsBody string) (string, error) {

    encoded := "number="+recipient+"&message="+smsBody

    body := strings.NewReader(encoded)

    req, err := http.NewRequest("POST", "http://textbelt.com/text", body)
    
    if(err != nil){return "", err}

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := http.DefaultClient.Do(req)
    if(err != nil){return "", err}
    defer resp.Body.Close()

    response, err := ioutil.ReadAll(resp.Body)
    if(err != nil){return "", err}

    return string(response), nil
}

func validEmail(email string) (bool){
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
 	return Re.MatchString(email)
}

func isValidDomain(domain string) (bool) {return true}
// MUXes =========
// Test sendEmailHandler with CURL
// curl -k -X POST \
// -d 'recipient={recipient email}' \
// -d 'body={email body}' \
// -d 'subject={optional subject}' \
// http://localhost:8080/email
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
    recipient := r.FormValue("recipient")
    subject := r.FormValue("subject")
    body := r.FormValue("body")

    resp := &Response{}

    if !validEmail(recipient){
    	resp = &Response{false, recipient, body, "Recipient not valid."}
    }else if !(recipient == "" || body == ""){
        if(subject == ""){
            subject = "Hello"
        }
        err := sendEmail(recipient, subject, body)
        checkForErrorAndReportToWeb(err, w)
        resp = &Response{true, recipient, body, "None"}
    }else{
        errorMsg := "Necessary field is empty."
        resp = &Response{false, recipient, body, errorMsg}
    }

    jsonOut, err := json.Marshal(resp)
    checkForErrorAndReportToWeb(err, w)

    fmt.Fprint(w, string(jsonOut))
}

// Test sendSmsHandler with CURL
// curl -k -X POST \
// -d 'recipient={phone number}' \
// -d 'body={sms body}' \
// http://localhost:8080/sms
func sendSmsHandler(w http.ResponseWriter, r *http.Request) {
    recipient := r.FormValue("recipient")
    body := r.FormValue("body")

    resp := &Response{}

    if !(recipient == "" || body == ""){
        textBeltResponse, err := sendSms(recipient, body)
        checkForErrorAndReportToWeb(err, w)

        if(strings.Contains(textBeltResponse, "\"message\": ")){
            startIndex := strings.Index(textBeltResponse, "\"message\": ") + len("\"message\": ")
            myError := strings.Replace(textBeltResponse[startIndex:], "\"", "", -1)
            myError = strings.Replace(myError, "\n", "", -1)
            myError = strings.Replace(myError, "}", "", -1)

            resp = &Response{false, 
                recipient, 
                body, 
                myError}
        }else{
            resp = &Response{true, recipient, body, "None"}
        }
    }else{
        errorMsg := "Necessary field is empty."
        resp = &Response{false, recipient, body, errorMsg}
    }

    jsonOut, err := json.Marshal(resp)
    checkForErrorAndReportToWeb(err, w)

    fmt.Fprint(w, string(jsonOut))
}

func getBrokenLinksHandler(w http.ResponseWriter, r *http.Request){
    domain := r.FormValue("domain")

    response := &GetBrokenLinksResponse{}

    if domain == ""{
        errorMsg := "no domain specified"
        response = responseGetBrokenLinksFailedToStart(domain, errorMsg)
    }else{
        response = responseGetBrokenLinksStarted(domain)
    }

    jsonOut, err:= json.Marshal(response)
    checkForError(err)

    fmt.Fprint(w, string(jsonOut))
}

func brokenLinksReportHandler(w http.ResponseWriter, r *http.Request) {
    ticketId := r.FormValue("ticketid")

    response := &GetBrokenLinksResponse{}

    response = responseGetBrokenLinksProcessReport(ticketId)

    jsonOut, err:= json.Marshal(response)
    checkForError(err)

    fmt.Fprint(w, string(jsonOut))
}

func spellingHandler(w http.ResponseWriter, r *http.Request){
    domain := r.FormValue("domain")

    response := &SiteSpellingResponse{}

    if domain == ""{
        errorMsg := "no domain specified"
        response = responseSiteSpellingCheckFailedToStart(domain, errorMsg)
    }else{
        response = responseSiteSpellingCheckStarted(domain)
    }

    jsonOut, err:= json.Marshal(response)
    checkForError(err)

    fmt.Fprint(w, string(jsonOut))
}

func siteSpellingReportHandler(w http.ResponseWriter, r *http.Request) {
    ticketId := r.FormValue("ticketid")

    response := &SiteSpellingResponse{}

    response = responseSiteSpellingCheckProcessReport(ticketId)

    jsonOut, err:= json.Marshal(response)
    checkForError(err)

    fmt.Fprint(w, string(jsonOut))
}

// This does a brute force formating of the json so that the saving process of the data doesnt have to marshal a huge string into a struct, instead it dumps the json string raw into elasticsearch.
func getProjectOutFileHandler(w http.ResponseWriter, r *http.Request) {
    ticketId := r.FormValue("ticketid")

    response := &ProjectOutputResponse{}

    response = getProjectOutput(ticketId)

    out := "{\n"
    out += " \"SaveTime\":\"" + response.SaveTime + "\",\n"
    out += " \"Project\":" + response.RawOutput + "\n}"

    fmt.Fprint(w, out)
}

func saveProjectOutFileHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    var suites []*json.RawMessage
    var cases []*json.RawMessage
    var steps []*json.RawMessage

    var obj_step map[string]*json.RawMessage
    var obj_case map[string]*json.RawMessage
    var obj_suite map[string]*json.RawMessage
    var project map[string]*json.RawMessage
    
    json.Unmarshal([]byte(r.Form["project"][0]), &project)

    // projectHistoryId, _:= strconv.Atoi(string(*project["project_history_id"]))
    raw_suites := *project["suites"]
    json.Unmarshal(raw_suites, &suites)

    // var raw_suites map[string]*json.RawMessage

    for _, vss := range suites{
        json.Unmarshal([]byte(string(*vss)), &obj_suite)

        raw_cases := *obj_suite["cases"]
        json.Unmarshal(raw_cases, &cases)

        for _, vcs := range cases{
            json.Unmarshal([]byte(string(*vcs)), &obj_case)

            raw_steps := *obj_case["steps"]
            json.Unmarshal(raw_steps, &steps)

            step_struct_list := []*Step{}
            for _, vsteps := range steps{
                // fmt.Fprintln(w, string(*vsteps))
                json.Unmarshal([]byte(string(*vsteps)), &obj_step)

                pid, _ := strconv.Atoi(string(*obj_step["project_id"]))
                phid,_ := strconv.Atoi(string(*obj_step["project_history_id"]))
                tcshid, _ := strconv.Atoi(string(*obj_step["test_case_structure_history_id"]))
                so, _ := strconv.Atoi(string(*obj_step["step_ordinal"]))

                new_step := &Step{pid,
                    phid,
                    tcshid,
                    so,
                    string(*obj_step["case_description"]),
                    string(*obj_step["case_name"]),
                    string(*obj_step["final_status"])}

                step_struct_list = append(step_struct_list, new_step)
            }
        }

        // new_suite := &Suite{string(*obj_suite["project_id"]),
        //     string(*obj_suite["project_history_id"]),
        //     string(*obj_suite["project_structure_history_id"]),
        //     string(*obj_suite["suite_ordinal"]),
        //     string(*obj_suite["suite_id"]),
        //     string(*obj_suite["suite_description"]),
        //     string(*obj_suite["final_status"]),
        //     string(*obj_suite["suite_name"]),
        //     cases}
    }
    
    // fmt.Fprint(w, reflect.TypeOf(suites))

//     for _, obj_suite := range r.Form["suites"]{
//         for _, obj_case := range obj_suite["cases"]{
//             steps := []Step{}
//             for _, obj_step := range obj_case["steps"]{
//                 new_step := &Step{obj_step["project_id"],
//                     obj_step["project_history_id"],
//                     obj_step["test_case_structure_history_id"],
//                     obj_step["step_ordinal"],
//                     obj_step["case_description"],
//                     obj_step["case_name"],
//                     obj_step["final_status"]}
//             }
//             // new_case := &Case{obj_case["project_id"],
//             //     obj_case["project_history_id"],
//             //     obj_case["suite_structure_history_id"],
//             //     obj_case["case_id"],
//             //     obj_case["case_ordinal"],
//             //     obj_case["case_description"],
//             //     obj_case["final_status"],
//             //     steps}
//         }
//         // new_suite := &Suite{obj_suite["project_id"],
//         //     obj_suite["project_history_id"],
//         //     obj_suite["project_structure_history_id"],
//         //     obj_suite["suite_ordinal"],
//         //     obj_suite["suite_id"],
//         //     obj_suite["suite_description"],
//         //     obj_suite["final_status"],
//         //     obj_suite["suite_name"],
//         //     cases}
//     }
//     // var project = &Project{r.Form["project_id"],
//     //     r.Form["project_history_id"],
//     //     r.Form["project_structure_history_id"],
//     //     r.Form["project_name"],
//     //     r.Form["project_description"],
//     //     r.Form["final_status"],
//     //     suites}

//     // fmt.Fprint(w, r.Form)
// }

// func saveProjectOutFileHandler(w http.ResponseWriter, r *http.Request) {

//     data := r.FormValue("project")

//     response := &ProjectOutputResponse{}

//     if data == ""{
//         errorMsg := "No data to save."
//         response = failedToSaveProjectOutput(data, errorMsg)
//     }else{
//         response = saveProjectOutput(data)
//     }

//     jsonOut, err:= json.Marshal(response)
//     checkForError(err)

//     fmt.Fprint(w, string(jsonOut))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<h1>Make a request...</h1>")
}

func main() {
    log.Println("Listenning in port :3001")

    http.HandleFunc("/email", sendEmailHandler)
    http.HandleFunc("/sms", sendSmsHandler)

    http.HandleFunc("/getbrokenlinks", getBrokenLinksHandler)
    http.HandleFunc("/brokenlinksreport", brokenLinksReportHandler)

    http.HandleFunc("/spelling", spellingHandler)
    http.HandleFunc("/spellingreport", siteSpellingReportHandler) 
    
    http.HandleFunc("/saveprojectoutoutput", saveProjectOutFileHandler)
    http.HandleFunc("/getprojectoutoutput", getProjectOutFileHandler)

    // http.HandleFunc("/notifysirqul", notifySirlqulHandler)

    http.HandleFunc("/", defaultHandler)

    http.ListenAndServe(":3001", nil)
}