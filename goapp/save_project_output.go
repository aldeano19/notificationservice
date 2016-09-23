package main

import (
    // "log"
    "time"
)

type Step struct{
    ProjectId                    int
    ProjectHistoryId             int
    TestCaseStructureHistoryId   int
    StepOrdinal                  int
    StepDescription              string
    StepName                     string
    FinalStatus                  string
}

type Case struct{
    ProjectId                    int
    ProjectHistoryId             int
    SuiteStructureHistoryId      int
    CaseId                       int
    CaseOrdinal                  int
    CaseDescription              string
    CaseName                     string
    FinalStatus                  string
    Steps                        []Step
}

type Suite struct{
    ProjectId                    int
    ProjectHistoryId             int
    ProjectStructureHistoryId    int
    SuiteOrdinal                 int
    SuiteId                      int
    SuiteDescription             int
    FinalStatus                  string
    SuiteName                    string
    Cases                        []Case
}

type Project struct{
    ProjectId                    int
    ProjectHistoryId             int
    ProjectStructureHistoryId    int
    ProjectName                  string
    ProjectDescription           string
    FinalStatus                  string
    Suites                       []Suite
}

type ProjectOutputDocument2 struct{
    SaveTime    string
    Project
}


type ProjectOutputResponse struct{
    RawOutput   string
    SaveTime    string
    TicketId    string
    Error       string
}

type ProjectOutputDocument struct{
    RawOutput   string
    SaveTime    string
}

func failedToSaveProjectOutput(rawProjectOutput, errorMsg string) (*ProjectOutputResponse) {

    mTime := time.Now()

    return &ProjectOutputResponse{rawProjectOutput,
        mTime.Format(time.ANSIC), // time format Mon Jan _2 15:04:05 2006
        NO_TICKET,
        errorMsg}
}

// func saveProjectOutput(rawProjectOutput string) (*ProjectOutputResponse) {
    
// }
// func saveProjectOutput(rawProjectOutput string) (*ProjectOutputResponse) {
//     mTime := time.Now()

//     saveOutputDoc :=  &ProjectOutputDocument{rawProjectOutput,
//         mTime.Format(time.ANSIC)}

//     client, err := initClient() // client is of type *elastic.Client
//     checkForError(err)
//     createIndexTestcasecentralIfNotExists(client)

//     ticketId := updateProjectOutputDocument(client, saveOutputDoc)

//     return &ProjectOutputResponse{rawProjectOutput,
//         mTime.Format(time.ANSIC),
//         ticketId,
//         NO_ERROR}
// }
