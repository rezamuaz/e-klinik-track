package imagekio

import (
	"context"
	"fmt"
	"time"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func UploadSingle(c context.Context, url string, folder string, filename string, extention string) (*uploader.UploadResponse, error) {
	var err error
	ik := imagekit.NewFromParams(
		imagekit.NewParams{
			PublicKey:   "public_kQrMJTZ4GBB33wOvt0iOGDkeZ84=",
			PrivateKey:  "private_veIvwkmiZTJkhsqCtbrPGnA1IzM=",
			UrlEndpoint: "https://ik.imagekit.io/04emdmsez"})
	var b1 *bool
	b1 = new(bool) // b1 now points to a bool with a zero value (false)
	*b1 = false
	result, err := ik.Uploader.Upload(c, url, uploader.UploadParam{
		Folder:   folder,
		FileName: fmt.Sprintf("%s.%s", filename, extention), UseUniqueFileName: b1},
	)
	if err != nil {
		return &uploader.UploadResponse{}, err
	}
	time.Sleep(2 * time.Second)
	return result, nil
}
