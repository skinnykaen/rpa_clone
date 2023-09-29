package models

type CourseApiMediaCollectionCore struct {
	BannerImage *AbsoluteMediaCore `json:"banner_image"`
	CourseImage *MediaCore         `json:"course_image"`
	CourseVideo *MediaCore         `json:"course_video"`
	Image       *ImageCore         `json:"image"`
}

func (ht *CourseAPIMediaCollectionHTTP) ToCore() *CourseApiMediaCollectionCore {
	bannerImageCore := &AbsoluteMediaCore{}
	bannerImageCore = ht.BannerImage.ToCore()
	courseImageCore := &MediaCore{}
	courseImageCore = ht.CourseImage.ToCore()
	courseVideoCore := &MediaCore{}
	courseVideoCore = ht.CourseVideo.ToCore()
	imageCore := &ImageCore{}
	imageCore = ht.Image.ToCore()
	return &CourseApiMediaCollectionCore{
		BannerImage: bannerImageCore,
		CourseImage: courseImageCore,
		CourseVideo: courseVideoCore,
		Image:       imageCore,
	}
}

func (ht *CourseAPIMediaCollectionHTTP) FromCore(courseApiMediaCollection *CourseApiMediaCollectionCore) {
	ht.BannerImage.FromCore(courseApiMediaCollection.BannerImage)
	ht.CourseImage.FromCore(courseApiMediaCollection.CourseImage)
	ht.CourseVideo.FromCore(courseApiMediaCollection.CourseVideo)
	ht.Image.FromCore(courseApiMediaCollection.Image)
}
