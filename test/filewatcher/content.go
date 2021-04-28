package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// content of a .content file
type content struct {
	Coverpagenumber int  `json:"coverPageNumber"`
	Dummydocument   bool `json:"dummyDocument"`
	Extrametadata   struct {
		Lastballpointcolor     string `json:"LastBallpointColor"`
		Lastballpointsize      string `json:"LastBallpointSize"`
		Lastballpointv2Color   string `json:"LastBallpointv2Color"`
		Lastballpointv2Size    string `json:"LastBallpointv2Size"`
		Lastcalligraphycolor   string `json:"LastCalligraphyColor"`
		Lastcalligraphysize    string `json:"LastCalligraphySize"`
		Lastclearpagecolor     string `json:"LastClearPageColor"`
		Lastclearpagesize      string `json:"LastClearPageSize"`
		Lasterasesectioncolor  string `json:"LastEraseSectionColor"`
		Lasterasesectionsize   string `json:"LastEraseSectionSize"`
		Lasterasercolor        string `json:"LastEraserColor"`
		Lasterasersize         string `json:"LastEraserSize"`
		Lasterasertool         string `json:"LastEraserTool"`
		Lastfinelinercolor     string `json:"LastFinelinerColor"`
		Lastfinelinersize      string `json:"LastFinelinerSize"`
		Lastfinelinerv2Color   string `json:"LastFinelinerv2Color"`
		Lastfinelinerv2Size    string `json:"LastFinelinerv2Size"`
		Lasthighlightercolor   string `json:"LastHighlighterColor"`
		Lasthighlightersize    string `json:"LastHighlighterSize"`
		Lasthighlighterv2Color string `json:"LastHighlighterv2Color"`
		Lasthighlighterv2Size  string `json:"LastHighlighterv2Size"`
		Lastmarkercolor        string `json:"LastMarkerColor"`
		Lastmarkersize         string `json:"LastMarkerSize"`
		Lastmarkerv2Color      string `json:"LastMarkerv2Color"`
		Lastmarkerv2Size       string `json:"LastMarkerv2Size"`
		Lastpaintbrushcolor    string `json:"LastPaintbrushColor"`
		Lastpaintbrushsize     string `json:"LastPaintbrushSize"`
		Lastpaintbrushv2Color  string `json:"LastPaintbrushv2Color"`
		Lastpaintbrushv2Size   string `json:"LastPaintbrushv2Size"`
		Lastpen                string `json:"LastPen"`
		Lastpencilcolor        string `json:"LastPencilColor"`
		Lastpencilsize         string `json:"LastPencilSize"`
		Lastpencilv2Color      string `json:"LastPencilv2Color"`
		Lastpencilv2Size       string `json:"LastPencilv2Size"`
		Lastreservedpencolor   string `json:"LastReservedPenColor"`
		Lastreservedpensize    string `json:"LastReservedPenSize"`
		Lastselectiontoolcolor string `json:"LastSelectionToolColor"`
		Lastselectiontoolsize  string `json:"LastSelectionToolSize"`
		Lastsharppencilcolor   string `json:"LastSharpPencilColor"`
		Lastsharppencilsize    string `json:"LastSharpPencilSize"`
		Lastsharppencilv2Color string `json:"LastSharpPencilv2Color"`
		Lastsharppencilv2Size  string `json:"LastSharpPencilv2Size"`
		Lastsolidpencolor      string `json:"LastSolidPenColor"`
		Lastsolidpensize       string `json:"LastSolidPenSize"`
		Lasttool               string `json:"LastTool"`
		Lastundefinedcolor     string `json:"LastUndefinedColor"`
		Lastundefinedsize      string `json:"LastUndefinedSize"`
		Lastzoomtoolcolor      string `json:"LastZoomToolColor"`
		Lastzoomtoolsize       string `json:"LastZoomToolSize"`
	} `json:"extraMetadata"`
	Filetype      string   `json:"fileType"`
	Fontname      string   `json:"fontName"`
	Lineheight    int      `json:"lineHeight"`
	Margins       int      `json:"margins"`
	Orientation   string   `json:"orientation"`
	Pagecount     int      `json:"pageCount"`
	Pages         []string `json:"pages"`
	Textalignment string   `json:"textAlignment"`
	Textscale     int      `json:"textScale"`
	Transform     struct {
		M11 int `json:"m11"`
		M12 int `json:"m12"`
		M13 int `json:"m13"`
		M21 int `json:"m21"`
		M22 int `json:"m22"`
		M23 int `json:"m23"`
		M31 int `json:"m31"`
		M32 int `json:"m32"`
		M33 int `json:"m33"`
	} `json:"transform"`
}

type metadata struct {
	Deleted          bool   `json:"deleted"`
	Lastmodified     string `json:"lastModified"`
	Lastopenedpage   int    `json:"lastOpenedPage"`
	Metadatamodified bool   `json:"metadatamodified"`
	Modified         bool   `json:"modified"`
	Parent           string `json:"parent"`
	Pinned           bool   `json:"pinned"`
	Synced           bool   `json:"synced"`
	Type             string `json:"type"`
	Version          int    `json:"version"`
	Visiblename      string `json:"visibleName"`
}

func findMostRecent(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var modTime time.Time
	var names []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name())))
			}
		}
	}
	if len(names) == 1 {
		return names[0], nil
	}
	return "", errors.New("expected only one result")
}
