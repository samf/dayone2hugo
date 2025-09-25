package main

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"
)

const (
	journalName = "Journal.json"
)

type DOZip struct {
	reader *zip.Reader
	dir    map[string]*zip.File
}

func NewDOZip(input *os.File) (*DOZip, error) {
	fi, err := input.Stat()
	if err != nil {
		err = fmt.Errorf("cannot stat input file: %w", err)
		return nil, err
	}
	reader, err := zip.NewReader(input, fi.Size())
	if err != nil {
		err = fmt.Errorf("cannot open zip reader: %w", err)
		return nil, err
	}

	dir := map[string]*zip.File{}
	for _, f := range reader.File {
		dir[f.Name] = f
	}

	res := &DOZip{
		reader: reader,
		dir:    dir,
	}

	return res, nil
}

func (doz *DOZip) getJournal() (*Journal, error) {
	jfile := doz.dir[journalName]
	if jfile == nil {
		return nil, fmt.Errorf("no %v file found", journalName)
	}

	journal, err := NewDOJournal(jfile)
	if err != nil {
		err = fmt.Errorf("cannot parse %v: %w", journalName, err)
		return nil, err
	}

	return journal, nil
}

func (doz *DOZip) findPhotoFile(name string) (*zip.File, error) {
	for n, f := range doz.dir {
		if strings.Contains(n, name) {
			return f, nil
		}
	}

	return nil, fmt.Errorf("could not find zip component containing %q",
		name,
	)
}
