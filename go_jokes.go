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

/**************************************************
Name 		: init_log
Description 	: To initialize the logger object
***************************************************/
func init_log() {

    f, err := os.OpenFile("/var/log/test_server.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)   //
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    
    log.SetOutput(f) // setting file object as the logging stream
    log.Println("Start Logging..")

    
}
/*********************************************************
Name            : server_HTTP
Description     : Handler for HTTP requests to this server
**********************************************************/
func server_HTTP(responseWriter http.ResponseWriter, request *http.Request) {

    url_uinames := "http://uinames.com/api/"
    response, err := http.Get(url_uinames)

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
    bytes := []byte(responseString)
    var person1 Person
    json.Unmarshal(bytes, &person1)
	fname := person1.Name
	lname := person1.Surname
        log.Println("Person name from uinames.com: %s %s",fname,lname) 
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
    bytes = []byte(responseString)
    var joke1 Joke
    json.Unmarshal(bytes, &joke1)
    joke_text := joke1.Value.Joke
    fmt.Fprint(responseWriter, joke_text)
}
/*********************************************************
Name            : main
Description     : Main function for HTTP server
**********************************************************/
func main() {
    // Initializing logging mechanism
    init_log()









    // Initializing HTTP handler
    http.HandleFunc("/", server_HTTP)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
