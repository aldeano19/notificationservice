package main

const (
	// elastic index for this app
	TCC_INDEX = "testcasecentral"

	// elastic types for this app
	TCC_TYPE_BROKEN_LINKS = "brokenlinks"
	TCC_TYPE_SITE_SPELLING = "sitespelling"
	TCC_TYPE_PROJECT_OUTPUT = "projectoutput"
	
	// elastic host
	ELASTIC_HOST = "http://elasticsearch:9200"
	
	// process statuses
	STATUS_NORUN = "norun"
	STATUS_RUNNING = "running"
	STATUS_SUCCESS = "success"

	// full date format
	// FULL_DATE_FORMAT = 

	// no error
	NO_ERROR = "None"

	// no ticket
	NO_TICKET = "None"

	// english alphabet and whitespace
	ALPHA = "abcdefghijklmnopqrstuvwxy ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// special characters
	EMPTY_STRING = ""
	NEW_LINE_CHAR = "\n"

	// word collection file
	WORD_COLLECTION = "words.txt"
) 