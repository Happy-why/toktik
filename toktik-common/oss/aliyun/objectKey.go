package aliyun

import (
	"strings"
)

func (o *OSS) CreateObjectKey(suffix string, directories ...string) string {
	var builder strings.Builder
	for i, v := range directories {
		builder.WriteString(v)
		if i != len(directories)-1 { //TODO 这里绝对可以优化，懒得测性能，下次看到一定优化
			builder.WriteString("/")
		}
	}
	builder.WriteString(suffix)
	return builder.String()
}

// bucket_name:why_bucket 视频目录：video/user_id/video_name 封面目录：cover/user_id/cover_name
// video_name：video_id + title.mp4  cover：video_id + title.jpg
