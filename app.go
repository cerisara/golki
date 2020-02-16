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
var blikes [npaps]*widget.Button
var isliked [npaps]bool
var days = [...]string {"Lundi","Mardi","Mercredi","Jeudi","Vendredi"}
var curday = 0


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

func setMultiLines(s string, l *widget.Label) {
    if l!=nil {
        ss := s[0:10]+"\n"+s[10:]
        l.SetText(ss)
    }
}

type winMesure struct {
    widget.BaseWidget
    winsize int
    lab2update *widget.Label
}
func (t *winMesure) Resize(wins fyne.Size) {
    t.winsize = wins.Width
    setMultiLines("tétéjjjjjjôïkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkü", t.lab2update)
}
func NewWinMesure(labelWidget *widget.Label) *winMesure {
    w := &winMesure{lab2update: labelWidget}
    w.ExtendBaseWidget(w)
    return w
}

func main() {
    app := app.New()
    w := app.NewWindow("GOLKi")
    glob := widget.NewVBox()

    testlab := widget.NewLabel("")
    glob.Append(NewWinMesure(testlab))
    glob.Append(testlab)

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

    // scroller is not fluid enough
    // paps := widget.NewGroupWithScroller("Papers")
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

    /*
    t := widget.NewEntry()
    t.OnCursorChanged = func() {
        fmt.Println("ttyypass")
    }
    */

    w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(glob,nil,nil,nil),glob,paps))
    fmt.Printf("uinsize %v %v\n",w.Canvas().Size(),paps.Size())
    w.ShowAndRun()
}

