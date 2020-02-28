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
    menufct MenuFct
    jumpEnd bool
}

type TapLab struct {
    widget.Label
    self widget.Label
    // je pourrais avoir besoin d'un textProvider pour changer la fontsize, mais je prefere laisser pour le moment
    // txtprov *textProvider
    i int
    tg *TextGroup
}

type tapLabRenderer struct {
        taplab  *TapLab
        labrenderer fyne.WidgetRenderer
        // objects []fyne.CanvasObject
}
func (b *TapLab) CreateRenderer() fyne.WidgetRenderer {
    // lr est un textRenderer: je le garde sous la main pour pouvoir réutiliser ses capacités
    lr := b.self.CreateRenderer()
    /*
    canvastexts := lr.texts
    for ct := range canvastexts {
        ct.TextSize = 8
    }
    */
    // mais je crée mon propre renderer pour changer la couleur et text fontsize
    r := &tapLabRenderer{b,lr}
    // text := canvas.NewText(b.Text, theme.TextColor())
    // text.TextStyle.Bold = true
    // objects := []fyne.CanvasObject{ text }
    return r
}
func (b *tapLabRenderer) Objects() []fyne.CanvasObject {
        return b.labrenderer.Objects()
}
func (b *tapLabRenderer) MinSize() fyne.Size {
        return b.labrenderer.MinSize()
}
func (b *tapLabRenderer) Layout(size fyne.Size) {
        b.labrenderer.Layout(size)
}
func (b *tapLabRenderer) BackgroundColor() color.Color {
        return theme.ButtonColor()
        // return color.RGBA{R: 0xff, G: 0x01, B: 0x01, A: 0xff}
}
func (b *tapLabRenderer) Refresh() {
        b.labrenderer.Refresh()
}
func (b *tapLabRenderer) Destroy() {
        b.labrenderer.Destroy()
}

func (b *TapLab) Tapped(*fyne.PointEvent) {
    fmt.Printf("tapped %d\n",b.i)
    b.tg.menufct(b.i)
}
func (b *TapLab) TappedSecondary(*fyne.PointEvent) {
}
func NewTapLab(s string,j int) *TapLab {
    l := &TapLab{}
    ll := widget.NewLabel(s)
    l.Label = *ll
    l.i = j
    l.self = *ll
    // l := &TapLab{widget.Label{Text:s},j}
    l.ExtendBaseWidget(l)
    return l
}

// --------------------------------------------------------------------------------

func (b *TextGroup) CreateRenderer() fyne.WidgetRenderer {
        var objects []fyne.CanvasObject
        labelsBox := widget.NewVBox()

        bprev := widget.NewButtonWithIcon("",theme.MoveUpIcon(),func () {
            b.curpage--
            if b.curpage<0 {
                b.menufct(-3)
                b.curpage=0
            }
            b.Refresh()
        })
        bnext := widget.NewButtonWithIcon("",theme.MoveDownIcon(),func () {
            b.curpage++
            if b.curpage>=len(b.pageidx) {
                b.curpage--
                b.menufct(-2)
            }
            b.Refresh()
        })
        bsettings := widget.NewButtonWithIcon("",theme.MenuIcon(),func () {
            b.menufct(-1)
        })
        pnbuttons := widget.NewHBox(bprev,layout.NewSpacer(),bsettings,layout.NewSpacer(),bnext)
        buth = pnbuttons.MinSize().Height

        whole := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil,pnbuttons,nil,nil),labelsBox,pnbuttons)
        objects = append(objects,whole)
        r := &textGroupRenderer{objects, labelsBox, whole, b}
        return r
}

func (tg *TextGroup) SetTexts(s []string) {
    tg.txts = s
    tg.ExtendBaseWidget(tg)
    tg.word2width = make(map[string]int)
    tg.curpage = 0
    tg.pageidx =[]int{0}
    tg.Refresh()
}

type MenuFct func(int)

func NewTextGroup(txt []string, mf MenuFct) *TextGroup{
    tg := &TextGroup{txts: txt, menufct: mf}
    tg.ExtendBaseWidget(tg)
    tg.word2width = make(map[string]int)
    tg.curpage = 0
    tg.pageidx =[]int{0}
    tg.jumpEnd = false
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
    for ;; {
        b.createLabels(b.tg,size)
        if !b.tg.jumpEnd {break}
        if b.tg.curpage>=len(b.tg.pageidx)-1 {break}
        b.tg.curpage++
    }
    b.tg.jumpEnd=false
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
        newlab.tg = b.tg
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

