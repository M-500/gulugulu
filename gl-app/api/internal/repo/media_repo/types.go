package media_repo

import (
	"gorm.io/gorm"
)

type MediaAsset struct {
	gorm.Model
	UserID          int64  `gorm:"column:user_id;type:bigint;not null;comment:上传用户ID"`
	ResourceType    string `gorm:"column:resource_type;type:varchar(32);not null;comment:资源类型"`
	Bucket          string `gorm:"column:bucket;type:varchar(128);not null;comment:存储桶"`
	ObjectKey       string `gorm:"column:object_key;type:varchar(512);comment:对象存储桶路径"`
	OriginName      string `gorm:"column:origin_name;type:varchar(256);comment:原始文件名"`
	ContentType     string `gorm:"column:content_type;type:varchar(128);comment:内容类型"`
	Ext             string `gorm:"column:ext;type:varchar(32);comment:文件扩展名"`
	FileSize        int64  `gorm:"column:file_size;type:bigint;comment:文件大小"`
	Status          string `gorm:"column:status;type:varchar(32);comment:状态"`
	FormalBucket    string `gorm:"column:formal_bucket;type:varchar(128);comment:正式存储桶"`
	FormalObjectKey string `gorm:"column:formal_object_key;type:varchar(512);comment:正式对象存储桶路径"`
	DurationMs      int64  `gorm:"column:duration_ms;type:bigint;comment:持续时间(毫秒)"`
	Width           int64  `gorm:"column:width;type:bigint;comment:宽度"`
	Height          int64  `gorm:"column:height;type:bigint;comment:高度"`
	BoundWorkID     int64  `gorm:"column:bound_work_id;type:bigint;comment:绑定作品ID"`
	BoundCommentID  int64  `gorm:"column:bound_comment_id;type:bigint;not null;default:0;index;comment:绑定评论ID"`
	ProcessError    string `gorm:"column:process_error;type:varchar(1024);comment:处理错误信息"`
}

func (MediaAsset) TableName() string { return "media_asset" }

type ProcessTask struct {
	gorm.Model
	WorkID       int64  `gorm:"column:work_id;type:bigint;not null;comment:作品ID"`
	TaskType     string `gorm:"column:task_type;type:varchar(32);not null;comment:任务类型"`
	Status       string `gorm:"column:status;type:varchar(32);not null;comment:状态"`
	Stage        string `gorm:"column:stage;type:varchar(32);comment:阶段"`
	Progress     int64  `gorm:"column:progress;type:bigint;comment:进度"`
	RetryCount   int64  `gorm:"column:retry_count;type:bigint;comment:重试次数"`
	ErrorMessage string `gorm:"column:error_message;type:varchar(1024);comment:错误信息"`
}

func (ProcessTask) TableName() string { return "media_process_task" }

type workState struct {
	ID            int64          `gorm:"column:id;primaryKey"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
	UserID        int64          `gorm:"column:user_id"`
	Type          string         `gorm:"column:type"`
	ProcessStatus string         `gorm:"column:process_status"`
	ReviewStatus  string         `gorm:"column:review_status"`
	PublishStatus string         `gorm:"column:publish_status"`
}

func (workState) TableName() string { return "work" }

type ProcessAsset struct {
	ID              int64  `gorm:"column:id"`
	Role            string `gorm:"column:role"`
	Bucket          string `gorm:"column:bucket"`
	ObjectKey       string `gorm:"column:object_key"`
	OriginName      string `gorm:"column:origin_name"`
	ContentType     string `gorm:"column:content_type"`
	FormalBucket    string `gorm:"column:formal_bucket"`
	FormalObjectKey string `gorm:"column:formal_object_key"`
	Status          string `gorm:"column:status"`
}

type AssetPromotion struct {
	AssetID         int64
	FormalBucket    string
	FormalObjectKey string
}

type PendingWork struct {
	WorkID int64  `gorm:"column:work_id"`
	Type   string `gorm:"column:type"`
}
