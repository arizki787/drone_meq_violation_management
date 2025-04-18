package config

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func Credentials() (cld *cloudinary.Cloudinary, ctx context.Context) {
	var err error
	cld, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	cld.Config.URL.Secure = true
	ctx = context.Background()
	return cld, ctx
}
