package main

import (
    "fmt"
    "strings"
    "image/color"
    "fyne.io/fyne"
    "fyne.io/fyne/canvas"
    "fyne.io/fyne/theme"
    "fyne.io/fyne/widget"
)

const padding = 2

type TextGroup struct {
    widget.BaseWidget
    txts []string
    word2width map[string]int
    lineh int
    curpage int
}

func (b *TextGroup) Tapped(*fyne.PointEvent) {
    // TODO
}
func (b *TextGroup) TappedSecondary(*fyne.PointEvent) {
}

func (b *TextGroup) CreateRenderer() fyne.WidgetRenderer {
        var objects []fyne.CanvasObject
        return &textGroupRenderer{objects, b}
}

func NewTextGroup(txt []string) *TextGroup{
        tg := &TextGroup{txts: txt}
        tg.ExtendBaseWidget(tg)
        tg.word2width = make(map[string]int)
        tg.curpage = 0
        return tg
}

// ========

type textGroupRenderer struct {
        objects []fyne.CanvasObject
        tg  *TextGroup
}

func (b *textGroupRenderer) MinSize() fyne.Size {
        baseSize := fyne.NewSize(100,150)
        // b.label.MinSize()
        baseSize = baseSize.Add(fyne.NewSize(24, 24))
        return baseSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

func (b *textGroupRenderer) Layout(size fyne.Size) {
        inner := size.Subtract(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
        inner = inner.Subtract(fyne.NewSize(24, 24))
        b.recalcText(inner)
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
    // on android, we need to refresh the global app window after we change the layout
    // if t.Appwin != nil {t.Appwin.Refresh()}
}

// ==================

func (b *textGroupRenderer) createLabels(t *TextGroup, size fyne.Size) {
    wmax := int(0.95*float32(size.Width))
    hmax := int(0.95*float32(size.Height))
    b.objects = b.objects[:0]
    posendlab := 0
    for i:=t.curpage;i<len(t.txts);i++ {
        var w2w = make([]int,10)
        sfin := ""
        ss := strings.Split(t.txts[i]," ")
        cum:=padding
        nlines := 1
        for j:=0;j<len(ss);j++ {
            s := strings.TrimSuffix(ss[j],"\n")
            sl := t.word2width[s]
            w2w = append(w2w,sl)
            cum += sl
            cum += 3
            if cum>=wmax {
                sfin += "\n"
                nlines++
                cum=padding+sl+3
            }
            sfin += s+" "
        }

        // estimate the height of this piece of text
        prevpos := posendlab
        posendlab += t.lineh * nlines
        posendlab += t.lineh // interline
        if posendlab>=hmax {
            // t.curpage = i
            break
        }
        newlab := widget.NewLabel(sfin)
        b.objects = append(b.objects,newlab)
        newlab.Move(fyne.NewPos(0,prevpos))
        // fmt.Printf("move label %d %d\n",i,prevpos)
    }
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








/*

Old stuff
I keep them just for the records

// label normal mais qui ne fixe pas de mindwidth
// warning: ce reactLabel ne se mets pas a jour lorsqu'on change son texte avec SetText !!?
type reactLabel struct {
    widget.Label
    text0 string
    OnTapped func()
}

func (b *reactLabel) Tapped(*fyne.PointEvent) {
    b.OnTapped()
}
func (b *reactLabel) TappedSecondary(*fyne.PointEvent) {
}

func NewReactLabel(txt string) *reactLabel {
    label := &reactLabel{OnTapped: func() {fmt.Println("TAPPED")}}
    label.ExtendBaseWidget(label)
    label.text0 = txt
    label.SetText("") // on le mettra a jour lors du premier call a Resize
    return label
}
// juste pour qu'on puisse reduire la taille lors du test et pour qu'il ne cree par une fenetre en se basant sur la taille de ce texte
func (l *reactLabel) MinSize() fyne.Size {
    return fyne.NewSize(50,50)
}

/////////////////////////////////////////////////////////
type winMeasurer struct {
    widget.Label
    winw int
    Labels2wrap []*widget.Label
    txt0 []string
}
func (_ *winMeasurer) MinSize() fyne.Size {
    return fyne.NewSize(0,1)
}
func NewWinMeasurer() *winMeasurer {
    w := &winMeasurer{}
    w.ExtendBaseWidget(w)
    w.SetText("")
    return w
}
func (w *winMeasurer) AddLabel2wrap(l *widget.Label, s string) {
    w.Labels2wrap = append(w.Labels2wrap,l)
    w.txt0 = append(w.txt0,s)
}
func (t *winMeasurer) Resize(winsize fyne.Size) {
    t.winw = winsize.Width
    for i:=0;i<len(t.Labels2wrap);i++ {
        ll := t.Labels2wrap[i]
        s := t.txt0[i]
        var ss string
        if t.winw < 600 {
            ss = s[:10]+"\n"+s[10:]
        } else {
            ss = s
        }
        // fmt.Printf("winsize %v %s\n",t.winw,s)
        ll.SetText(ss)
    }

    // estimate the size of the full text
    tt := canvas.NewText(t.alltxt, color.Black)
    fontsize := theme.TextSize()
    tt.TextSize = fontsize
    tt.TextStyle = t.TextStyle
    fullsize := tt.MinSize()
    if fullsize.Width>=winsize.Width {
        // need to wrap
        s := t.Text
        szPerLetter := float64(fullsize.Width)/float64(len(s))
        nletPerLine := int32(float64(winsize.Width)*0.8 / float64(szPerLetter))
        ss := s[:nletPerLine]+"\n"+s[nletPerLine:]
        fmt.Printf("cutline %v %v %s\n",nletPerLine,len(s),ss)
        t.SetText(ss)
        t.Refresh()
    }
    // fullsize := textMinSize(t.alltxt,,t.textStyle())
}

*/
