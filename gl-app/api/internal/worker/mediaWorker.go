package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	corequeue "github.com/zeromicro/go-zero/core/queue"
	"github.com/zeromicro/go-zero/core/service"
	"gl-app/api/internal/queue"
	"gl-app/api/internal/svc"
)

type MediaWorker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	svcCtx     *svc.ServiceContext
	kafkaQueue corequeue.MessageQueue
	stopOnce   sync.Once
}

var _ service.Service = (*MediaWorker)(nil)

type processAsset struct {
	ID              int64  `db:"id"`
	Role            string `db:"role"`
	Bucket          string `db:"bucket"`
	ObjectKey       string `db:"object_key"`
	OriginName      string `db:"origin_name"`
	ContentType     string `db:"content_type"`
	FormalBucket    string `db:"formal_bucket"`
	FormalObjectKey string `db:"formal_object_key"`
	Status          string `db:"status"`
}

type probeOutput struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int64  `json:"width"`
		Height    int64  `json:"height"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func NewMediaWorker(ctx context.Context, svcCtx *svc.ServiceContext) *MediaWorker {
	workerCtx, cancel := context.WithCancel(ctx)
	mediaWorker := &MediaWorker{
		ctx:    workerCtx,
		cancel: cancel,
		svcCtx: svcCtx,
	}
	mediaWorker.kafkaQueue = kq.MustNewQueue(svcCtx.Config.MediaQueue.KqConf, mediaWorker)
	return mediaWorker
}

// Start 实现 service.Service，由 servicegroup 与 HTTP 服务一同启动。
func (w *MediaWorker) Start() {
	if err := w.ensureFormalBucket(); err != nil {
		logx.Must(err)
	}
	if err := os.MkdirAll(w.svcCtx.Config.MediaWorker.TempDir, 0o750); err != nil {
		logx.Must(fmt.Errorf("创建媒体临时目录失败: %w", err))
	}
	// A worker can die after claiming a database task but before committing the
	// Kafka offset. Re-open only very old tasks to avoid competing with a
	// legitimately long transcode still running on another worker.
	_, _ = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_process_task SET status='pending',stage='queued'
		WHERE status='processing' AND updated_at < DATE_SUB(NOW(), INTERVAL 2 HOUR)`)
	go w.runDispatcher()
	go w.runScheduler()

	logx.Info("Kafka媒体消费者已启动")
	w.kafkaQueue.Start()
}

// Stop 实现 service.Service，确保系统退出时停止后台扫描和Kafka消费。
func (w *MediaWorker) Stop() {
	w.stopOnce.Do(func() {
		w.cancel()
		w.kafkaQueue.Stop()
	})
}

// Consume implements kq.ConsumeHandler.
func (w *MediaWorker) Consume(ctx context.Context, _ string, value string) error {
	message, err := queue.DecodeMessage(value)
	if err != nil {
		// A malformed message can never succeed after retry. Commit it and keep
		// the error visible in logs instead of blocking a Kafka partition.
		logx.WithContext(ctx).Errorf("忽略非法Kafka媒体消息: %v", err)
		return nil
	}
	err = w.processWork(message.WorkID)
	if err != nil {
		w.markFailed(message.WorkID, err)
		if message.Attempt < w.svcCtx.Config.MediaQueue.MaxRetry {
			if retryErr := w.svcCtx.MediaQueue.PublishProcessWork(ctx, message.WorkID, message.Attempt+1); retryErr != nil {
				// ForceCommit=false makes kq leave the current offset
				// uncommitted, so Kafka will deliver it again.
				return retryErr
			}
		}
	}
	return nil
}

