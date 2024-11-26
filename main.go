package main

import (
	"encoding/xml"
	"io"
	"log/slog"
	"os"
)

func main() {
	f, _ := os.Open("./data/doc1.xml")
	dec := xml.NewDecoder(f)
	for {
		tok, err := dec.Token()
		if tok == nil {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("fail", "err", err.Error())
			break
		}
		switch elm := tok.(type) {
		case xml.StartElement:
			slog.Info("start", "elm", elm)
		case xml.EndElement:
			slog.Info("end", "elm", elm)
		case xml.ProcInst:
			slog.Info("inst", "elm", string(elm.Inst))
		case xml.CharData:
			slog.Info("data", "elm", string(elm))
		default:
			slog.Info("other", "elm", elm)
		}
	}
}
