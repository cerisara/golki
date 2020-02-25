package main

import (
    "fmt"
    "strings"
        "net/http"
        "io/ioutil"
        "encoding/json"
)

func getAPNote(url string) string {
    resp, err := http.Get(url)
    if err != nil {
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body)
}

func aptest() []string {
    fmt.Println("detson in aptest")
    var ss []string
    u := "https://bctpub.duckdns.org/polson/outbox?page=1"

    /*
    * this does not work on android, may be because fyne.io uses a vendor/ dir, which creates 2 contexts in android and thus no asset manager
    f, err := asset.Open("url.txt")
    fmt.Printf("detson %v %v\n",f, err)
    if err != nil {
        fmt.Println("detson cannot open url.txt")
        return ss
    }
    fmt.Println("detson f ok")
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        // only keep the last url
        u = scanner.Text()
    }
    */
    fmt.Println("detson url: "+u)

    s := getAPNote(u)
    var jsonobj map[string]interface{}
    json.Unmarshal([]byte(s), &jsonobj)

    z := jsonobj["orderedItems"].([]interface{})
    for _,x := range z {
        y := x.(map[string]interface{})
        u := y["object"].(map[string]interface{})
        msg := u["content"].(string)
        /*
        date := u["published"].(string)
        layout := "2006-01-02T15:04:05.000Z"
        date,_ = time.Parse(layout,date)
        */
        from := u["attributedTo"].(string)
        stmp := strings.Split(from,"/")
        from = stmp[len(stmp)-1]
        oneline := from+": "+msg
        ss=append(ss, oneline)
    }
    return ss
}

