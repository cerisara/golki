package main

import (
    "fmt"
    "image/color"
    "fyne.io/fyne"
    "fyne.io/fyne/widget"
    "fyne.io/fyne/canvas"
    "fyne.io/fyne/theme"
)

// label normal mais qui ne fixe pas de mindwidth

type reactLabel struct {
    widget.Label
    OnTapped func()
}

func (b *reactLabel) Tapped(*fyne.PointEvent) {
    b.OnTapped()
}
func (b *reactLabel) TappedSecondary(*fyne.PointEvent) {
}

func NewReactLabel(txt string) *reactLabel {
    label := &reactLabel{OnTapped: func() {}}
    label.ExtendBaseWidget(label)
    label.SetText(txt)
    return label
}

type reactLabelRenderer struct {
    label *canvas.Text

    objects []fyne.CanvasObject
    lab  *reactLabel
}

func (b *reactLabelRenderer) MinSize() fyne.Size {
    baseSize := b.label.MinSize()
    // w := baseSize.Width + theme.Padding()*2
    w := 50
    h := baseSize.Height + theme.Padding()*2
    h+=100
    nov := fyne.NewSize(w,h)
    return nov
}

func (b *reactLabelRenderer) Layout(size fyne.Size) {
    inner := size.Subtract(fyne.NewSize(0, theme.Padding()*2))
    b.label.Resize(inner)
}

func (b *reactLabelRenderer) ApplyTheme() {
    b.label.Color = theme.TextColor()
    b.label.TextSize = theme.TextSize()

    b.Refresh()
}

func (b *reactLabelRenderer) BackgroundColor() color.Color {
    return theme.TextColor()
}

func (b *reactLabelRenderer) Refresh() {
    fmt.Println("refresh")
    s := b.lab.alltxt[:10]+"\n "+b.lab.alltxt[10:]
    b.label.Text=s
    b.Layout(b.lab.Size())
    canvas.Refresh(b.lab)
}

func (b *reactLabelRenderer) Objects() []fyne.CanvasObject {
    return b.objects
}

func (b *reactLabelRenderer) Destroy() {
}

func (b *reactLabel) CreateRenderer() fyne.WidgetRenderer {
    var objects []fyne.CanvasObject
    text := canvas.NewText(b.Text, theme.TextColor())
    text.TextSize = theme.TextSize()
    text.Alignment = fyne.TextAlignLeading

    objects = append(objects, text)
    return &reactLabelRenderer{text, objects, b}
}

// this function is trigerred when the widget is resized, so we can calcule the new linebreaks
func (t *reactLabel) Resize(winsize fyne.Size) {
    t.winsize = winsize.Width

    /*
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
    */
    // fullsize := textMinSize(t.alltxt,,t.textStyle())
}

