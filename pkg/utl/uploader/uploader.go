package uploader

type (
	Uploader interface {
		GetPresignUploadUrl(string) (string, error)
		GetPresignFetchUrl(string) (string, error)
	}
)
