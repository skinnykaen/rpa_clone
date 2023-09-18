package models

type MediaCore struct {
	URI string `json:"uri"`
}

func (ht *MediaHTTP) ToCore() *MediaCore {
	return &MediaCore{
		URI: ht.URI,
	}
}

func (ht *MediaHTTP) FromCore(media *MediaCore) {
	ht.URI = media.URI
}
