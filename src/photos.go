package main

import "log"

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
}

type PhotoWallet struct {
	wallet map[string]*Photo
}

func NewPhotoWallet(j *Entry) *PhotoWallet {
	res := &PhotoWallet{
		wallet: make(map[string]*Photo, len(j.Photos)),
	}

	for _, photo := range j.Photos {
		res.wallet[photo.Identifier] = photo
	}

	return res
}

func (pw *PhotoWallet) fixPhotoSrc(given string) string {
	p := pw.wallet[(given)]
	if p == nil {
		log.Printf("warning: no photo found for %q", (given))
		return given
	}

	return p.fixPhotoSrc(given)
}

func (p *Photo) fixPhotoSrc(bgiven string) string {
	res := p.MD5

	res += "." + p.Type

	return res
}
