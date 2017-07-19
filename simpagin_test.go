package simpagin

import (
	"fmt"
	"testing"
)

func TestWrongActivePage(t *testing.T) {
	pg := New(0, 10, 0, 0)
	t.Run("reset active page", func(t *testing.T) {
		if pg.ActivePage == 0 {
			t.Fail()
		}
	})

	t.Run("reset items on page", func(t *testing.T) {
		if pg.ItemsOnPage == 0 {
			t.Fail()
		}
	})

	t.Run("reset frame length", func(t *testing.T) {
		if pg.FrameLength == 0 {
			t.Fail()
		}
	})
}

func TestNoItems(t *testing.T) {
	if v := New(0, 0, 0, 0).Render(); v != "" {
		t.Log(v)
		t.Fail()
	}
}

func TestRightValuesForSides(t *testing.T) {
	paginator := New(11, 160, 8, 7)
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
	paginator := New(11, 160, 8, 7)
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

func TestPageRenderer(t *testing.T) {
	pg := New(
		10,  // Active page which items we already display
		120, // Total count of items
		8,   // We show only 8 items in each page
		10,  // And our paginator rendered as 10 pages list
	)
	pg.SetRenderer(func(p Page) string {
		switch p.Type {
		case PageLeft:
			if p.Number == 0 {
				return "<li class=\"disabled\"><span>&laquo;</span></li>"
			}
			return fmt.Sprintf("<li><a href=\"?p=%d\">&laquo;</a></li>", p.Number)
		case PageMiddle:
			if p.IsActive {
				return fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", p.Number)
			}
			return fmt.Sprintf("<li><a href=\"?p=%d\">%d</a></li>", p.Number, p.Number)
		case PageRight:
			if p.Number == 0 {
				return "<li class=\"disabled\"><span>&raquo;</span></li>"
			}
			return fmt.Sprintf("<li><a href=\"?p=%d\">&raquo;</a></li>", p.Number)
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
	if pg.RightPage.String() != "<li><a href=\"?p=11\">&raquo;</a></li>" {
		t.Error(pg.RightPage)
	}
}

func TestRender(t *testing.T) {
	result := New(11, 150, 3, 5).SetRenderer(func(p Page) string { return fmt.Sprintf("%v,", p.Number) }).Render()
	if result != "10,9,10,11,12,13,12," {
		t.Fail()
	}
}

func TestLastPage(t *testing.T) {
	result := New(8, 15, 2, 3).SetRenderer(func(p Page) string { return fmt.Sprintf("%v,", p.Number) }).Render()
	if result != "7,6,7,8,0," {
		t.Fail()
	}
}

func comp(a, b *Page) bool {
	if a.Index != b.Index || a.Number != b.Number || a.IsActive != b.IsActive {
		return false
	}
	return true
}
