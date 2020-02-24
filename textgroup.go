package main

import (
    "fmt"
    "strings"
    "image/color"
    "fyne.io/fyne"
    "fyne.io/fyne/layout"
    "fyne.io/fyne/canvas"
    "fyne.io/fyne/theme"
    "fyne.io/fyne/widget"
)

const padding = 2
const interline = 10
const spacew = 2
var buth int = 0

type TextGroup struct {
    widget.BaseWidget
    txts []string
    word2width map[string]int
    lineh int
    curpage int
    pageidx []int
}

type TapLab struct {
    widget.Label
    i int
}

func (b *TapLab) Tapped(*fyne.PointEvent) {
    fmt.Printf("tapped %d\n",b.i)
    // TODO
}
func (b *TapLab) TappedSecondary(*fyne.PointEvent) {
}
func NewTapLab(s string,j int) *TapLab {
    l := &TapLab{widget.Label{Text:s},j}
    l.ExtendBaseWidget(l)
    return l
}

func settings(tg *TextGroup) {
    s := aptest()
    tg.reset(s)
    tg.Refresh()
}

func (b *TextGroup) CreateRenderer() fyne.WidgetRenderer {
        var objects []fyne.CanvasObject
        labelsBox := widget.NewVBox()

        bprev := widget.NewButtonWithIcon("",theme.MoveUpIcon(),func () {
            b.curpage--
            if b.curpage<0 {b.curpage=0}
            b.Refresh()
        })
        bnext := widget.NewButtonWithIcon("",theme.MoveDownIcon(),func () {
            b.curpage++
            if b.curpage>=len(b.pageidx) {b.curpage--}
            b.Refresh()
        })
        bsettings := widget.NewButtonWithIcon("",theme.MenuIcon(),func () {
            settings(b)
        })
        pnbuttons := widget.NewHBox(bprev,layout.NewSpacer(),bsettings,layout.NewSpacer(),bnext)
        buth = pnbuttons.MinSize().Height

        whole := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil,pnbuttons,nil,nil),labelsBox,pnbuttons)
        objects = append(objects,whole)
        r := &textGroupRenderer{objects, labelsBox, whole, b}
        return r
}

func (tg *TextGroup) reset(s []string) {
        tg.txts = s
        tg.ExtendBaseWidget(tg)
        tg.word2width = make(map[string]int)
        tg.curpage = 0
        tg.pageidx =[]int{0}
}

func NewTextGroup(txt []string) *TextGroup{
        tg := &TextGroup{txts: txt}
        tg.ExtendBaseWidget(tg)
        tg.word2width = make(map[string]int)
        tg.curpage = 0
        tg.pageidx =[]int{0}
        return tg
}

// ========

type textGroupRenderer struct {
        objects []fyne.CanvasObject
        labelsBox *widget.Box
        whole *fyne.Container
        tg  *TextGroup
}

func (b *textGroupRenderer) MinSize() fyne.Size {
        baseSize := fyne.NewSize(150,200)
        // b.label.MinSize()
        baseSize = baseSize.Add(fyne.NewSize(24, 24))
        return baseSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

func (b *textGroupRenderer) Layout(size fyne.Size) {
        b.whole.Resize(size)
        b.recalcText(size)
}

func (b *textGroupRenderer) ApplyTheme() {
        b.Refresh()
}

func (b *textGroupRenderer) BackgroundColor() color.Color {
        return theme.ButtonColor()
}

func (b *textGroupRenderer) Refresh() {
        b.Layout(b.tg.Size())
        canvas.Refresh(b.tg)
}

func (b *textGroupRenderer) Objects() []fyne.CanvasObject {
        return b.objects
}

func (b *textGroupRenderer) Destroy() {
}

func (b *textGroupRenderer) recalcText(size fyne.Size) {
    calcTxtSize(b.tg)
    b.createLabels(b.tg,size)
}

// ==================

func (b *textGroupRenderer) createLabels(t *TextGroup, size fyne.Size) {
    wmax := int(0.99*float32(size.Width))
    hmax := int(0.99*float32(size.Height))-buth
    b.labelsBox.Children = b.labelsBox.Children[:0]
    posendlab := 0
    for i:=t.pageidx[t.curpage];i<len(t.txts);i++ {
        sfin := ""
        ss := strings.Split(t.txts[i]," ")
        cum:=padding
        nlines := 1
        for j:=0;j<len(ss);j++ {
            s := strings.TrimSuffix(ss[j],"\n")
            sl := t.word2width[s]
            cum += sl+spacew
            if cum>wmax {
                sfin += "\n"
                nlines++
                cum=padding+sl+3
            }
            sfin += s+" "
        }

        // estimate the height of this piece of text
        posendlab += t.lineh * nlines
        posendlab += interline
        if posendlab>=hmax && i>t.pageidx[t.curpage] {
            t.pageidx = t.pageidx[:t.curpage+1]
            t.pageidx = append(t.pageidx,i)
            break
        }
        newlab := NewTapLab(sfin,i)
        b.labelsBox.Append(newlab)
    }
    // b.labelsBox.Refresh()
}

func calcTxtSize(t *TextGroup) {
    fontsize := theme.TextSize()
    var fullsize fyne.Size
    for i:=0;i<len(t.txts);i++ {
        ss := strings.Split(t.txts[i]," ")
        for j:=0;j<len(ss);j++ {
            s := strings.TrimSuffix(ss[j],"\n")
            _, ok := t.word2width[s]
            if !ok {
                tt := canvas.NewText(s,color.Black)
                tt.TextSize = fontsize
                tt.TextStyle = tt.TextStyle
                fullsize = tt.MinSize()
                t.word2width[s]=fullsize.Width
                if fullsize.Height > t.lineh { t.lineh = fullsize.Height }
            }
        }
    }
    if false {
        for x,n := range t.word2width {
            fmt.Printf("size %v %v\n",x,n)
        }
        // utilise cela pour estimer la taille d'un espace = 3 pixels
        s := "même pas"
        tt := canvas.NewText(s,color.Black)
        tt.TextSize = fontsize
        tt.TextStyle = tt.TextStyle
        fullsize = tt.MinSize()
        fmt.Printf("même pas %v\n",fullsize.Width)
    }
}

