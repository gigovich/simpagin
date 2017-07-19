// Use of this source code is governed by BSD license
// 2014. Moscow. Givi Khojanashvili <gigovich@gmail.com>

// Package simpagin realizes simple pagination tool.
package simpagin

import (
	"strconv"
	"strings"
)

const (
	// PageLeft in paginator.
	PageLeft = iota

	// PageMiddle in paginator.
	PageMiddle

	// PageRight in paginator.
	PageRight
)

// PageRenderer should return builded HTML.
type PageRenderer func(p Page) string

// Page object contains all attributes for render them in paginator.
type Page struct {
	Index    int          // Object index (position) in total object list
	Number   int          // Page number valid value is from range [1..(itemsCount / frameLength)]
	IsActive bool         // Is this page active, paginator contains only one active page
	Type     int          // PageMiddle is all pages in frame, PageLeft and PageRight are scroller pages
	Renderer PageRenderer // Custom function to render page element as string called by Page.Render()
}

// String calls Page.Renderer to render string representation of page or some other page data
//
// If Renderer does not set, empty string will be returned. You can set
// Renderer function for all pages by Paginator.SetRenderer method.
func (p Page) String() string {
	if p.Renderer != nil {
		return p.Renderer(p)
	}
	if p.Number > 0 {
		return strconv.Itoa(p.Number)
	}
	return ""
}

// Paginator is main struct which renders pages.
type Paginator struct {
	ActivePage  int     // Active page number
	LeftPage    *Page   // Page for left scroller, if active page is too close to start it must be nil
	RightPage   *Page   // Page for right scroller, if active page is too close to end it must be nil
	ItemsCount  int     // Total items count
	PagesCount  int     // Auto calculated field which equals to ItemsCount / ItemsOnPage
	ItemsOnPage int     // How much items contains each page
	FrameLength int     // Number of pages displayed in paginator
	PageList    []*Page // You must fetch this slice to display each paginated page
}

// New returns new Paginator struct, with calculated fields which you can use,
// to render paginator. If 'activePage' argument is wrong (less then 1 or items on this page out of itemsCount) it
// will be reset to 1. If 'frameLength' is less then 2, it will be reset to 2.
//
// Exemple of usage:
// 	pg := New(
// 		10,  // Active page which items we displaying now
// 		120, // Total count of items
// 		8,   // We show only 8 items in each page
// 		10,  // And our paginator rendered as 10 pages list
// 	)
// 	fmt.Printf("<a href=\"/page/%s/\">&lt</a>", pg.LeftPage.Number)
// 	for _, page := range pg.PageList {
// 		if page.IsActive {
// 			fmt.Print(page.Number)
// 		} else {
// 			fmt.Printf("<a href=\"/page/%s/\">&lt</a>", pg.LeftPage.Number)
// 		}
// 	}
// 	fmt.Printf("<a href=\"/page/%s/\">&gt</a>", pg.RightPage.Number)
func New(activePage, itemsCount, itemsOnPage, frameLength int) *Paginator {
	if itemsOnPage < 1 {
		itemsOnPage = 1
	}

	// Calculate PagesCount
	pagesCount := itemsCount / itemsOnPage
	if itemsCount%itemsOnPage != 0 {
		pagesCount++
	}

	if activePage < 1 {
		activePage = 1
	}

	if activePage > pagesCount {
		activePage = 1
	}

	if frameLength < 2 {
		frameLength = 2
	}

	pg := &Paginator{
		ActivePage:  activePage,
		ItemsCount:  itemsCount,
		ItemsOnPage: itemsOnPage,
		FrameLength: frameLength,
		LeftPage:    &Page{Type: PageLeft},
		RightPage:   &Page{Type: PageRight},
		PagesCount:  pagesCount,
	}

	// Calculate side indexes
	distanceToLeftSide := (frameLength / 2)
	distanceToRightSide := frameLength - distanceToLeftSide
	frameStartIndex := 1
	if activePage > distanceToLeftSide+1 {
		frameStartIndex = activePage - distanceToLeftSide
		if activePage > pagesCount-distanceToRightSide {
			frameStartIndex -= activePage - (pagesCount - distanceToRightSide) - 1
		}
	}
	if activePage > 1 {
		pg.LeftPage = &Page{(activePage - 1) * itemsOnPage, activePage - 1, false, PageLeft, nil}
	}
	if activePage < pagesCount {
		pg.RightPage = &Page{(activePage + 1) * itemsOnPage, activePage + 1, false, PageRight, nil}
	}
	pages := make([]*Page, min(frameLength, pagesCount))
	for i := 0; i < len(pages); i++ {
		pageNumber := i + frameStartIndex
		pages[i] = &Page{
			Index:  (pageNumber - 1) * itemsOnPage,
			Number: pageNumber,
			Type:   PageMiddle,
		}
		if pageNumber == activePage {
			pages[i].IsActive = true
		}
	}
	pg.PageList = pages
	return pg
}

