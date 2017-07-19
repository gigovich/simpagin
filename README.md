Simpagin - simple pagination tool
=================================
Simple lib for pagination like google.

Install package:
```bash
    go get github.com/gigovich/simpagin
```

## Create paginator
This is paginator instance constructor. Remeber that, page numbers started from 1.
```go
func New(activePage, itemsCount, itemsOnPage, frameLength int) *Paginator
```
* **activePage**, **itemsOnPage** and **frameLength** will be reseted to default values if they are less (1, 1, 2).
* **itemsCount** - 0 is valid value for this param.

### Example - most simple
```go
pg := simpagin.New(
    10,  // Active page which items we displaying now
    120, // Total count of items
    8,   // We show only 8 items in each page
    10,  // And our paginator rendered as 10 pages list
)
// Show paginator
fmt.Printf(`<a href="/page/%s/">&lt</a>`, pg.LeftPage.Number)
for _, page := range pg.PageList {
    if page.IsActive {
        fmt.Print(page.Number)
    }
    fmt.Printf(`<a href="/page/%s/">&lt</a>`, pg.LeftPage.Number)
}
fmt.Printf(`<a href="/page/%s/">&gt</a>`, pg.RightPage.Number)
```

### Example 2 - customize page render method
```go
pg := simpagin.New(
    10,  // Active page which items we displaying now
    120, // Total count of items
    8,   // We show only 8 items in each page
    10,  // And our paginator rendered as 10 pages list
)

pg.SetRenderer(func (p Page) string {
    switch p.Type {
    case simpagin.PageLeft:
        if p.Number == 0 {
            return `<li class="disabled"><span>&laquo;</span></li>`
        }
        return fmt.Sprintf(`<li><a href="?p=%d">&laquo;</a></li>`, p.Number)
    case simpagin.PageMiddle:
        if p.IsActive {
            return fmt.Sprintf(`<li class="active"><span>%d</span></li>`, p.Number)
        }
        return fmt.Sprintf(`<li><a href="?p=%d">%d</a></li>`, p.Number, p.Number)
    case simpagin.PageRight:
        if p.Number == 0 {
            return `<li class="disabled"><span>&raquo;</span></li>`
        }
        return fmt.Sprintf(`<li><a href="?p=%d">&raquo;</a></li>`, p.Number)
    }
    return ""
})

// Now each page element know how to render himself
fmt.Print(pg.LeftPage)
for _, page := range pg.PageList {
    fmt.Print(page)
}

fmt.Print(pg.RightPage)
```

### Example 3 - customize page render and use navigator
```go
pageRender := func (p Page) string {
	switch p.Type {
	case simpagin.PageLeft:
		if p.Number == 0 {
			return `<li class="disabled"><span>&laquo;</span></li>`
		}
		return fmt.Sprintf(`<li><a href="?p=%d">&laquo;</a></li>`, p.Number)
	case simpagin.PageMiddle:
		if p.IsActive {
			return fmt.Sprintf(`<li class="active"><span>%d</span></li>`, p.Number)
		}
		return fmt.Sprintf(`<li><a href="?p=%d">%d</a></li>`, p.Number, p.Number)
	case simpagin.PageRight:
		if p.Number == 0 {
			return `<li class="disabled"><span>&raquo;</span></li>`
		}
		return fmt.Sprintf(`<li><a href="?p=%d">&raquo;</a></li>`, p.Number)
	}
	return ""
}

renderedPaginator := simpagin.New(
	10,  // Active page which items we displaying now
	120, // Total count of items
	8,   // We show only 8 items in each page
	10,  // And our paginator rendered as 10 pages list
).SetRenderer(pageRender).Render()

fmt.Println(renderedPaginator)
```

## Use with html/template
Golang `html/template` package escapes raw strings, so you should use `template.HTML` to wrap rendered paginator string.
