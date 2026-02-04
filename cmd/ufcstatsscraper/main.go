package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// FighterSummary is what we pull from the fighter index / listing pages.
type FighterSummary struct {
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Nickname  string        `json:"nickname"`
	Height    string        `json:"height"`
	Weight    string        `json:"weight"`
	Reach     string        `json:"reach"`
	Stance    string        `json:"stance"`
	Wins      string        `json:"wins"`
	Losses    string        `json:"losses"`
	Draws     string        `json:"draws"`
	DetailURL string        `json:"detail_url"`
	Details   FighterDetail `json:"details,omitempty"`
}

// FighterDetail is the full profile scraped from an individual fighter page.
type FighterDetail struct {
	Name   string `json:"name"`
	Record string `json:"record"`

	// Physical
	Height string `json:"height"`
	Weight string `json:"weight"`
	Reach  string `json:"reach"`
	Stance string `json:"stance"`
	DOB    string `json:"dob"`

	// Career stats
	SLpM      string `json:"slpm"`       // Significant Strikes Landed per Minute
	StrikeAcc string `json:"strike_acc"` // Striking Accuracy
	SApM      string `json:"sapm"`       // Significant Strikes Absorbed per Minute
	StrikeDef string `json:"strike_def"` // Strike Defence %
	TDAvg     string `json:"td_avg"`     // Takedown Avg per 15 min
	TDAcc     string `json:"td_acc"`     // Takedown Accuracy
	TDDef     string `json:"td_def"`     // Takedown Defence %
	SubAvg    string `json:"sub_avg"`    // Submission Avg per 15 min

	Fights []FightRecord `json:"fights"`
}

