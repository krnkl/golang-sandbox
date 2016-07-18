Closing HTTP Response Body


// When you make requests using the standard http library you get a 
// http response variable. If you don't read the response body you 
// still need to close it. Note that you must do it for empty 
// responses too. It's very easy to forget especially for new Go developers.

// Some new Go developers do try to close the response body, 
// but they do it in the wrong place.

package main

import (  
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {  
    resp, err := http.Get("https://api.ipify.org?format=json")
    defer resp.Body.Close()//not ok
    if err != nil {
        fmt.Println(err)
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(string(body))
}
// This code works for successful requests, but if the http request 
// fails the resp variable might be nil, which will cause a runtime panic.

// The most common why to close the response body is by using a defer 
// call after the http response error check.

package main

import (  
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {  
    resp, err := http.Get("https://api.ipify.org?format=json")
    if err != nil {
        fmt.Println(err)
        return
    }

    defer resp.Body.Close()//ok, most of the time :-)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(string(body))
}
// Most of the time when your http request fails the resp variable will 
// be nil and the err variable will be non-nil. However, when you get a 
// redirection failure both variables will be non-nil. This means you 
// can still end up with a leak.

// You can fix this leak by adding a call to close non-nil response 
// bodies in the http response error handling block. Another option is 
// to use one defer call to close response bodies for all failed and 
// successful requests.

package main

import (  
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {  
    resp, err := http.Get("https://api.ipify.org?format=json")
    if resp != nil {
        defer resp.Body.Close()
    }

    if err != nil {
        fmt.Println(err)
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(string(body))
}
// The orignal implementation for resp.Body.Close() also reads and discards 
// the remaining response body data. This ensured that the http connection 
// could be reused for another request if the keepalive http connection behavior 
// is enabled. The latest http client behavior is different. Now it's your 
// responsibility to read and discard the remaining response data. If you 
// don't do it the http connection might be closed instead of being reused. 
// This little gotcha is supposed to be documented in Go 1.5.

// If reusing the http connection is important for your application you 
// might need to add something like this at the end of your response 
// processing logic:

_, err = io.Copy(ioutil.Discard, resp.Body)  
// It will be necessary if you don't read the entire response body right 
// away, which might happen if you are processing json API responses with 
// code like this:

json.NewDecoder(resp.Body).Decode(&data)