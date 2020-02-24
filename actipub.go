package main

import (
    "fmt"
    "bufio"
        "net/http"
        "io/ioutil"
        "encoding/json"
        "golang.org/x/mobile/asset"
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
    f, err := asset.Open("url.txt")
    fmt.Printf("detson %v %v\n",f, err)
    if err != nil {
        fmt.Println("detson cannot open url.txt")
        return ss
    }
    fmt.Println("detson f ok")
    defer f.Close()
    u := "https://bctpub.duckdns.org/polson/outbox?page=1"
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        // only keep the last url
        u = scanner.Text()
    }
    fmt.Println("detson url: "+u)

    s := getAPNote(u)
    var jsonobj map[string]interface{}
    json.Unmarshal([]byte(s), &jsonobj)
    z := jsonobj["orderedItems"].([]interface{})
    for _,x := range z {
        y := x.(map[string]interface{})
        u := y["object"].(map[string]interface{})
        ss=append(ss, u["content"].(string))
    }
    return ss
}

