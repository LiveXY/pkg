package seal

import (
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/livexy/pkg/strx"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

func FontFileName(fontdata draw2d.FontData) string {
	fontname := fontdata.Name
	if fontdata.Style&draw2d.FontStyleBold != 0 {
		fontname += "b"
	}
	if fontdata.Style&draw2d.FontStyleItalic != 0 {
		fontname += "i"
	}
	fontname += ".ttf"
	return fontname
}

func Init(fontpath string) {
	draw2d.SetFontFolder(fontpath)
	draw2d.SetFontNamer(FontFileName)
}

type CompanySeal struct {
	hexcolor             color.Color
	gc                   *draw2dimg.GraphicContext
	dest                 *image.RGBA
	companyName          string
	fontName             string
	securityCode         string
	typeName             string
	fontdata             draw2d.FontData
	companyNameFontSize  float64
	typeNameFontSize     float64
	radius               float64
	securityCodeFontSize float64
	innerLineWidth       float64
	lineWidth            float64
	isSecurityCode       bool
	hasInnerLine         bool
}

func NewCompanySeal(w, h int, companyName, fontName string) *CompanySeal {
	s := &CompanySeal{
		companyName: companyName, fontName: fontName,
		dest: image.NewRGBA(image.Rect(0, 0, w, h)),
	}
	s.hexcolor = color.RGBA{0xe5, 0x0, 0x20, 0xff}
	s.radius = float64(w) / float64(2)
	s.lineWidth = float64(4) / float64(75) * s.radius
	s.hasInnerLine = true
	s.innerLineWidth = float64(1) / float64(75) * s.radius
	s.gc = draw2dimg.NewGraphicContext(s.dest)
	s.companyNameFontSize = float64(14) / float64(75) * s.radius
	s.typeNameFontSize = float64(11) / float64(75) * s.radius
	s.securityCodeFontSize = 0.12 * s.radius
	s.fontdata = draw2d.FontData{Name: fontName}
	return s
}
func (s *CompanySeal) IsSecurityCode(isSecurityCode bool) *CompanySeal {
	s.isSecurityCode = isSecurityCode
	return s
}
func (s *CompanySeal) HideInnerLine() *CompanySeal {
	s.hasInnerLine = false
	return s
}
func (s *CompanySeal) SetTypeName(typeName string) *CompanySeal {
	s.typeName = typeName
	return s
}
func (s *CompanySeal) SetColor(c string) *CompanySeal {
	s.hexcolor = strx.ToColor(c, color.RGBA{0xe5, 0x0, 0x20, 0xff})
	return s
}
func (s *CompanySeal) SetSecurityCode(securityCode string) *CompanySeal {
	s.securityCode = securityCode
	return s
}
func (s *CompanySeal) drawOuterLine() {
	s.gc.SetStrokeColor(s.hexcolor)
	s.gc.SetLineWidth(s.lineWidth)
	s.gc.BeginPath()
	s.gc.ArcTo(s.radius, s.radius, s.radius-s.lineWidth, s.radius-s.lineWidth, 0, math.Pi*2)
	s.gc.Stroke()
	s.gc.Restore()
}
func (s *CompanySeal) drawInnerLine() {
	s.gc.SetStrokeColor(s.hexcolor)
	s.gc.SetLineWidth(s.innerLineWidth)
	s.gc.BeginPath()
	s.gc.ArcTo(s.radius, s.radius, s.radius-s.lineWidth-s.radius/float64(15), s.radius-s.lineWidth-s.radius/float64(15), 0, math.Pi*2)
	s.gc.Stroke()
	s.gc.Restore()
}
func (s *CompanySeal) drawStar() {
	var R = s.radius
	var r = R / float64(3)
	var c = 360 / float64(5) * math.Pi / 180
	var d = c / float64(2)
	var e = d / float64(2)
	var l = r * math.Sin(e) / math.Sin(d+e)
	var lsd = l * math.Sin(d)
	var lcd = l * math.Cos(d)
	var lsc = l * math.Sin(c)
	var lcc = l * math.Cos(c)
	var rsc = r * math.Sin(c)
	var rcc = r * math.Cos(c)
	var rsd = r * math.Sin(d)
	var rcd = r * math.Cos(d)
	var p0 = []float64{R, float64(2) / float64(3) * R}
	var p1 = []float64{R + lsd, R - lcd}
	var p2 = []float64{R + rsc, R - rcc}
	var p3 = []float64{R + lsc, R + lcc}
	var p4 = []float64{R + rsd, R + rcd}
	var p5 = []float64{R, R + l}
	var p6 = []float64{R - rsd, R + rcd}
	var p7 = []float64{R - lsc, R + lcc}
	var p8 = []float64{R - rsc, R - rcc}
	var p9 = []float64{R - lsd, R - lcd}
	var aPs = [][]float64{p0, p1, p2, p3, p4, p5, p6, p7, p8, p9}
	s.gc.Save()
	s.gc.SetFillColor(s.hexcolor)
	s.gc.BeginPath()
	for _, v := range aPs {
		s.gc.LineTo(v[0], v[1])
	}
	s.gc.Close()
	s.gc.Fill()
	s.gc.Restore()
}
func (s *CompanySeal) drawLetter(letter string, angle, x, y float64) {
	s.gc.Save()
	s.gc.Rotate(angle)
	s.gc.FillStringAt(letter, x, y)
	s.gc.Restore()
}
func (s *CompanySeal) drawCompanyName() {
	s.drawText(s.companyName, s.companyNameFontSize, false, false)
}
func (s *CompanySeal) drawTypeName() {
	s.drawText(s.typeName, s.typeNameFontSize, true, false)
}
func (s *CompanySeal) drawSecurityCode() {
	s.drawText(s.securityCode, s.securityCodeFontSize, false, true)
}
func (s *CompanySeal) drawText(text string, fontSize float64, isTypeName, isSecurityCode bool) {
	texts := strings.Split(text, "")
	if isTypeName && len(texts) > 5 {
		text = strings.Join(texts[:5], "")
	} else {
		if isSecurityCode && len(texts) > 13 {
			texts = texts[:13]
		} else if len(texts) > 19 {
			texts = texts[:19]
		}
	}
	step := 0.32
	astartpos := []float64{0, -0.15, -0.31, -0.5, -0.63, -0.83, -0.95, -1.12, -1.27, -1.44, -1.61, -1.77, -1.93, -2.1, -2.23, -2.4, -2.57, -2.73, -2.89}
	startindex := len(texts) - 1
	startpos := astartpos[startindex]
	linetextcap := 0.65
	s.gc.Save()
	s.gc.SetFillColor(s.hexcolor)
	s.gc.SetFontData(s.fontdata)
	s.gc.SetFontSize(fontSize)
	s.gc.Translate(s.radius, s.radius)
	if isTypeName {
		left, _, right, _ := s.gc.GetStringBounds(text)
		s.gc.FillStringAt(text, -(right-left)/float64(2), s.radius/float64(2))
	} else {
		for i, letter := range texts {
			if isSecurityCode {
				s.drawLetter(letter, 0.57-float64(i)*0.1, -s.securityCodeFontSize/float64(2), s.radius*0.8)
			} else {
				s.drawLetter(letter, startpos+float64(i)*step, -s.companyNameFontSize/float64(2), -s.radius*linetextcap)
			}
		}
	}
	s.gc.Restore()
}

func (s *CompanySeal) Save(sealPath string) error {
	s.drawOuterLine()
	if s.hasInnerLine {
		s.drawInnerLine()
	}
	s.drawStar()
	if len(s.companyName) > 0 {
		s.drawCompanyName()
	}
	l := len(strings.Split(s.companyName, ""))
	if len(s.typeName) > 0 && l < 19 {
		s.drawTypeName()
	}
	if len(s.securityCode) > 0 && l <= 15 {
		s.drawSecurityCode()
	}
	return draw2dimg.SaveToPngFile(sealPath, s.dest)
}
