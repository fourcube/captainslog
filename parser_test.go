package main_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	c "github.com/fourcube/captainslog"
)

func TestParsesEmptyFile(t *testing.T) {
	res := c.Parse("")

	if len(res) != 0 {
		t.Errorf("Expected no entries, got %v", res)
	}

}

func TestParseSingleEntry(t *testing.T) {
	sample := `## May 26, 2015 at 7:31pm (CEST)

Employee assignments planarc project are also colored by team now. What I have to figure out is how to chain the loading of the data. Aggregate stores?`

	res := c.Parse(sample)

	if len(res) != 1 {
		t.Errorf("Expected 1 entry, got %v", res)
	}

	if len(res[0].Lines) < 1 {
		t.Errorf("Expected log entry to have text. Has none.")
	}
}

func TestParseMultipleEntries(t *testing.T) {
	sample := `## May 26, 2015 at 7:31pm (CEST)

Employee assignments planarc project are also colored by team now. What I have to figure out is how to chain the loading of the data. Aggregate stores?

## May 26, 2016 at 7:31pm (CEST)

Employee assignments planarc project are also colored by team now. What I have to figure out is how to chain the loading of the data. Aggregate stores?`

	res := c.Parse(sample)

	if len(res) != 2 {
		t.Errorf("Expected 1 entry, got %v", res)
	}

	if len(res[0].Lines) < 1 {
		t.Errorf("Expected log entry to have text. Has none.")
	}

	if len(res[1].Lines) < 1 {
		t.Errorf("Expected log entry to have text. Has none.")
	}

	if res[1].Year() != 2016 {
		t.Errorf("Expected log entry to have year 2016. Has %v.", res[1].Year())
	}
}

func TestParseCaptainslogFile(t *testing.T) {
	logs, err := ioutil.ReadFile(os.Getenv("CAPTAINSLOG"))
	if err != nil {
		t.Skip("No $CAPTAINSLOG file")
	}

	res := c.Parse(string(logs))
	if len(res) < 1 {
		t.Errorf("That failed badly.")
	}

	for _, r := range res {
		log.Printf("Date: %v\nText: %s", r.Time, r.Lines)
	}
}
