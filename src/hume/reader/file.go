package reader

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"hume/record"
	"io"
	"os"
	"unicode/utf8"
)

type FileReader struct {
	BaseReader
	InputFile  string   `json:"input_file",omitempty`
	Delim      string   `json:"delim"`
	HasHeader  bool     `json:"has_header"`
	Header     []string `json:"header"`
	IdField    string   `json:"id_field"`
	file       *os.File
	reader     *csv.Reader
	size       int64
	pos        int
	lineNumber int32
	log        *log.Entry
}

func (fr *FileReader) Init() error {
	var err error
	fr.lineNumber = 0

	fr.log = log.WithFields(log.Fields{
		"name": "FileReader",
		"file": fr.InputFile,
	})

	fr.log.Debug("Initializing file reader")

	fr.file, err = os.Open(fr.InputFile)
	if err != nil {
		fr.log.Error(err)
		return err
	}

	fs, err := fr.file.Stat()
	if err != nil {
		fr.log.Error(err)
		return err
	}
	fr.size = fs.Size()

	delimRune, size := utf8.DecodeRuneInString(fr.Delim)
	if size == 0 {
		err = fmt.Errorf("Invalid Delimiter: %s", fr.Delim)
		fr.log.Error(err)
		return err
	}

	fr.reader = csv.NewReader(fr.file)
	fr.reader.Comma = delimRune
	fr.reader.LazyQuotes = true
	fr.reader.FieldsPerRecord = -1

	fr.lineNumber += 1
	if fr.HasHeader {
		header, err := fr.reader.Read()
		if err != nil {
			if err == io.EOF {
				err = fmt.Errorf("Input file %s is empty", fr.InputFile)
			}
			fr.log.Error(err)
			return err
		}
		fr.Header = header
		fr.lineNumber += 1
	}

	fr.log.Debug("Done initalizing")

	return err
}

func (fr *FileReader) ChannelIterator(bufferSize int) <-chan *record.Record {
	output := make(chan *record.Record, bufferSize)
	go func() {
		var err error
		for err == nil {
			r := record.Record{}
			r.Number = fr.lineNumber
			r.Fields = &fr.Header

			r.Values, err = fr.reader.Read()
			if err != nil {
				break
			}
			r.InitMap()

			r.Id = uuid.NewRandom().String()
			if fr.IdField != "" {
				recId, ok := r.Map[fr.IdField]
				if ok {
					r.Id = recId
				} else {
					fr.log.WithFields(log.Fields{
						"file_line": fr.lineNumber,
					}).Errorf("Missing record id field %s", fr.IdField)
				}
			}

			output <- &r
			fr.lineNumber += 1
		}

		if err == io.EOF {
			fr.log.WithFields(log.Fields{
				"file_line": fr.lineNumber,
			}).Debug("End of file reached")
		} else {
			fr.log.Error(err)
		}
		close(output)
	}()
	return output
}

func (fr *FileReader) Close() error {
	return fr.file.Close()
}
