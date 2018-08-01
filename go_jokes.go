package main

import (
    "fmt"                 // format
    "io/ioutil"           // "ReadAll" - to read response 
    "net/http"            // "Get" request
    "net/url"             // URL Parsing
    "log"                 // for Logging
    "os"                  // for file handling
    "encoding/json"       // For JSON Marshal/Unmarshal
)

// structure for casting 'Person Data' json object
type Person struct {
    Name    string 
    Surname string 
    Gender  string 
    Region  string
}

// structure for casting 'Random joke' json object
type Joke struct {
    Type    string
    Value struct {
        Id int
        Joke string
        Categories string
    }
}

//Declare global variables
var (
    f *os.File
    err error
)

/**************************************************
Name 		: initLog
Description 	: To initialize the logger object
***************************************************/
func initLog() {

    f, err = os.OpenFile("/var/log/test_server.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)   //
    if err != nil {
        log.Fatal(err)
    }
    
    log.SetOutput(f) // setting file object as the logging stream
    log.Println("Start Logging..")

    
}
/*********************************************************
Name            : serverHTTP
Description     : Handler for HTTP requests to this server
**********************************************************/
func serverHTTP(responseWriter http.ResponseWriter, request *http.Request) {

    urlUinames := "http://uinames.com/api/"
    response, err := http.Get(urlUinames)

    if err != nil {
        log.Fatal(err)
    }

    defer response.Body.Close() //
    log.Println("received response from uinames.com") 

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
 
    responseString := string(responseData)
    log.Println("Response string from names api: ",responseString)
    bytes := []byte(responseString)
    var personObj Person
    json.Unmarshal(bytes, &personObj)
    fname := personObj.Name
    lname := personObj.Surname
    gender := personObj.Gender
    region := personObj.Region
    log.Println("Person name from uinames.com: ",fname, lname) 
    log.Println("Person gender and region from uinames.com: ",gender, region) 
    urlstr := "http://api.icndb.com/jokes/random?firstName=%s&lastName=%s"
    u, _ := url.Parse(urlstr)
    values, _ := url.ParseQuery(u.RawQuery)
    values.Set("firstName",fname)
    values.Set("lastName",lname)
    u.RawQuery = values.Encode()
    response, err = http.Get(u.String())

    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
 
    responseData, err = ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
 
    responseString = string(responseData)
    log.Println("Response string from jokes api: ",responseString)
    bytes = []byte(responseString)
    var jokeObj Joke
    json.Unmarshal(bytes, &jokeObj)
    jokeText := jokeObj.Value.Joke
    fmt.Fprint(responseWriter, jokeText)
}
/*********************************************************
Name            : main
Description     : Main function for HTTP server
**********************************************************/
func main() {
    // Initializing logging mechanism
    initLog()
    defer f.Close()

    // Initializing HTTP handler
    http.HandleFunc("/", serverHTTP)
    err := http.ListenAndServe(":8080", nil) 
    if err != nil {
        log.Fatal(err)
    }
}