// FightRecord is a single row in the fighter's fight history table.
type FightRecord struct {
	Result   string `json:"result"` // "win" | "loss" | "draw" | "nc"
	Opponent string `json:"opponent"`
	Event    string `json:"event"`
	Date     string `json:"date"`
	Method   string `json:"method"`
	Round    string `json:"round"`
	Time     string `json:"time"`
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

// ---------------------------------------------------------------------------
// Scraper: fighter index (listing page)
// ---------------------------------------------------------------------------

// scrapeFighterIndex fetches a single "char" page (e.g. char=a) and returns
// every FighterSummary on that page
func scrapeFighterIndex(c *colly.Collector, char string) ([]FighterSummary, error) {
	url := fmt.Sprintf("http://ufcstats.com/statistics/fighters?char=%s&page=all", char)
	var fighters = make([]FighterSummary, 0)

	// Each fighter is one <tr> inside <tbody>.
	// The first <tr> is an empty spacer row — we skip it by checking for links.
	c.OnHTML("table.b-statistics__table tbody tr.b-statistics__table-row", func(e *colly.HTMLElement) {
		cols := e.ChildAttrs("td.b-statistics__table-col a.b-link", "href")
		if len(cols) == 0 {
			return // spacer row
		}

		// Pull every <td> in order (matches the <thead> columns).
		var tds []string
		e.ForEach("td", func(_ int, sel *colly.HTMLElement) {
			tds = append(tds, clean(sel.Text))
		})

		// Minimum 11 columns expected; guard against malformed rows.
		if len(tds) < 11 {
			return
		}

		fighters = append(fighters, FighterSummary{
			FirstName: tds[0],
			LastName:  tds[1],
			Nickname:  tds[2],
			Height:    tds[3],
			Weight:    tds[4],
			Reach:     tds[5],
			Stance:    tds[6],
			Wins:      tds[7],
			Losses:    tds[8],
			Draws:     tds[9],
			// tds[10] = Belt (empty for most fighters)
			DetailURL: cols[0], // first link href is the fighter detail page
		})
	})

	err := c.Visit(url)
	return fighters, err
}

// scrapeFighterDetail fetches a single fighter's detail page and returns a populated FighterDetail.
func scrapeFighterDetail(c *colly.Collector, url string) (FighterDetail, error) {
	var f FighterDetail

	// Name & Record
	c.OnHTML("h2.b-content__title", func(e *colly.HTMLElement) {
		e.ForEach("span", func(i int, sel *colly.HTMLElement) {
			switch i {
			case 0:
				f.Name = clean(sel.Text)
			case 1:
				// "Record: 5-3-0"  →  "5-3-0"
				f.Record = strings.TrimPrefix(clean(sel.Text), "Record: ")
			}
		})
	})

	// Physical stats (left info box)
	// Each <li> has a title <i> and a text value.
	c.OnHTML("div.b-list__info-box_style_small-width li.b-list__box-list-item", func(e *colly.HTMLElement) {
		title := clean(e.ChildText("i.b-list__box-item-title"))
		// The value is the remaining text after the <i>; grab full text and trim the title.
		value := clean(strings.TrimPrefix(clean(e.Text), title))

		switch {
		case strings.HasPrefix(title, "Height"):
			f.Height = value
		case strings.HasPrefix(title, "Weight"):
			f.Weight = value
		case strings.HasPrefix(title, "Reach"):
			f.Reach = value
		case strings.HasPrefix(title, "STANCE"):
			f.Stance = value
		case strings.HasPrefix(title, "DOB"):
			f.DOB = value
		}
	})

	// Career statistics
	// These live inside the middle-width info box, split into left / right sub-boxes.
	c.OnHTML("div.b-list__info-box_style_middle-width .b-list__info-box-left li.b-list__box-list-item", func(e *colly.HTMLElement) {
		title := clean(e.ChildText("i.b-list__box-item-title"))
		value := clean(strings.TrimPrefix(clean(e.Text), title))

		switch {
		case strings.HasPrefix(title, "SLpM"):
			f.SLpM = value
		case strings.HasPrefix(title, "Str. Acc"):
			f.StrikeAcc = value
		case strings.HasPrefix(title, "SApM"):
			f.SApM = value
		case strings.HasPrefix(title, "Str. Def"):
			f.StrikeDef = value
		case strings.HasPrefix(title, "TD Avg"):
			f.TDAvg = value
		case strings.HasPrefix(title, "TD Acc"):
			f.TDAcc = value
		case strings.HasPrefix(title, "TD Def"):
			f.TDDef = value
		case strings.HasPrefix(title, "Sub. Avg"):
			f.SubAvg = value
		}
	})

	// Fight history table
	// Each <tr> (except the spacer) maps to one fight.
	c.OnHTML("table.b-fight-details__table tbody tr.b-fight-details__table-row", func(e *colly.HTMLElement) {
		// Result flag class: b-flag_style_green = win, b-flag_style_bordered = loss, etc.
		resultText := clean(e.ChildText("td a.b-flag i.b-flag__inner i.b-flag__text"))
		if resultText == "" {
			return // spacer row
		}

		// Opponent is the second <a> inside the Fighter column.
		var opponent string
		links := e.ChildAttrs("td.b-fight-details__table-col a.b-link", "href")
		_ = links // we only need the display text
		e.ForEach("td.b-fight-details__table-col a.b-link", func(i int, sel *colly.HTMLElement) {
			if i == 1 {
				opponent = clean(sel.Text)
			}
		})

		event := clean(e.ChildText("td a.b-link[href*='event-details']"))

		// Date is the second <p> inside the Event column — grab all <p> texts in that td.
		var date string
		e.ForEach("td.b-fight-details__table-col", func(_ int, td *colly.HTMLElement) {
			// The event td contains both an <a> with the event name and a bare <p> with the date.
			if strings.Contains(td.Attr("class"), "") {
				ps := td.ChildTexts("p.b-fight-details__table-text")
				for _, p := range ps {
					p = clean(p)
					if p != "" && p != event && date == "" {
						// Looks like a date if it contains a comma or a month abbreviation.
						if strings.Contains(p, ",") {
							date = p
						}
					}
				}
			}
		})

		// Method column has the method type on the first <p> and details on the second.
		var method string
		// Walk all tds; the method td is the one whose first <p> matches known keywords.
		e.ForEach("td.b-fight-details__table-col", func(_ int, td *colly.HTMLElement) {
			texts := td.ChildTexts("p.b-fight-details__table-text")
			if len(texts) >= 1 {
				first := clean(texts[0])
				switch first {
				case "KO (Punch)", "KO", "TKO (Punches)", "TKO", "SUB", "U-DEC", "S-DEC", "M-DEC", "DQ", "NC", "OVERRULE":
					method = first
					if len(texts) > 1 && clean(texts[1]) != "" {
						method += " (" + clean(texts[1]) + ")"
					}
				}
			}
		})

		// Round & Time are the last two single-<p> columns.
		var round, time string
		e.ForEach("td.b-fight-details__table-col", func(i int, td *colly.HTMLElement) {
			t := clean(td.Text)
			// Round is a bare single digit; Time contains a colon.
			if t != "" && len(t) <= 2 && t != "W/L" {
				if round == "" && t[0] >= '1' && t[0] <= '9' {
					round = t
				}
			}
			if strings.Contains(t, ":") && len(t) <= 5 {
				time = t
			}
		})

		f.Fights = append(f.Fights, FightRecord{
			Result:   resultText,
			Opponent: opponent,
			Event:    event,
			Date:     date,
			Method:   method,
			Round:    round,
			Time:     time,
		})
	})

	err := c.Visit(url)
	return f, err
}

// ---------------------------------------------------------------------------
// Collector factory — shared setup so both scrapers behave consistently.
// ---------------------------------------------------------------------------

func newCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
	)

	// Rotate user-agents so we don't get blocked immediately.
	extensions.RandomUserAgent(c)

	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*ufcstats.com*",
		Parallelism: 4,
		Delay:       20 * time.Millisecond,
	}); err != nil {
		log.Fatal(err)
	}

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("HTTP %d  %s  → %v", r.StatusCode, r.Request.URL.String(), err)
	})

	return c
}

func main() {
	indexCollector := newCollector()

	for i := 'a'; i <= 'z'; i++ {
		char := string(i)
		fighters, err := scrapeFighterIndex(indexCollector, char)
		if err != nil {
			log.Printf("error scraping index for %s: %v", char, err)
		}

		// ── 2. Pick the first fighter and scrape their detail page ──────────────
		if len(fighters) == 0 {
			log.Printf("no fighters scraped — nothing to do")
		}

		detailCollector := newCollector()

		// open file to write fighter details
		file, err := os.OpenFile(fmt.Sprintf("fighter_details_%s.json", string(i)), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Scrape details for every fighter.
		for i, f := range fighters {
			detail, err := scrapeFighterDetail(detailCollector, f.DetailURL)
			if err != nil {
				log.Printf("error scraping details for %s: %v", f.DetailURL, err)
			}

			fighters[i].Details = detail
		}

		log.Printf("scraped details for %d fighters", len(fighters))

		// write fighters to file as JSON
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(fighters); err != nil {
			log.Printf("error encoding fighters to JSON: %v", err)
		}

		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}
}