// SetRenderer set for all page object in the paginator PageRenderer function. This method returns paginator itself,
// so you can do chaincalls.
//
// Exemple of usage:
// 	pg := simpagin.New(
// 		10,  // Active page which items we displaying now
// 		120, // Total count of items
// 		8,   // We show only 8 items in each page
// 		10,  // And our paginator rendered as 10 pages list
// 	)
// 	pg.SetRenderer(func (p Page) string {
// 		switch p.Type {
// 		case simpagin.PageLeft:
// 			if p.Number == 0 {
// 				return "<li class=\"disabled\"><span>&laquo;</span></li>"
// 			} else {
// 				return fmt.Sprintf("<li><a href=\"?p=%d\">&laquo;</a></li>", p.Number)
// 			}
// 		case simpagin.PageMiddle:
// 			if p.IsActive {
// 				return fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", p.Number)
// 			} else {
// 				return fmt.Sprintf("<li><a href=\"?p=%d\">%d</a></li>", p.Number, p.Number)
// 			}
// 		case simpagin.PageRight:
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
func (p *Paginator) SetRenderer(f PageRenderer) *Paginator {
	p.LeftPage.Renderer = f
	for ind := range p.PageList {
		p.PageList[ind].Renderer = f
	}
	p.RightPage.Renderer = f
	return p
}

// Render collects first, all middle pages and last page render results, concat them and returns as string.
//
// Exemple of usage:
//  pageRender := func (p Page) string {
// 		switch p.Type {
// 		case simpagin.PageLeft:
// 			if p.Number == 0 {
// 				return "<li class=\"disabled\"><span>&laquo;</span></li>"
// 			}
// 			return fmt.Sprintf("<li><a href=\"?p=%d\">&laquo;</a></li>", p.Number)
// 		case simpagin.PageMiddle:
// 			if p.IsActive {
// 				return fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", p.Number)
// 			}
// 			return fmt.Sprintf("<li><a href=\"?p=%d\">%d</a></li>", p.Number, p.Number)
// 		case simpagin.PageRight:
// 			if p.Number == 0 {
// 				return "<li class=\"disabled\"><span>&raquo;</span></li>"
// 			}
// 			return fmt.Sprintf("<li><a href=\"?p=%d\">&raquo;</a></li>", p.Number, p.Number)
// 		}
// 		return ""
// 	}
//
// 	renderedPaginator := simpagin.New(
// 		10,  // Active page which items we displaying now
// 		120, // Total count of items
// 		8,   // We show only 8 items in each page
// 		10,  // And our paginator rendered as 10 pages list
// 	).SetRenderer(pageRender).Render()
//  fmt.Println(renderedPaginator)
func (p *Paginator) Render() string {
	l := make([]string, len(p.PageList)+2)
	l[0] = p.LeftPage.String()
	for i, page := range p.PageList {
		l[i+1] = page.String()
	}
	l[len(l)-1] = p.RightPage.String()
	return strings.Join(l, "")
}

// GetIndex of the first item on the page.
func (p *Paginator) GetIndex() int {
	if p.ActivePage-1 >= len(p.PageList) {
		return 0
	}
	return p.PageList[p.ActivePage-1].Index
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