func (w *MediaWorker) processWork(workID int64) error {
	result, err := w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_process_task SET status='processing',
		stage='preparing',progress=5,error_message='' WHERE work_id=? AND status IN ('pending','failed')`, workID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil
	}
	if _, err = w.svcCtx.SqlConn.ExecCtx(w.ctx,
		"UPDATE work SET process_status='processing',review_status='waiting_process' WHERE id=?", workID); err != nil {
		return err
	}

	var assets []processAsset
	if err = w.svcCtx.SqlConn.QueryRowsCtx(w.ctx, &assets, `SELECT ma.id,wa.role,ma.bucket,ma.object_key,
		ma.origin_name,ma.content_type,ma.formal_bucket,ma.formal_object_key,ma.status
		FROM work_asset wa JOIN media_asset ma ON ma.id=wa.media_asset_id
		WHERE wa.work_id=? ORDER BY CASE wa.role WHEN 'cover' THEN 0 ELSE 1 END,wa.sort`, workID); err != nil {
		return err
	}
	if len(assets) < 2 {
		return fmt.Errorf("作品素材不完整")
	}

	for index, asset := range assets {
		if asset.Status == "ready" && asset.FormalObjectKey != "" {
			continue
		}
		progress := int64(10 + (index * 75 / len(assets)))
		_, _ = w.svcCtx.SqlConn.ExecCtx(w.ctx,
			"UPDATE media_process_task SET stage=?,progress=? WHERE work_id=?",
			"processing_"+asset.Role, progress, workID)
		if asset.Role == "video" {
			if err = w.transcodeVideo(workID, asset); err != nil {
				return err
			}
		} else {
			if err = w.promoteImage(workID, asset); err != nil {
				return err
			}
		}
	}

	_, err = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_process_task SET status='succeeded',
		stage='waiting_review',progress=100,error_message='' WHERE work_id=?`, workID)
	if err != nil {
		return err
	}
	_, err = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE work SET process_status='succeeded',
		review_status='pending_review',publish_status='pending' WHERE id=?`, workID)
	return err
}

func (w *MediaWorker) promoteImage(workID int64, asset processAsset) error {
	ext := strings.ToLower(filepath.Ext(asset.OriginName))
	if ext == "" {
		ext = ".jpg"
	}
	target := fmt.Sprintf("works/%d/%s/%d%s", workID, asset.Role, asset.ID, ext)
	source := minio.CopySrcOptions{Bucket: asset.Bucket, Object: asset.ObjectKey}
	destination := minio.CopyDestOptions{Bucket: w.svcCtx.Config.Minio.FormalBucket, Object: target}
	if _, err := w.svcCtx.MinioClient.CopyObject(w.ctx, destination, source); err != nil {
		return fmt.Errorf("迁移图片%d失败: %w", asset.ID, err)
	}
	if _, err := w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_asset SET status='ready',formal_bucket=?,
		formal_object_key=?,process_error='' WHERE id=?`,
		w.svcCtx.Config.Minio.FormalBucket, target, asset.ID); err != nil {
		return err
	}
	_ = w.svcCtx.MinioClient.RemoveObject(w.ctx, asset.Bucket, asset.ObjectKey, minio.RemoveObjectOptions{})
	return nil
}

