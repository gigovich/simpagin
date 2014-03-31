Simpagin - simple pagination tool
=================================
Simple tool to orginize pagination in your web applications for Golang.
It uses mechanism of pagination like google.

To install package:

    go get github.com/gigovich/simpagin


### Example - most simple


	pg, err := simpagin.New(
		10,  // Active page which items we displaying now
		120, // Total count of items
		8,   // We show only 8 items in each page
		10,  // And our paginator rendered as 10 pages list
	)
	if err != nil {
		log.Fatal(err)
	}
	// Show paginator
	fmt.Printf("<a href=\"/page/%s/\">&lt</a>", pg.LeftPage.Number)
	for _, page := range pg.PageList {
		if page.IsActive {
			fmt.Print(page.Number)
		} else {
			fmt.Printf("<a href=\"/page/%s/\">&lt</a>", pg.LeftPage.Number)
		}
	}
	fmt.Printf("<a href=\"/page/%s/\">&gt</a>", pg.RightPage.Number)


### Example 2 - customize page stringer method


	pg, err := simpagin.New(
		10,  // Active page which items we displaying now
		120, // Total count of items
		8,   // We show only 8 items in each page
		10,  // And our paginator rendered as 10 pages list
	)
	if err != nil {
		log.Fatal(err)
	}
	pg.SetRenderer(func (p Page) string {
		switch p.Type {
		case simpagin.LEFT:
			if p.Number == 0 {
				return "<li class=\"disabled\"><span>&laquo;</span></li>"
			} else {
				return fmt.Sprintf("<li><a href=\"?p=%d\">&laquo;</a></li>", p.Number)
			}
		case simpagin.MIDDLE:
			if p.IsActive {
				return fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", p.Number)
			} else {
				return fmt.Sprintf("<li><a href=\"?p=%d\">%d</a></li>", p.Number, p.Number)
			}
		case simpagin.RIGHT:
			if p.Number == 0 {
				return "<li class=\"disabled\"><span>&raquo;</span></li>"
			} else {
				return fmt.Sprintf("<li><a href=\"?p=%d\">&raquo;</a></li>", p.Number, p.Number)
			}
		}
		return ""
	})
    // Now each page element know how to render himself
	fmt.Print(pg.LeftPage)
	for _, page := range pg.PageList {
		fmt.Print(page)
	}
	fmt.Print(pg.RightPage)
