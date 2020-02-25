package main

import (
    "fmt"
    "html"
    "time"
    "github.com/grokify/html-strip-tags-go"
    "strings"
        "net/http"
        "io/ioutil"
        //"encoding/json"
)

func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

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
func (o *APobject) Strings() []string {
    switch o.Typ {
    case 0: return o.User.Strings()
    case 1: return o.Toots.Strings()
    case 2: return nil
default: return nil
    }
}

func getJsonVal(s string, k string) (string, int) {
    i := strings.Index(s, "\""+k+"\"")
    if i<0 {return "",0}
    j := strings.Index(s[i:],":")
    if j<0 {return "",0}
    a := strings.Index(s[i+j:],"\"")
    if a<0 {return "",0}
    a++
    b := strings.Index(s[i+j+a:],"\"")
    if b<0 {return "",0}
    return s[i+j+a:i+j+a+b], i+j+a+b
}

type APuser struct {
    outbox string
    name string
}
func (o *APuser) Strings() []string {
    var r = []string{"User "+o.name}
    return r
}
func parseUser(s string) *APuser {
    outbox,_ := getJsonVal(s, "outbox")
    name,_ := getJsonVal(s, "name")
    if outbox != "" && name != "" {
        fmt.Println("outbox "+outbox)
        return &APuser{outbox,name}
    }
    return nil
}

type APtoots struct {
    from []string
    txt  []string
}
func parseToots(s string) *APtoots {
    x := s
    var froms []string
    var txts []string
    for ;; {
        a,i := getJsonVal(x, "content")
        if a=="" {break}
        b,j := getJsonVal(x, "attributedTo")
        if b=="" {break}
        stmp := strings.Split(b,"/")
        b = stmp[len(stmp)-1]
        a = strings.Replace(a,"\\u003c","<",-1)
        a = strings.Replace(a,"\\u003e",">",-1)
        a = strings.Replace(a,"\\u0026","&",-1)
        a = strip.StripTags(a)
        a = html.UnescapeString(a)
        // do not handle retoots
        if a!="" {
            froms = append(froms,b)
            txts = append(txts,a)
        }
        k := Max(i,j)
        x = x[k:]
    }
    return &APtoots{froms,txts}
}
func (o *APtoots) Strings() []string {
    var s []string
    for i:=0;i<len(o.from);i++ {
        s = append(s,o.from[i]+": "+o.txt[i])
    }
    return s
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
    // u = "https://bctpub.duckdns.org/polson"
    u = "https://mastodon.etalab.gouv.fr/@cerisara"
    // u = "https://mastodon.etalab.gouv.fr/@cerisara/103514562679577450"
    u = "https://mastodon.etalab.gouv.fr/users/cerisara/outbox?page=true"

    fmt.Println("detson url: "+u)

    s := getAPJson(u)
    x := parseAPjson(s)
    if x != nil {
        ss = x.Strings()
    }

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

