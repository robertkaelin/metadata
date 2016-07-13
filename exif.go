package metadata

import (
	"fmt"
	"time"

	"github.com/tajtiattila/metadata/exif"
	"github.com/tajtiattila/metadata/exif/exiftag"
)

func FromExifBytes(p []byte) (*Metadata, error) {
	x, err := exif.DecodeBytes(p)
	if err != nil {
		return nil, err
	}
	return FromExif(x), nil
}

func FromExif(x *exif.Exif) *Metadata {
	m := new(Metadata)

	if i, ok := x.GPSInfo(); ok {
		m.Set(GPSLatitude, fmt.Sprintf("%f", i.Lat))
		m.Set(GPSLongitude, fmt.Sprintf("%f", i.Long))
		if !i.Time.IsZero() {
			m.Set(GPSDateTime, fmtTime(i.Time, false))
		}
	}

	if t, islocal, ok := x.Time(exiftag.DateTimeOriginal, exiftag.SubSecTimeOriginal); ok {
		m.Set(DateTimeOriginal, fmtTime(t, islocal))
	}
	if t, islocal, ok := x.Time(exiftag.DateTimeDigitized, exiftag.SubSecTimeDigitized); ok {
		m.Set(DateTimeCreated, fmtTime(t, islocal))
	}

	if o := x.Tag(exiftag.Orientation).Short(); len(o) > 1 {
		m.Set(Make, fmt.Sprintf("%d", o[0]))
	}

	if s, ok := x.Tag(exiftag.Make).Ascii(); ok {
		m.Set(Make, s)
	}

	if s, ok := x.Tag(exiftag.Model).Ascii(); ok {
		m.Set(Model, s)
	}
	return m
}

func fmtTime(t time.Time, islocal bool) string {
	x := Time{
		Time:   t,
		Prec:   6, // seconds
		HasLoc: !islocal,
	}
	return x.String()
}
