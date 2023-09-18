package models

type AbsoluteMediaCore struct {
	Uri         string `json:"uri"`
	UriAbsolute string `json:"uri_absolute"`
}

func (ht *AbsoluteMediaHTTP) ToCore() *AbsoluteMediaCore {
	return &AbsoluteMediaCore{
		Uri:         ht.URI,
		UriAbsolute: ht.URIAbsolute,
	}
}

func (ht *AbsoluteMediaHTTP) FromCore(absoluteMedia *AbsoluteMediaCore) {
	ht.URI = absoluteMedia.Uri
	ht.URIAbsolute = absoluteMedia.UriAbsolute
}
