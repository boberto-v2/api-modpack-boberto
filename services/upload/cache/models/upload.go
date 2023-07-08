package upload_cache_models

import "time"

type UploadCache struct {
	Id        string    `json:"id"`
	OutputDir string    `json:"output_dir"`
	CreateAt  time.Time `json:"create_at"`
	ExpireAt  time.Time `json:"expire_at"`
}