func (w *MediaWorker) transcodeVideo(workID int64, asset processAsset) error {
	baseDir, err := os.MkdirTemp(w.svcCtx.Config.MediaWorker.TempDir, fmt.Sprintf("work-%d-", workID))
	if err != nil {
		return fmt.Errorf("创建转码目录失败: %w", err)
	}
	defer os.RemoveAll(baseDir)

	input := filepath.Join(baseDir, "input"+filepath.Ext(asset.OriginName))
	outputDir := filepath.Join(baseDir, "hls")
	if err = os.MkdirAll(outputDir, 0o750); err != nil {
		return err
	}
	if err = w.svcCtx.MinioClient.FGetObject(w.ctx, asset.Bucket, asset.ObjectKey, input, minio.GetObjectOptions{}); err != nil {
		return fmt.Errorf("下载待转码视频失败: %w", err)
	}
	probe, err := w.probeVideo(input)
	if err != nil {
		return err
	}
	playlist := filepath.Join(outputDir, "index.m3u8")
	segmentPattern := filepath.Join(outputDir, "segment_%05d.ts")
	hlsTime := w.svcCtx.Config.MediaWorker.HlsTime
	if hlsTime <= 0 {
		hlsTime = 6
	}
	command := exec.CommandContext(w.ctx, w.svcCtx.Config.MediaWorker.FFmpegPath,
		"-y", "-i", input, "-c:v", "libx264", "-preset", "medium", "-crf", "23",
		"-c:a", "aac", "-b:a", "128k", "-movflags", "+faststart",
		"-hls_time", strconv.Itoa(hlsTime), "-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern, playlist)
	if output, commandErr := command.CombinedOutput(); commandErr != nil {
		return fmt.Errorf("ffmpeg转码失败: %w: %s", commandErr, tail(string(output), 1500))
	}

	prefix := fmt.Sprintf("works/%d/video/hls", workID)
	err = filepath.WalkDir(outputDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		relative, _ := filepath.Rel(outputDir, path)
		objectKey := prefix + "/" + filepath.ToSlash(relative)
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if strings.HasSuffix(path, ".m3u8") {
			contentType = "application/vnd.apple.mpegurl"
		} else if strings.HasSuffix(path, ".ts") {
			contentType = "video/mp2t"
		}
		info, statErr := os.Stat(path)
		if statErr != nil {
			return statErr
		}
		file, openErr := os.Open(path)
		if openErr != nil {
			return openErr
		}
		_, uploadErr := w.svcCtx.MinioClient.PutObject(w.ctx, w.svcCtx.Config.Minio.FormalBucket,
			objectKey, file, info.Size(), minio.PutObjectOptions{ContentType: contentType})
		closeErr := file.Close()
		if uploadErr == nil {
			uploadErr = closeErr
		}
		return uploadErr
	})
	if err != nil {
		return fmt.Errorf("上传HLS文件失败: %w", err)
	}

	manifestKey := prefix + "/index.m3u8"
	if _, err = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_asset SET status='ready',formal_bucket=?,
		formal_object_key=?,duration_ms=?,width=?,height=?,process_error='' WHERE id=?`,
		w.svcCtx.Config.Minio.FormalBucket, manifestKey, probe.durationMs, probe.width, probe.height, asset.ID); err != nil {
		return err
	}
	_ = w.svcCtx.MinioClient.RemoveObject(w.ctx, asset.Bucket, asset.ObjectKey, minio.RemoveObjectOptions{})
	return nil
}

type videoMetadata struct {
	durationMs int64
	width      int64
	height     int64
}

func (w *MediaWorker) probeVideo(input string) (videoMetadata, error) {
	command := exec.CommandContext(w.ctx, w.svcCtx.Config.MediaWorker.FFprobePath,
		"-v", "error", "-show_streams", "-show_format", "-of", "json", input)
	output, err := command.Output()
	if err != nil {
		return videoMetadata{}, fmt.Errorf("ffprobe探测失败: %w", err)
	}
	var probe probeOutput
	if err = json.Unmarshal(output, &probe); err != nil {
		return videoMetadata{}, err
	}
	duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
	metadata := videoMetadata{durationMs: int64(duration * 1000)}
	for _, stream := range probe.Streams {
		if stream.CodecType == "video" {
			metadata.width, metadata.height = stream.Width, stream.Height
			break
		}
	}
	if metadata.width == 0 || metadata.height == 0 {
		return videoMetadata{}, fmt.Errorf("文件中没有有效视频轨")
	}
	return metadata, nil
}

func (w *MediaWorker) markFailed(workID int64, processErr error) {
	message := tail(processErr.Error(), 1800)
	_, _ = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_process_task SET status='failed',stage='failed',
		error_message=? WHERE work_id=?`, message, workID)
	_, _ = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE work SET process_status='failed',
		review_status='waiting_process',publish_status='pending' WHERE id=?`, workID)
	_, _ = w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE media_asset ma JOIN work_asset wa ON wa.media_asset_id=ma.id
		SET ma.status=IF(ma.status='ready','ready','failed'),ma.process_error=? WHERE wa.work_id=?`, message, workID)
}

func (w *MediaWorker) runDispatcher() {
	type pendingWork struct {
		WorkID int64  `db:"work_id"`
		Type   string `db:"type"`
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			var works []pendingWork
			if err := w.svcCtx.SqlConn.QueryRowsCtx(w.ctx, &works, `SELECT t.work_id,w.type
				FROM media_process_task t
				JOIN work w ON w.id=t.work_id
				WHERE t.status='pending' AND t.updated_at < DATE_SUB(NOW(), INTERVAL 20 SECOND)
				LIMIT 100`); err != nil {
				continue
			}
			for _, item := range works {
				if item.Type == "image" {
					// 兼容升级前已经进入waiting_process的图片作品：
					// 直接执行MinIO迁移，不再补发Kafka消息。
					if err := w.processWork(item.WorkID); err != nil {
						w.markFailed(item.WorkID, err)
					}
					continue
				}
				_ = w.svcCtx.MediaQueue.PublishProcessWork(w.ctx, item.WorkID, 0)
			}
		}
	}
}

func (w *MediaWorker) runScheduler() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			_, err := w.svcCtx.SqlConn.ExecCtx(w.ctx, `UPDATE work SET publish_status='published',
				published_at=NOW(3) WHERE review_status='approved' AND publish_status='scheduled'
				AND scheduled_at<=NOW(3) AND deleted_at IS NULL`)
			if err != nil {
				logx.Errorf("执行定时发布扫描失败: %v", err)
			}
		}
	}
}

func (w *MediaWorker) ensureFormalBucket() error {
	bucket := w.svcCtx.Config.Minio.FormalBucket
	if bucket == "" {
		return errors.New("未配置MinIO正式资源桶")
	}
	exists, err := w.svcCtx.MinioClient.BucketExists(w.ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return w.svcCtx.MinioClient.MakeBucket(w.ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}

func tail(value string, size int) string {
	if len(value) <= size {
		return value
	}
	return value[len(value)-size:]
}
