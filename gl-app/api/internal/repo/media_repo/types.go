package media_repo

import (
	"time"

	"gorm.io/gorm"
)

type MediaAsset struct {
	gorm.Model
	UserID          int64  `gorm:"column:user_id"`
	ResourceType    string `gorm:"column:resource_type"`
	Bucket          string `gorm:"column:bucket"`
	ObjectKey       string `gorm:"column:object_key"`
	OriginName      string `gorm:"column:origin_name"`
	ContentType     string `gorm:"column:content_type"`
	Ext             string `gorm:"column:ext"`
	FileSize        int64  `gorm:"column:file_size"`
	Status          string `gorm:"column:status"`
	FormalBucket    string `gorm:"column:formal_bucket"`
	FormalObjectKey string `gorm:"column:formal_object_key"`
	DurationMs      int64  `gorm:"column:duration_ms"`
	Width           int64  `gorm:"column:width"`
	Height          int64  `gorm:"column:height"`
	BoundWorkID     int64  `gorm:"column:bound_work_id"`
	ProcessError    string `gorm:"column:process_error"`
}

func (MediaAsset) TableName() string { return "media_asset" }

type ProcessTask struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	WorkID       int64     `gorm:"column:work_id"`
	TaskType     string    `gorm:"column:task_type"`
	Status       string    `gorm:"column:status"`
	Stage        string    `gorm:"column:stage"`
	Progress     int64     `gorm:"column:progress"`
	RetryCount   int64     `gorm:"column:retry_count"`
	ErrorMessage string    `gorm:"column:error_message"`
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
