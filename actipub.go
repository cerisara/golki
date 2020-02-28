package main

import (
    "fmt"
    "html"
    "strconv"
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
    url,_ = strconv.Unquote("\""+url+"\"")

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error reading request. ", err)
    }
    req.Header.Set("Accept", "application/activity+json")
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
}
func NewAPobject(t int) *APobject {
    o := &APobject{t,nil,nil}
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

func debugJson(s string) {
    fmt.Println(s)
    /*
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
    nextURL string
    prevURL string
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
        aa,_ := strconv.Unquote("\""+a+"\"")
        a = aa
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
    next,_ := getJsonVal(s, "next")
    prev,_ := getJsonVal(s, "prev")
    return &APtoots{froms,txts,next,prev}
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
func parseOutbox(s string, k string) *APtoots {
    a,_ := getJsonVal(s, k)
    x := getAPJson(a)
    return parseToots(x)
}

func parseAPjson(s string, k string) *APobject {
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
        fmt.Println("page of toots detected")
        o := NewAPobject(1)
        o.Toots = parseToots(s)
        return o
    }
    i = strings.Index(s, "\"last\"")
    if i>=0 {
        // c'est une outbox
        fmt.Println("outbox detected")
        o := NewAPobject(1)
        o.Toots = parseOutbox(s,k)
        return o
    }
    return nil
}

var curObj *APobject

func GetPosts(u string) []string {
    var ss []string
    s := getAPJson(u)
    curObj = parseAPjson(s,"first")
    if curObj != nil {
        ss = curObj.Strings()
    }
    fmt.Printf("loaded cur page %d\n",len(ss))
    return ss
}

func GetNextPosts() []string {
    var ss []string
    if curObj.Toots.nextURL != "" {
        s := getAPJson(curObj.Toots.nextURL)
        newObj := parseAPjson(s,"next")
        if newObj != nil {
            curObj = newObj
            ss = curObj.Strings()
        }
    }
    fmt.Printf("loaded next page %d\n",len(ss))
    return ss
}
func GetPrevPosts() []string {
    var ss []string
    if curObj.Toots.prevURL != "" {
        s := getAPJson(curObj.Toots.prevURL)
        newObj := parseAPjson(s,"prev")
        if newObj != nil {
            curObj = newObj
            ss = curObj.Strings()
        }
    }
    fmt.Printf("loaded prev page %d\n",len(ss))
    return ss
}


