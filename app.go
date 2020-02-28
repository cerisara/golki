package main

import (
    "fmt"
    "fyne.io/fyne"
    "fyne.io/fyne/theme"
    "fyne.io/fyne/layout"
    "fyne.io/fyne/app"
    "fyne.io/fyne/widget"
)

const npaps = 10
var papers [npaps]*widget.Box
var blikes [npaps]*widget.Button
var isliked [npaps]bool
var days = [...]string {"Lundi","Mardi","Mercredi","Jeudi","Vendredi"}
var curday = 0
var tg *TextGroup

func paperchat(i int) {
  fmt.Printf("paper %d\n",i)
  /*
  papers[i].Disable()
  go func() {
      time.Sleep(time.Second)
      papers[i].Enable()
  }()
  */
}

func paperlike(i int) {
    if isliked[i] {
        blikes[i].SetIcon(likesvg)
        isliked[i]=false
    } else {
        blikes[i].SetIcon(likeredsvg)
        isliked[i]=true
    }
}

func menufct(bclicked int) {
    var users = []string {
        "cerisara@etalab",
        "rigel@olki-social",
    }
    var urls = []string {
        "https://mastodon.etalab.gouv.fr/users/cerisara/outbox",
        "https://olki-social.loria.fr/federation/actor/rigelk/outbox?page=1",
    }

    if bclicked == -1 {
        tg.SetTexts(users)
    } else {
        posts := GetPosts(urls[bclicked])
        tg.SetTexts(posts)
    }
}

func main() {
    app := app.New()
    w := app.NewWindow("GOLKi")
    glob := widget.NewVBox()

    titleday := &widget.Label{Text: days[curday], Alignment: fyne.TextAlignCenter}

    decday := func(next bool) {
        if next {
            curday++
            if curday>=len(days) {curday=0}
        } else {
            curday--
            if curday<0 {curday=len(days)-1}
        }
        titleday.SetText(days[curday])
    }

    prevday := widget.NewButton("",func(){decday(false)})
    prevday.SetIcon(theme.NavigateBackIcon())
    nextday := widget.NewButton("",func(){decday(true)})
    nextday.SetIcon(theme.NavigateNextIcon())
    titlebar := widget.NewHBox(prevday,layout.NewSpacer(),titleday,layout.NewSpacer(),nextday)
    glob.Append(titlebar)


    news := widget.NewGroup("Fediverse")
    news1 := &widget.Label{Text: "Read mode", Alignment: fyne.TextAlignCenter}
    news.Append(news1)
    // news2 := &widget.Label{Text: "news 2", Alignment: fyne.TextAlignCenter}
    // news.Append(news2)
    glob.Append(news)

    // scroller is not fluid enough
    // paps := widget.NewGroupWithScroller("Papers")

    var txts = []string{}
    txts = append(txts,"Click on menu button bellow to choose a User thread")
    tg = NewTextGroup(txts, menufct)

    /*
    paps := widget.NewGroup("Papers")
    shownpaps := 2
    for i:=0;i<shownpaps;i++ {
        j := i
        var s string
        if false && i<len(titles) {
            s = titles[i]
        } else {
            s = "pap\ner "+strconv.Itoa(i)
        }
        // lab := &widget.Label{Text: s, Alignment: fyne.TextAlignCenter}
        lab := &widget.Label{Text: s}
        blikes[i] = widget.NewButton("", func() {paperlike(j)})
        blikes[i].SetIcon(likesvg)
        isliked[i]=false
        bchat := widget.NewButton("", func() {paperchat(j)})
        bchat.SetIcon(theme.VisibilityIcon())
        papers[i] = widget.NewHBox(blikes[i],layout.NewSpacer(),lab,layout.NewSpacer(),bchat)
        paps.Append(papers[i])
    }
*/

    /*
    t := widget.NewEntry()
    t.OnCursorChanged = func() {
        fmt.Println("ttyypass")
    }
    */

    content := fyne.NewContainerWithLayout(layout.NewBorderLayout(glob,nil,nil,nil),glob,tg)
    w.SetContent(content)
    // tg.Appwin = content

    /*
    li := canvas.NewLine(color.Black)
    content.AddObject(li)

    go func() {
        time.Sleep(4*time.Second)
        for i:=0;i<len(content.Objects);i++ {
            fmt.Printf("object %d %T %v %v\n",i,content.Objects[i],content.Objects[i],content.Objects[i].Position())
        }
    }()
    */

    w.ShowAndRun()
}

