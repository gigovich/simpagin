// Use of this source code is governed by BSD license
// 2014. Moscow. Givi Khojanashvili <gigovich@gmail.com>

// Simpagin is a simple pagination tool
package simpagin

import (
	"fmt"
	"log"
)

const (
	LEFT = iota
	MIDDLE
	RIGHT
)

type PageRenderer func(p Page) string

type Page struct {
	Index    int          // Object index (position) in total object list
	Number   int          // Page number valid value is from range [1..(itemsCount / frameLength)]
	IsActive bool         // Is this page active, paginator contains only one active page
	Type     int          // MIDDLE is all pages in frame, LEFT and RIGHT are scroller pages
	Renderer PageRenderer // Custom function to render page element as string called by Page.Render()
}

// String calls Page.Renderer to render string representation of page or some other page data
//
// If Renderer does not set, empty string will be returned. You can set
// Renderer function for all pages by Paginator.SetRenderer method.
func (p Page) String() string {
	if p.Renderer != nil {
		return p.Renderer(p)
	} else {
		log.Print("String renderer function is not set. See Paginator.SetRenderer ")
		return ""
	}
}

type Paginator struct {
	ActivePage  int     // Active page number
	LeftPage    *Page   // Page for left scroller, if active page is too close to start it must be nil
	RightPage   *Page   // Page for right scroller, if active page is too close to end it must be nil
	ItemsCount  int     // Total items count
	PagesCount  int     // Auto calculated field whish equals to ItemsCount / ItemsOnPage
	ItemsOnPage int     // How much items contains each page
	FrameLength int     // Number of pages displayed in paginator
	PageList    []*Page // You must fetch this slice to display each paginated page
}

// New returns new Paginator struct, with calculated fields which you can use,
// to render paginator.
//
// Exemple of usage:
// 	pg, err := New(
// 		10,  // Active page which items we displaying now
// 		120, // Total count of items
// 		8,   // We show only 8 items in each page
// 		10,  // And our paginator rendered as 10 pages list
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("<a href=\"/page/%s/\">&lt</a>", pg.LeftPage.Number)
// 	for _, page := range pg.PageList {
// 		if page.IsActive {
// 			fmt.Print(page.Number)
// 		} else {
// 			fmt.Printf("<a href=\"/page/%s/\">&lt</a>", pg.LeftPage.Number)
// 		}
// 	}
// 	fmt.Printf("<a href=\"/page/%s/\">&gt</a>", pg.RightPage.Number)
func New(activePage, itemsCount, itemsOnPage, frameLength int) (*Paginator, error) {
	if activePage*itemsOnPage > itemsCount {
		return nil, fmt.Errorf("Wrong page number for paginate")
	}
	if frameLength < 2 {
		return nil, fmt.Errorf("Paginated frame can't be less then 2")
	}
	pg := &Paginator{
		ActivePage:  activePage,
		ItemsCount:  itemsCount,
		ItemsOnPage: itemsOnPage,
		FrameLength: frameLength,
		LeftPage:    &Page{Type: LEFT},
		RightPage:   &Page{Type: RIGHT},
	}
	// Calculate PagesCount
	pg.PagesCount = itemsCount / itemsOnPage
	if itemsCount%itemsOnPage > 0 {
		pg.PagesCount++
	}
	// Calculate side indexes
	distanceToLeftSide := (frameLength / 2)
	distanceToRightSide := frameLength - (distanceToLeftSide + 1)
	frameStartIndex := 1
	if activePage > (distanceToLeftSide+1) && pg.PagesCount > frameLength {
		pg.LeftPage = &Page{(activePage - 1) * itemsOnPage, activePage - 1, false, LEFT, nil}
		frameStartIndex = activePage - distanceToLeftSide
	}
	if pg.PagesCount > frameLength && activePage < (pg.PagesCount-(distanceToRightSide+1)) {
		pg.RightPage = &Page{(activePage + 1) * itemsOnPage, activePage + 1, false, RIGHT, nil}
	}
	pages := make([]*Page, min(frameLength, pg.PagesCount))
	for i := 0; i < len(pages); i++ {
		pageNumber := i + frameStartIndex
		pages[i] = &Page{
			Index:  (pageNumber - 1) * itemsOnPage,
			Number: pageNumber,
			Type:   MIDDLE,
		}
		if pageNumber == activePage {
			pages[i].IsActive = true
		}
	}
	pg.PageList = pages
	return pg, nil
}

// SetRenderer set for all page object in the paginator PageRenderer function.
//
// Exemple of usage:
// 	pg, err := simpagin.New(
// 		10,  // Active page which items we displaying now
// 		120, // Total count of items
// 		8,   // We show only 8 items in each page
// 		10,  // And our paginator rendered as 10 pages list
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pg.SetRenderer(func (p Page) string {
// 		switch p.Type {
// 		case simpagin.LEFT:
// 			if p.Number == 0 {
// 				return "<li class=\"disabled\"><span>&laquo;</span></li>"
// 			} else {
// 				return fmt.Sprintf("<li><a href=\"?p=%d\">&laquo;</a></li>", p.Number)
// 			}
// 		case simpagin.MIDDLE:
// 			if p.IsActive {
// 				return fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", p.Number)
// 			} else {
// 				return fmt.Sprintf("<li><a href=\"?p=%d\">%d</a></li>", p.Number, p.Number)
// 			}
// 		case simpagin.RIGHT:
// 			if p.Number == 0 {
// 				return "<li class=\"disabled\"><span>&raquo;</span></li>"
// 			} else {
// 				return fmt.Sprintf("<li><a href=\"?p=%d\">&raquo;</a></li>", p.Number, p.Number)
// 			}
// 		}
// 		return ""
// 	})
// 	fmt.Print(pg.LeftPage)
// 	for _, page := range pg.PageList {
// 		fmt.Print(page)
// 	}
// 	fmt.Print(pg.RightPage)
func (p *Paginator) SetRenderer(f PageRenderer) {
	p.LeftPage.Renderer = f
	for ind := range p.PageList {
		p.PageList[ind].Renderer = f
	}
	p.RightPage.Renderer = f
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
