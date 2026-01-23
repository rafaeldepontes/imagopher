package processor

type Repository interface {
	FindImageByID(id uint64) (any, error)
	FindImages() (any, error)
	TransformImage(img any) error
	UploadImage(img any) (uint64, error)
}
