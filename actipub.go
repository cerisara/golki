package main

import (
        "net/http"
        "fmt"
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

func aptest() {
    u := "https://bctpub.duckdns.org/polson/outbox?page=1"
    s := getAPNote(u)
    var jsonobj map[string]interface{}
    json.Unmarshal([]byte(s), &jsonobj)
    z := jsonobj["orderedItems"].([]interface{})
    for _,x := range z {
        y := x.(map[string]interface{})
        u := y["object"].(map[string]interface{})
        fmt.Println(u["content"])
    }
}

