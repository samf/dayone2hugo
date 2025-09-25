package main

import (
	"fmt"
	"log"
	"strings"
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
}

type PhotoWallet struct {
	wallet        map[string]*Photo
	entryPhotos   []*Photo
	friendlyNames []string
}

func NewPhotoWallet(j *Entry) *PhotoWallet {
	res := &PhotoWallet{
		wallet:      make(map[string]*Photo, len(j.Photos)),
		entryPhotos: make([]*Photo, len(j.Photos)),
	}

	for i, photo := range j.Photos {
		res.wallet[photo.Identifier] = photo
		res.entryPhotos[i] = photo
	}

	return res
}

func (pw *PhotoWallet) setFriendlyNames(names []string) error {
	if len(names) == 0 {
		return nil
	}

	if len(names) != len(pw.wallet) {
		return fmt.Errorf(
			"there are %v photos but you have %v names; these must match",
			len(pw.wallet),
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
	p := pw.wallet[given]
	if p == nil {
		log.Printf("warning: no photo found for %q", (given))
		return given
	}

	return p.getFName()
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
