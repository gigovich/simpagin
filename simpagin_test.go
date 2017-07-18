package simpagin

import (
	"fmt"
	"testing"
)

func TestWrongActivePage(t *testing.T) {
	_, err := New(2, 10, 10, 10)
	if err == nil {
		t.Fail()
	}
}

func TestRightValuesForSides(t *testing.T) {
	paginator, err := New(11, 160, 8, 7)
	if err != nil {
		t.Error(err)
	}
	if paginator.LeftPage.Number != 10 {
		t.Errorf("Left side page number: 10 != %d", paginator.LeftPage.Number)
	}
	if paginator.LeftPage.Index != 80 {
		t.Errorf("Left side page index: 80 != %d", paginator.LeftPage.Index)
	}
	if paginator.RightPage.Number != 12 {
		t.Errorf("Right side page number: 12 != %d", paginator.RightPage.Number)
	}
	if paginator.RightPage.Index != 96 {
		t.Errorf("Right side page index: 96 != %d", paginator.RightPage.Index)
	}
}

func TestPagesValues(t *testing.T) {
	paginator, err := New(11, 160, 8, 7)
	if err != nil {
		t.Error(err)
	}
	rightValues := []Page{
		{56, 8, false, PageMiddle, nil},
		{64, 9, false, PageMiddle, nil},
		{72, 10, false, PageMiddle, nil},
		{80, 11, true, PageMiddle, nil},
		{88, 12, false, PageMiddle, nil},
		{96, 13, false, PageMiddle, nil},
		{104, 14, false, PageMiddle, nil},
	}
	for ind := range paginator.PageList {
		if !comp(paginator.PageList[ind], &rightValues[ind]) {
			t.Errorf("Values %v != %v", paginator.PageList[ind], rightValues[ind])
		}
	}
}

func TestRenderer(t *testing.T) {
	pg, err := New(
		10,  // Active page which items we already display
		120, // Total count of items
		8,   // We show only 8 items in each page
		10,  // And our paginator rendered as 10 pages list
	)
	if err != nil {
		t.Error(err)
	}
	pg.SetRenderer(func(p Page) string {
		switch p.Type {
		case PageLeft:
			if p.Number == 0 {
				return "<li class=\"disabled\"><span>&laquo;</span></li>"
			} else {
				return fmt.Sprintf("<li><a href=\"?p=%d\">&laquo;</a></li>", p.Number)
			}
		case PageMiddle:
			if p.IsActive {
				return fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", p.Number)
			} else {
				return fmt.Sprintf("<li><a href=\"?p=%d\">%d</a></li>", p.Number, p.Number)
			}
		case PageRight:
			if p.Number == 0 {
				return "<li class=\"disabled\"><span>&raquo;</span></li>"
			} else {
				return fmt.Sprintf("<li><a href=\"?p=%d\">&raquo;</a></li>", p.Number, p.Number)
			}
		}
		return ""
	})
	if pg.LeftPage.String() != "<li><a href=\"?p=9\">&laquo;</a></li>" {
		t.Error(pg.LeftPage)
	}
	for _, page := range pg.PageList {
		if page.IsActive {
			if page.String() != fmt.Sprintf(
				"<li class=\"active\"><span>%d</span></li>", page.Number) {
				t.Error(page)
			}
		} else {
			if page.String() != fmt.Sprintf(
				"<li><a href=\"?p=%d\">%d</a></li>", page.Number, page.Number) {
				t.Error(page)
			}
		}
	}
	if pg.RightPage.String() != "<li class=\"disabled\"><span>&raquo;</span></li>" {
		t.Error(pg.RightPage)
	}
}

func comp(a, b *Page) bool {
	if a.Index != b.Index || a.Number != b.Number || a.IsActive != b.IsActive {
		return false
	}
	return true
}
