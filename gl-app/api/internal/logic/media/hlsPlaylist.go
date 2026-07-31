package media

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"gl-app/api/internal/svc"

	"github.com/minio/minio-go/v7"
)

const hlsPlaylistContentType = "application/vnd.apple.mpegurl"

// BuildSignedHLSPlaylist 读取HLS清单并将其中的分片地址替换为短期有效的签名地址。
// App公开播放、创作者预览和审核预览共用这一套清单生成规则。
func BuildSignedHLSPlaylist(ctx context.Context, svcCtx *svc.ServiceContext, bucket, objectKey string) (string, error) {
	if bucket == "" || objectKey == "" {
		return "", fmt.Errorf("视频播放清单不存在")
	}

	object, err := svcCtx.MinioClient.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("读取视频播放清单失败: %w", err)
	}
	defer object.Close()

	content, err := io.ReadAll(io.LimitReader(object, 2<<20))
	if err != nil {
		return "", fmt.Errorf("读取视频播放清单失败: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	basePath := path.Dir(objectKey)
	for index, line := range lines {
		segment := strings.TrimSpace(line)
		if segment == "" || strings.HasPrefix(segment, "#") {
			continue
		}

		presigned, signErr := svcCtx.MinioClient.PresignedGetObject(ctx, bucket, path.Join(basePath, segment), 30*time.Minute, nil)
		if signErr != nil {
			return "", fmt.Errorf("生成视频分片签名地址失败: %w", signErr)
		}
		lines[index] = presigned.String()
	}

	return strings.Join(lines, "\n"), nil
}
