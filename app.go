package main

import (
    "fmt"
    "time"
    "strconv"
    "fyne.io/fyne"
    "fyne.io/fyne/layout"
    "fyne.io/fyne/app"
    "fyne.io/fyne/widget"
)

const npaps = 10
var papers [npaps]*widget.Button

func paperhdl(i int) {
  fmt.Printf("paper %d\n",i)
  papers[i].Disable()
  go func() {
      time.Sleep(time.Second)
      papers[i].Enable()
  }()
}

func main() {
    app := app.New()

    w := app.NewWindow("GOLKi")
    glob := widget.NewVBox()
    titleday := &widget.Label{Text: "Monday", Alignment: fyne.TextAlignCenter}
    glob.Append(titleday)

    news := widget.NewGroup("News")
    news1 := &widget.Label{Text: "news 1", Alignment: fyne.TextAlignCenter}
    news.Append(news1)
    news2 := &widget.Label{Text: "news 2", Alignment: fyne.TextAlignCenter}
    news.Append(news2)
    glob.Append(news)

    paps := widget.NewGroupWithScroller("Papers")
    for i:=0;i<npaps;i++ {
        //papers[i] = &widget.Label{Text: "paper "+strconv.Itoa(i), Alignment: fyne.TextAlignCenter}
        j := i
        papers[i] = widget.NewButton("paper "+strconv.Itoa(i), func() {paperhdl(j)})
        paps.Append(papers[i])
    }

    /*
    t := widget.NewEntry()
    t.OnCursorChanged = func() {
        fmt.Println("ttyypass")
    }
    */

    glob.Refresh()
    w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(glob,nil,nil,nil),glob,paps))
    w.ShowAndRun()
}

