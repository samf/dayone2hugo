package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/barasher/go-exiftool"
)

const (
	CaptionAbstract = "Caption-Abstract"
)

type Photo struct {
	FileSize          int64   `json:"fileSize"`
	OrderInEntry      int     `json:"orderInEntry"`
	CreationDevice    string  `json:"creationDevice"`
	Duration          int     `json:"duration"`
	Favorite          bool    `json:"favorite"`
	Type              string  `json:"type"`
	Identifier        string  `json:"identifier"`
	Date              string  `json:"date"`
	ExposureBiasValue float64 `json:"exposureBiasValue"`
	Height            int     `json:"height"`
	Width             int     `json:"width"`
	MD5               string  `json:"md5"`
	IsSketch          bool    `json:"isSketch"`

	friendlyName string
	caption      string
}

type PhotoWallet struct {
	byID          map[string]*Photo
	byFName       map[string]*Photo
	entryPhotos   []*Photo
	friendlyNames []string
}

func NewPhotoWallet(j *Entry) *PhotoWallet {
	res := &PhotoWallet{
		byID:        make(map[string]*Photo, len(j.Photos)),
		byFName:     make(map[string]*Photo, len(j.Photos)),
		entryPhotos: make([]*Photo, len(j.Photos)),
	}

	for i, photo := range j.Photos {
		res.byID[photo.Identifier] = photo
		res.entryPhotos[i] = photo
	}

	return res
}

func (pw *PhotoWallet) setFriendlyNames(names []string) error {
	if len(names) == 0 {
		return nil
	}

	if len(names) != len(pw.byID) {
		return fmt.Errorf(
			"there are %v photos but you have %v names; these must match",
			len(pw.byID),
			len(names),
		)
	}

	pw.friendlyNames = names

	for i := range names {
		pw.entryPhotos[i].friendlyName = names[i]
	}

	return nil
}

func (pw *PhotoWallet) fixPhotoSrc(given string) string {
	p := pw.byID[given]
	if p == nil {
		log.Printf("warning: no photo found for %q", (given))
		return given
	}

	fname := p.getFName()
	pw.byFName[fname] = p
	return fname
}

func (pw *PhotoWallet) setCaptions(et *exiftool.Exiftool) error {
	for _, p := range pw.entryPhotos {
		err := p.setCaption(et)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Photo) getFName() string {
	res := p.MD5
	if p.friendlyName != "" {
		res = p.friendlyName
	}

	if r := strings.IndexRune(res, '.'); r == -1 {
		res += "." + p.Type
	}

	return res
}

func (p *Photo) setCaption(et *exiftool.Exiftool) error {
	fname := p.getFName()
	md := et.ExtractMetadata(fname)[0]
	if err := md.Err; err != nil {
		err = fmt.Errorf("cannot get metadata for %q: %w", fname, md.Err)
		return err
	}

	var (
		captionVal any
		ok         bool
	)
	if captionVal, ok = md.Fields[CaptionAbstract]; !ok {
		return nil
	}
	p.caption, _ = captionVal.(string)

	return nil
}
