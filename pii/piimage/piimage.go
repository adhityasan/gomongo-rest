package piimage

import (
	"io/ioutil"
	"net/http"
	"net/textproto"
)

// ImageStruct standart Struct for image
type ImageStruct struct {
	Data   []byte               `schema:"data,omitempty" bson:"data,omitempty"`
	Name   string               `schema:"name,omitempty" bson:"name,omitempty"`
	Size   int64                `schema:"size,omitempty" bson:"size,omitempty"`
	Header textproto.MIMEHeader `schema:"header,omitempty" bson:"header,omitempty"`
}

// ImageStructHandler create image Struct with Data , Name, Size, Header from a request
func ImageStructHandler(fieldname string, r *http.Request) (*ImageStruct, error) {
	file, handler, err := r.FormFile(fieldname)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	newImageStruct := new(ImageStruct)
	newImageStruct.Data = fileBytes
	newImageStruct.Name = handler.Filename
	newImageStruct.Size = handler.Size
	newImageStruct.Header = handler.Header

	return newImageStruct, nil
}
