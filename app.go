package main

import (
    "fmt"
//    "time"
    "strconv"
    "fyne.io/fyne"
    "fyne.io/fyne/theme"
    "fyne.io/fyne/layout"
    "fyne.io/fyne/app"
    "fyne.io/fyne/widget"
)

const npaps = 10
var papers [npaps]*widget.Box
var days = [...]string {"Lundi","Mardi","Mercredi","Jeudi","Vendredi"}
var curday = 0

var likesvg = &fyne.StaticResource{
        StaticName: "like.svg",
        StaticContent: []byte{
60,63,120,109,108,32,118,101,114,115,105,111,110,61,34,49,46,48,34,32,101,110,99,111,100,105,110,103,61,34,117,116,102,45,56,34,63,62,13,10,60,33,45,45,32,83,118,103,32,86,101,99,116,111,114,32,73,99,111,110,115,32,58,32,104,116,116,112,58,47,47,119,119,119,46,111,110,108,105,110,101,119,101,98,102,111,110,116,115,46,99,111,109,47,105,99,111,110,32,45,45,62,13,10,60,33,68,79,67,84,89,80,69,32,115,118,103,32,80,85,66,76,73,67,32,34,45,47,47,87,51,67,47,47,68,84,68,32,83,86,71,32,49,46,49,47,47,69,78,34,32,34,104,116,116,112,58,47,47,119,119,119,46,119,51,46,111,114,103,47,71,114,97,112,104,105,99,115,47,83,86,71,47,49,46,49,47,68,84,68,47,115,118,103,49,49,46,100,116,100,34,62,13,10,60,115,118,103,32,118,101,114,115,105,111,110,61,34,49,46,49,34,32,120,109,108,110,115,61,34,104,116,116,112,58,47,47,119,119,119,46,119,51,46,111,114,103,47,50,48,48,48,47,115,118,103,34,32,120,109,108,110,115,58,120,108,105,110,107,61,34,104,116,116,112,58,47,47,119,119,119,46,119,51,46,111,114,103,47,49,57,57,57,47,120,108,105,110,107,34,32,120,61,34,48,112,120,34,32,121,61,34,48,112,120,34,32,118,105,101,119,66,111,120,61,34,48,32,48,32,49,48,48,48,32,49,48,48,48,34,32,101,110,97,98,108,101,45,98,97,99,107,103,114,111,117,110,100,61,34,110,101,119,32,48,32,48,32,49,48,48,48,32,49,48,48,48,34,32,120,109,108,58,115,112,97,99,101,61,34,112,114,101,115,101,114,118,101,34,62,13,10,60,109,101,116,97,100,97,116,97,62,32,83,118,103,32,86,101,99,116,111,114,32,73,99,111,110,115,32,58,32,104,116,116,112,58,47,47,119,119,119,46,111,110,108,105,110,101,119,101,98,102,111,110,116,115,46,99,111,109,47,105,99,111,110,32,60,47,109,101,116,97,100,97,116,97,62,13,10,60,103,62,60,112,97,116,104,32,100,61,34,77,57,55,57,46,57,44,54,51,49,46,54,99,48,45,51,51,46,57,45,49,51,46,55,45,53,57,46,56,45,51,54,46,49,45,55,52,46,56,99,49,51,46,54,45,55,46,56,44,51,53,46,57,45,52,54,46,54,44,51,53,46,57,45,56,48,46,57,99,45,49,46,51,45,53,54,46,54,45,53,52,46,57,45,49,50,54,46,55,45,49,51,53,46,57,45,49,50,54,46,55,108,45,50,49,52,46,53,44,48,67,54,54,54,44,50,52,55,46,49,44,54,53,57,44,49,51,55,46,54,44,54,49,50,46,57,44,54,52,46,49,67,53,56,53,46,50,44,49,57,46,57,44,53,53,49,46,56,44,49,48,44,53,50,56,46,55,44,49,48,99,45,55,55,46,56,44,48,45,57,51,46,51,44,55,54,46,54,45,57,51,46,51,44,49,48,50,46,54,99,48,44,56,55,46,49,45,49,50,46,51,44,49,50,54,46,57,45,50,55,46,55,44,49,54,51,46,54,99,45,49,50,46,49,44,50,51,46,55,45,55,51,46,56,44,49,49,49,46,49,45,49,52,57,46,57,44,49,49,49,46,49,72,49,49,53,46,49,99,45,53,53,46,49,44,48,45,57,53,46,49,44,52,55,46,49,45,57,53,44,49,49,52,46,49,108,51,55,44,51,56,57,46,52,99,52,46,53,44,54,49,46,52,44,51,52,44,57,56,46,54,44,57,51,46,51,44,57,56,46,54,99,48,44,48,44,50,57,46,53,44,48,46,51,44,53,49,46,54,44,48,46,51,99,49,49,46,51,44,48,44,50,49,46,56,45,48,46,49,44,50,56,45,48,46,51,99,49,56,46,51,45,48,46,54,44,51,51,46,51,45,49,50,46,54,44,52,54,46,54,45,50,50,46,54,99,57,46,52,45,55,46,49,44,49,54,46,57,45,49,51,46,50,44,50,51,46,55,45,49,51,46,50,99,51,44,48,44,50,48,46,51,44,49,48,46,49,44,52,49,46,54,44,49,57,46,51,99,50,49,46,54,44,57,46,52,44,53,54,44,49,55,46,49,44,56,56,46,57,44,49,55,46,49,104,50,57,56,46,53,99,49,48,48,44,48,44,49,52,56,46,49,45,51,56,46,51,44,49,52,56,46,49,45,49,48,50,46,49,99,48,45,49,56,46,56,45,52,46,50,45,51,48,45,49,52,46,54,45,52,52,46,55,99,52,48,46,50,45,57,46,57,44,55,53,46,56,45,51,53,46,54,44,55,53,46,56,45,56,48,46,49,99,48,45,49,54,46,53,45,55,46,54,45,52,51,46,49,45,49,56,45,53,48,67,57,52,48,46,51,44,55,48,56,44,57,55,57,46,57,44,54,55,50,46,50,44,57,55,57,46,57,44,54,51,49,46,54,122,32,77,49,56,49,46,53,44,57,49,54,46,51,99,45,50,51,44,48,45,52,51,46,52,45,50,50,46,49,45,52,56,46,54,45,53,49,46,57,76,57,52,46,55,44,53,50,53,46,49,99,48,45,51,49,44,50,52,46,52,45,54,53,44,53,55,46,49,45,54,53,118,49,46,52,108,57,54,44,48,46,51,118,52,53,52,99,45,52,46,55,44,48,46,52,45,57,46,52,45,50,46,49,45,49,50,46,49,45,49,46,56,67,50,49,55,46,56,44,57,49,52,46,54,44,49,56,49,46,57,44,57,49,54,46,51,44,49,56,49,46,53,44,57,49,54,46,51,122,32,77,56,50,56,46,57,44,54,56,53,46,57,99,48,44,48,45,53,46,52,44,48,45,53,46,52,44,48,108,48,46,53,44,51,55,99,52,46,51,44,49,46,55,44,50,48,46,52,44,54,46,55,44,50,48,46,52,44,51,56,46,56,99,48,44,51,52,46,52,45,51,52,46,57,44,52,56,46,51,45,55,50,46,51,44,52,56,46,51,104,48,46,49,108,45,48,46,55,44,51,53,46,51,99,49,53,46,50,44,51,46,53,44,50,48,46,55,44,49,55,44,50,48,46,54,44,51,56,46,55,99,48,44,51,55,46,57,45,51,57,46,54,44,51,50,46,51,45,54,51,46,57,44,50,57,46,57,72,52,49,57,46,51,99,45,49,52,46,56,44,48,45,55,56,46,57,45,51,50,46,51,45,57,54,46,50,45,52,48,46,56,99,45,48,46,51,45,48,46,50,45,48,46,53,45,48,46,49,45,48,46,56,45,48,46,50,86,52,53,48,46,51,99,55,53,46,57,45,51,48,46,54,44,49,51,57,46,57,45,49,49,51,46,57,44,49,53,53,46,53,45,49,52,52,46,54,99,49,55,46,54,45,52,49,46,56,44,49,50,46,52,45,55,48,46,50,44,49,50,46,52,45,49,54,54,46,49,99,48,45,52,55,46,54,44,49,55,46,49,45,53,52,44,51,53,46,52,45,53,52,99,49,54,46,49,44,48,44,51,50,46,51,44,52,44,52,53,46,51,44,50,52,46,55,99,51,56,46,52,44,54,49,46,49,44,50,49,46,53,44,49,55,53,46,51,45,50,54,46,50,44,50,55,52,108,45,56,46,51,44,52,49,46,51,108,51,48,57,46,50,45,48,46,49,99,52,54,46,54,44,50,46,54,44,52,52,46,50,44,49,53,44,52,52,46,50,44,53,48,46,50,108,45,48,46,51,44,49,49,99,48,46,49,44,48,46,53,45,49,46,50,44,51,48,46,53,45,49,54,46,56,44,52,51,46,55,99,45,49,51,46,49,44,49,49,46,49,45,50,51,46,55,44,49,49,46,53,45,50,51,46,55,44,49,49,46,53,108,45,48,46,51,44,51,54,46,52,99,48,44,48,44,49,53,46,56,44,50,46,50,44,50,52,46,49,44,55,46,57,99,49,50,46,52,44,56,46,55,44,49,56,44,50,53,46,50,44,49,55,46,52,44,52,51,46,50,67,56,56,57,46,51,44,54,53,50,46,49,44,56,53,48,46,49,44,54,56,53,46,57,44,56,50,56,46,57,44,54,56,53,46,57,122,34,47,62,60,47,103,62,13,10,60,47,115,118,103,62}}

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


    news := widget.NewGroup("News")
    news1 := &widget.Label{Text: "news 1", Alignment: fyne.TextAlignCenter}
    news.Append(news1)
    news2 := &widget.Label{Text: "news 2", Alignment: fyne.TextAlignCenter}
    news.Append(news2)
    glob.Append(news)

    paps := widget.NewGroupWithScroller("Papers")
    for i:=0;i<npaps;i++ {
        j := i
        lab := &widget.Label{Text: "paper "+strconv.Itoa(i), Alignment: fyne.TextAlignCenter}
        blike := widget.NewButton("", func() {paperlike(j)})
        blike.SetIcon(likesvg)
        bchat := widget.NewButton("", func() {paperchat(j)})
        bchat.SetIcon(theme.VisibilityIcon())
        papers[i] = widget.NewHBox(blike,layout.NewSpacer(),lab,layout.NewSpacer(),bchat)
        paps.Append(papers[i])
    }

    /*
    t := widget.NewEntry()
    t.OnCursorChanged = func() {
        fmt.Println("ttyypass")
    }
    */

    w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(glob,nil,nil,nil),glob,paps))
    w.ShowAndRun()
}

