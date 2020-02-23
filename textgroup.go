package main

import (
    "fmt"
    "time"
    "strings"
    "image/color"
    "fyne.io/fyne"
    "fyne.io/fyne/canvas"
    "fyne.io/fyne/layout"
    "fyne.io/fyne/theme"
    "fyne.io/fyne/widget"
)

const padding = 2

/*
   This file proposes a new "TextGroup" widget that allows to stack several texts into a stream of messages,
   where each message is automatically wrapped according to the available width.
   The messages are presented one after the other, until the available height is filled.
   Then, a button to show the next (or previous) page of messages is shown.
   I have chosen this "static" display of a list of messages instead of the more common scrolling stream,
   because the first tests I have done with a scrolling list implement with fyne.io were not smooth enough.

   Implementation details:
   - First, a WidthProgressBar is placed on top of the group and waits for 1 second to be sure that the display is stable;
     this widget then computes the available width, and computes the wrapping of all texts statically.
   - The widget then generates the labels with the computed wrapping and manages the prev/next buttons

   Limitations: this widget does not support changing orientation, or changing the window size dynamically for now.
*/

type Bord struct {
   canvas.Line
   w int
}

func NewBord(wh int) *Bord {
    b := &Bord{}
    b.Line = *canvas.NewLine(color.Black)
    b.w=wh
    b.Position1 = fyne.NewPos(0,0)
    b.Position2 = fyne.NewPos(wh,0)
    return b
}

func (b *Bord)MinSize() fyne.Size {
    return fyne.NewSize(b.w,3)
}

type TextGroup struct {
    widget.Group
    t0 int64
    winw, winh, curpage, lineh int
    txts []string
    word2width map[string]int
    Appwin *fyne.Container
    LabObj []*widget.Label
}

func NewTextGroup(tit string, txt []string) *TextGroup {
    txtg := &TextGroup{t0:0, winw:100, winh:180, curpage:0, lineh:0}
    g := widget.NewGroup(tit)
    txtg.Group = *g
    txtg.ExtendBaseWidget(txtg)
    txtg.txts = txt
    txtg.word2width = make(map[string]int)
    return txtg
}

func calcLabels(t *TextGroup) {
    t.t0 = -t.t0
    fmt.Printf("win fixed %d %d\n",t.winw,t.winh)
    calcTxtSize(t)
    createLabels(t)
    // on android, we need to refresh the global app window after we change the layout
    if t.Appwin != nil {t.Appwin.Refresh()}
}

func (t *TextGroup) Resize(winsize fyne.Size) {
    if t.t0<0 {return}
    t.winw = winsize.Width
    t.winh = winsize.Height
    if t.t0==0 {
        t.t0 = time.Now().Unix() // in sec
        go func() {
            time.Sleep(time.Second)
            // now we assume the display is stable
            calcLabels(t)
        }()
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
            }
        }
    }
    t.lineh = fullsize.Height
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

func createLabels(t *TextGroup) {
    wmax := int(0.95*float32(t.winw))
    hmax := int(0.95*float32(t.winh))
    posendlab := 0
    t.LabObj = make([]*widget.Label,0)
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
        posendlab += t.lineh * nlines
        posendlab += t.lineh // interline
        fmt.Printf("posendlab %v %v\n",posendlab,t.winh)
        if posendlab>=hmax {
            t.curpage = i
            fmt.Println("cut at: "+t.txts[i])
            break
        }
        lab := widget.NewLabel(sfin)
        t.LabObj = append(t.LabObj,lab)
        l := NewBord(t.winh)
        // l := canvas.NewLine(color.Black)
        // l.Position1 = fyne.NewPos(0,0)
        // l.Position2 = fyne.NewPos(t.winh,10)
        clab := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil,nil,nil,nil),l)
        t.Append(clab)
    }

    go func() {
        time.Sleep(time.Second)
        // draw lines to separate labels
        /*
        var lab *widget.Label
        var lines [5]*canvas.Line
        for i:=0;i<len(t.LabObj);i++ {
            lines[i] = canvas.NewLine(color.Black)
            // quand on addobject, il recalcule toutes les positions precedentes !
            t.Appwin.AddObject(lines[i])
        }
        for i:=0;i<len(t.LabObj);i++ {
            lab = t.LabObj[i]
            j := lab.Position().Y
            lines[i].Position1 = fyne.NewPos(0,j)
            lines[i].Position2 = fyne.NewPos(t.winh,j)
            fmt.Printf("zzz %v %v \n",lines[i].Position1,lines[i].Position2)
        }
        for i:=0;i<len(t.LabObj);i++ {
            fmt.Printf("yyy %v %v \n",lines[i].Position1,lines[i].Position2)
        }
        */
    }()
    
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
