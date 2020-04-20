package main

import "encoding/xml"

type CrosswordCompiler struct {
	XMLName           xml.Name `xml:"crossword-compiler"`
	Text              string   `xml:",chardata"`
	Xmlns             string   `xml:"xmlns,attr"`
	RectangularPuzzle struct {
		Text     string `xml:",chardata"`
		Xmlns    string `xml:"xmlns,attr"`
		Alphabet string `xml:"alphabet,attr"`
		Metadata struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Creator     string `xml:"creator"`
			Copyright   string `xml:"copyright"`
			Description string `xml:"description"`
		} `xml:"metadata"`
		Crossword struct {
			Text string `xml:",chardata"`
			Grid struct {
				Text     string `xml:",chardata"`
				Width    string `xml:"width,attr"`
				Height   string `xml:"height,attr"`
				GridLook struct {
					Text                   string `xml:",chardata"`
					ThickBorder            string `xml:"thick-border,attr"`
					NumberingScheme        string `xml:"numbering-scheme,attr"`
					CellSizeInPixels       string `xml:"cell-size-in-pixels,attr"`
					ClueSquareDividerWidth string `xml:"clue-square-divider-width,attr"`
					Arrows                 struct {
						Text           string `xml:",chardata"`
						Stem           string `xml:"stem,attr"`
						BendStart      string `xml:"bend-start,attr"`
						BendEnd        string `xml:"bend-end,attr"`
						BendSideOffset string `xml:"bend-side-offset,attr"`
						HeadWidth      string `xml:"head-width,attr"`
						HeadLength     string `xml:"head-length,attr"`
					} `xml:"arrows"`
				} `xml:"grid-look"`
				Cells []Cell `xml:"cell"`
			} `xml:"grid"`
			Words []Word `xml:"word"`
		} `xml:"crossword"`
	} `xml:"rectangular-puzzle"`
}

type Cell struct {
	Text     string `xml:",chardata"`
	X        string `xml:"x,attr"`
	Y        string `xml:"y,attr"`
	Type     string `xml:"type,attr"`
	Solution string `xml:"solution,attr"`
	Clue     []struct {
		Text string `xml:",chardata"`
		Word string `xml:"word,attr"`
	} `xml:"clue"`
	Arrow []struct {
		Text string `xml:",chardata"`
		From string `xml:"from,attr"`
		To   string `xml:"to,attr"`
	} `xml:"arrow"`
}

type Word struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	X    string `xml:"x,attr"`
	Y    string `xml:"y,attr"`
}