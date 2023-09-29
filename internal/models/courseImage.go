package models

import (
	"gorm.io/gorm"
)

type ImageCore struct {
	Raw   string `json:"raw"`
	Small string `json:"small"`
	Large string `json:"large"`
}

type ImageDB struct {
	gorm.Model
	Raw   string
	Small string
	Large string
}

func (ht *ImageHTTP) ToCore() *ImageCore {
	return &ImageCore{
		Raw:   ht.Raw,
		Small: ht.Small,
		Large: ht.Large,
	}
}

func (ht *ImageHTTP) FromCore(imageCore *ImageCore) {
	ht.Raw = imageCore.Raw
	ht.Small = imageCore.Small
	ht.Large = imageCore.Large
}
