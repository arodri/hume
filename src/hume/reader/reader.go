package reader

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"hume/record"
)

func GetReader(config []byte) Reader {
	var err error

	br := BaseReader{}
	err = json.Unmarshal(config, &br)
	if err != nil {
		logrus.Error(err)
		logrus.Fatal("Missing reader name and/or type")
	}

	var r Reader

	switch br.Type {
	default:
		logrus.Panicf("Unrecognized reader type: %#v", br)
	case "FileReader":
		fr := FileReader{}
		err = json.Unmarshal(config, &fr)
		r = &fr
	}

	if err != nil {
		logrus.Error(err)
		logrus.Fatal(fmt.Sprintf("Invalid configuration for %#v", br))
	}

	return r
}

type BaseReader struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (br BaseReader) GetName() string {
	return br.Name
}

func (br BaseReader) GetType() string {
	return br.Type
}

type Reader interface {
	GetName() string
	GetType() string
	Init() error
	ChannelIterator(bufferSize int) <-chan *record.Record
	Close() error
}
