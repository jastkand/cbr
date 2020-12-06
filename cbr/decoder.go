package cbr

import (
	"bytes"
	"encoding/xml"
	"io"

	"golang.org/x/text/encoding/charmap"
)

func xmlDecoder(body []byte) *xml.Decoder {
	d := xml.NewDecoder(bytes.NewReader(body))
	d.CharsetReader = identReader
	return d
}

func identReader(encoding string, input io.Reader) (io.Reader, error) {
	switch encoding {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return input, nil
	}
}
