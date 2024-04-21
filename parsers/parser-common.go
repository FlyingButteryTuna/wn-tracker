package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type CommonParser struct{}

func (p *CommonParser) RemoveAttrsFromElement(elem string, doc *goquery.Selection) {
	doc.Find(elem).Each(func(_ int, s *goquery.Selection) {
		attrs := make([]html.Attribute, len(s.Nodes[0].Attr))
		copy(attrs, s.Nodes[0].Attr)

		for _, attr := range attrs {
			s.RemoveAttr(attr.Key)
		}
	})
}
func (p *CommonParser) RemoveRP(doc *goquery.Selection) {
	doc.Find("rp").Each(func(_ int, s *goquery.Selection) {
		s.Remove()
	})
}

func (p *CommonParser) ReplaceAttrInElement(elems *goquery.Selection, attrName string, attrValue string) {
	elems.Each(func(_ int, s *goquery.Selection) {
		s.SetAttr(attrName, attrValue)
	})
}
