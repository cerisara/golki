package main

import (
    "fmt"
    "time"
    "strings"
        "net/http"
        "io/ioutil"
        //"encoding/json"
)

func getAPJson(url string) string {

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error reading request. ", err)
    }
    req.Header.Set("Accept", "application/json")
    client := &http.Client{Timeout: time.Second * 10}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error reading response. ", err)
    }

    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body)
}

type APobject struct {
    Typ int // 0: user, 1: toot, 2: outbox
    User *APuser
    Toots *APtoots
    Outbox *APoutbox
}
func NewAPobject(t int) *APobject {
    o := &APobject{t,nil,nil,nil}
    return o
}
func (o *APobject) String() string {
    switch o.Typ {
    case 0: return o.User.String()
    case 1: return ""
    case 2: return ""
default: return ""
    }
}

func getJsonVal(s string, k string) string {
    i := strings.Index(s, "\""+k+"\"")
    if i<0 {return ""}
    j := strings.Index(s[i:],":")
    if j<0 {return ""}
    a := strings.Index(s[i+j:],"\"")
    if a<0 {return ""}
    a++
    b := strings.Index(s[i+j+a:],"\"")
    if b<0 {return ""}
    return s[i+j+a:i+j+a+b]
}

type APuser struct {
    outbox string
    name string
}
func (o *APuser) String() string {
    return "User "+o.name
}
func parseUser(s string) *APuser {
    outbox := getJsonVal(s, "outbox")
    name := getJsonVal(s, "name")
    if outbox != "" && name != "" {
        return &APuser{outbox,name}
    }
    return nil
}

type APtoots struct {
}
func parseToots(s string) *APtoots {
    return nil
}

type APoutbox struct {
}
func parseOutbox(s string) *APoutbox {
    return nil
}

func parseAPjson(s string) *APobject {
    i := strings.Index(s, "\"outbox\"")
    if i>=0 {
        fmt.Println("user detected")
        // c'est un user
        o := NewAPobject(0)
        o.User = parseUser(s)
        return o
    }
    i = strings.Index(s, "\"content\"")
    if i>=0 {
        // c'est un ou plusieurs toots
        o := NewAPobject(1)
        o.Toots = parseToots(s)
        return o
    }
    i = strings.Index(s, "\"last\"")
    if i>=0 {
        // c'est une outbox
        o := NewAPobject(2)
        o.Outbox = parseOutbox(s)
        return o
    }
    return nil
}

func aptest() []string {
    fmt.Println("detson in aptest")
    var ss []string
    u := "https://bctpub.duckdns.org/polson/outbox?page=1"
    u = "https://bctpub.duckdns.org/polson"
    u = "https://mastodon.etalab.gouv.fr/@cerisara"
    fmt.Println("detson url: "+u)

    s := getAPJson(u)
    x := parseAPjson(s)
    if x != nil {ss = append(ss,x.String())}

    /*
    ** maniere exacte de parser le json

    var jsonobj map[string]interface{}
    json.Unmarshal([]byte(s), &jsonobj)

    z := jsonobj["orderedItems"].([]interface{})
    for _,x := range z {
        y := x.(map[string]interface{})
        u := y["object"].(map[string]interface{})
        msg := u["content"].(string)
        // date := u["published"].(string)
        // layout := "2006-01-02T15:04:05.000Z"
        // date,_ = time.Parse(layout,date)
        from := u["attributedTo"].(string)
        stmp := strings.Split(from,"/")
        from = stmp[len(stmp)-1]
        oneline := from+": "+msg
        ss=append(ss, oneline)
    }
    */
    return ss
}

