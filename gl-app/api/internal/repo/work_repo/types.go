package work_repo

import (
	"time"

	"gorm.io/gorm"
)

type Work struct {
	ID                int64          `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at"`
	UserID            int64          `gorm:"column:user_id"`
	Type              string         `gorm:"column:type"`
	Title             string         `gorm:"column:title"`
	Content           string         `gorm:"column:content"`
	Visibility        string         `gorm:"column:visibility"`
	VisibilityUserIDs *string        `gorm:"column:visibility_user_ids;type:json"`
	CollectionID      int64          `gorm:"column:collection_id"`
	Original          bool           `gorm:"column:original"`
	CoverAssetID      int64          `gorm:"column:cover_asset_id"`
	ProcessStatus     string         `gorm:"column:process_status"`
	ReviewStatus      string         `gorm:"column:review_status"`
	PublishStatus     string         `gorm:"column:publish_status"`
	ScheduledAt       *time.Time     `gorm:"column:scheduled_at"`
	PublishedAt       *time.Time     `gorm:"column:published_at"`
	ReviewedBy        int64          `gorm:"column:reviewed_by"`
	ReviewedAt        *time.Time     `gorm:"column:reviewed_at"`
	ReviewReason      string         `gorm:"column:review_reason"`
	IdempotencyKey    string         `gorm:"column:idempotency_key"`
}

func (Work) TableName() string { return "work" }

type WorkAsset struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	WorkID       int64     `gorm:"column:work_id"`
	MediaAssetID int64     `gorm:"column:media_asset_id"`
	Role         string    `gorm:"column:role"`
	Sort         int64     `gorm:"column:sort"`
}

func (WorkAsset) TableName() string { return "work_asset" }

type Topic struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	Name           string    `gorm:"column:name"`
	NormalizedName string    `gorm:"column:normalized_name"`
	ViewNum        int64     `gorm:"column:view_num"`
	CommentNum     int64     `gorm:"column:comment_num"`
}

func (Topic) TableName() string { return "topic" }

type WorkTopic struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"column:created_at"`
	WorkID    int64     `gorm:"column:work_id"`
	TopicID   int64     `gorm:"column:topic_id"`
	Sort      int64     `gorm:"column:sort"`
}

func (WorkTopic) TableName() string { return "work_topic" }

type collection struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	UserID    int64          `gorm:"column:user_id"`
}

func (collection) TableName() string { return "collection" }

type mediaAsset struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
	UserID       int64          `gorm:"column:user_id"`
	ResourceType string         `gorm:"column:resource_type"`
	Status       string         `gorm:"column:status"`
	BoundWorkID  int64          `gorm:"column:bound_work_id"`
}

func (mediaAsset) TableName() string { return "media_asset" }

type mediaProcessTask struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	WorkID     int64  `gorm:"column:work_id"`
	TaskType   string `gorm:"column:task_type"`
	Status     string `gorm:"column:status"`
	Stage      string `gorm:"column:stage"`
	Progress   int64  `gorm:"column:progress"`
	RetryCount int64  `gorm:"column:retry_count"`
}

func (mediaProcessTask) TableName() string { return "media_process_task" }

type CreateAssetInput struct {
	MediaID int64
	Sort    int64
}

type CreateTopicInput struct {
	Name string
	Sort int64
}

type CreateWorkInput struct {
	UserID            int64
	Type              string
	Title             string
	Content           string
	Visibility        string
	VisibilityUserIDs string
	CollectionID      int64
	Original          bool
	CoverAssetID      int64
	ScheduledAt       *time.Time
	IdempotencyKey    string
	Assets            []CreateAssetInput
	Topics            []CreateTopicInput
}

// AssetView、TopicView 是查询投影，避免业务层感知表连接细节。
type AssetView struct {
	MediaID    int64  `gorm:"column:media_asset_id"`
	Role       string `gorm:"column:role"`
	Sort       int64  `gorm:"column:sort"`
	Bucket     string `gorm:"column:formal_bucket"`
	ObjectKey  string `gorm:"column:formal_object_key"`
	DurationMs int64  `gorm:"column:duration_ms"`
	Width      int64  `gorm:"column:width"`
	Height     int64  `gorm:"column:height"`
}

type TopicView struct {
	ID   int64  `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

type OwnerWorkDetail struct {
	ID            int64      `gorm:"column:id"`
	Type          string     `gorm:"column:type"`
	Title         string     `gorm:"column:title"`
	Content       string     `gorm:"column:content"`
	Visibility    string     `gorm:"column:visibility"`
	CollectionID  int64      `gorm:"column:collection_id"`
	Original      bool       `gorm:"column:original"`
	ProcessStatus string     `gorm:"column:process_status"`
	ReviewStatus  string     `gorm:"column:review_status"`
	PublishStatus string     `gorm:"column:publish_status"`
	ReviewReason  string     `gorm:"column:review_reason"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	ScheduledAt   *time.Time `gorm:"column:scheduled_at"`
	PublishedAt   *time.Time `gorm:"column:published_at"`
}

type WorkSummary struct {
	AllCount       int64 `gorm:"column:all_count"`
	PublishedCount int64 `gorm:"column:published_count"`
	ReviewingCount int64 `gorm:"column:reviewing_count"`
	RejectedCount  int64 `gorm:"column:rejected_count"`
}

type CreatorListQuery struct {
	UserID  int64
	Status  string
	Keyword string
	Limit   int64
	Offset  int64
}

type CreatorWorkItem struct {
	ID             int64      `gorm:"column:id"`
	Type           string     `gorm:"column:type"`
	Title          string     `gorm:"column:title"`
	DurationMs     int64      `gorm:"column:duration_ms"`
	Visibility     string     `gorm:"column:visibility"`
	ReviewStatus   string     `gorm:"column:review_status"`
	PublishStatus  string     `gorm:"column:publish_status"`
	ReviewReason   string     `gorm:"column:review_reason"`
	ScheduledAt    *time.Time `gorm:"column:scheduled_at"`
	PublishedAt    *time.Time `gorm:"column:published_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	CoverBucket    string     `gorm:"column:cover_bucket"`
	CoverObjectKey string     `gorm:"column:cover_object_key"`
}

type WorkStatus struct {
	ID            int64      `gorm:"column:id"`
	ProcessStatus string     `gorm:"column:process_status"`
	ReviewStatus  string     `gorm:"column:review_status"`
	PublishStatus string     `gorm:"column:publish_status"`
	ScheduledAt   *time.Time `gorm:"column:scheduled_at"`
	PublishedAt   *time.Time `gorm:"column:published_at"`
	Stage         string     `gorm:"column:stage"`
	Progress      int64      `gorm:"column:progress"`
	ErrorMessage  string     `gorm:"column:error_message"`
}

type AuditSummary struct {
	PendingCount  int64 `gorm:"column:pending_count"`
	ApprovedCount int64 `gorm:"column:approved_count"`
	RejectedCount int64 `gorm:"column:rejected_count"`
}

type AuditListQuery struct {
	Status   string
	WorkType string
	Keyword  string
	Limit    int64
	Offset   int64
}

type AuditWorkItem struct {
	ID             int64      `gorm:"column:id"`
	Type           string     `gorm:"column:type"`
	Title          string     `gorm:"column:title"`
	ContentExcerpt string     `gorm:"column:content_excerpt"`
	DurationMs     int64      `gorm:"column:duration_ms"`
	Visibility     string     `gorm:"column:visibility"`
	AuthorID       int64      `gorm:"column:user_id"`
	AuthorName     string     `gorm:"column:author_name"`
	ReviewStatus   string     `gorm:"column:review_status"`
	PublishStatus  string     `gorm:"column:publish_status"`
	ReviewReason   string     `gorm:"column:review_reason"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	ReviewedAt     *time.Time `gorm:"column:reviewed_at"`
	CoverBucket    string     `gorm:"column:cover_bucket"`
	CoverObjectKey string     `gorm:"column:cover_object_key"`
}

type AuditWorkDetail struct {
	ID            int64      `gorm:"column:id"`
	Type          string     `gorm:"column:type"`
	Title         string     `gorm:"column:title"`
	Content       string     `gorm:"column:content"`
	Visibility    string     `gorm:"column:visibility"`
	Original      bool       `gorm:"column:original"`
	AuthorID      int64      `gorm:"column:user_id"`
	AuthorName    string     `gorm:"column:author_name"`
	AuthorEmail   string     `gorm:"column:author_email"`
	ReviewStatus  string     `gorm:"column:review_status"`
	PublishStatus string     `gorm:"column:publish_status"`
	ReviewReason  string     `gorm:"column:review_reason"`
	ScheduledAt   *time.Time `gorm:"column:scheduled_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	ReviewedAt    *time.Time `gorm:"column:reviewed_at"`
	CoverBucket   string     `gorm:"column:cover_bucket"`
	CoverObject   string     `gorm:"column:cover_object_key"`
}

type AuditState struct {
	ProcessStatus string     `gorm:"column:process_status"`
	ReviewStatus  string     `gorm:"column:review_status"`
	ScheduledAt   *time.Time `gorm:"column:scheduled_at"`
}

type ReviewInput struct {
	WorkID        int64
	ReviewerID    int64
	ReviewStatus  string
	PublishStatus string
	Reason        string
	PublishedAt   *time.Time
}

type PublicWorkItem struct {
	ID             int64      `gorm:"column:id"`
	Type           string     `gorm:"column:type"`
	Title          string     `gorm:"column:title"`
	ContentExcerpt string     `gorm:"column:content_excerpt"`
	PublishedAt    *time.Time `gorm:"column:published_at"`
	CoverBucket    string     `gorm:"column:cover_bucket"`
	CoverObjectKey string     `gorm:"column:cover_object_key"`
	DurationMs     int64      `gorm:"column:duration_ms"`
	AuthorID       int64      `gorm:"column:author_id"`
	AuthorName     string     `gorm:"column:author_name"`
	AuthorAvatar   string     `gorm:"column:author_avatar"`
}

type PublicWorkDetail struct {
	ID           int64      `gorm:"column:id"`
	Type         string     `gorm:"column:type"`
	Title        string     `gorm:"column:title"`
	Content      string     `gorm:"column:content"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	AuthorID     int64      `gorm:"column:author_id"`
	AuthorName   string     `gorm:"column:author_name"`
	AuthorAvatar string     `gorm:"column:author_avatar"`
}

type VideoAsset struct {
	Bucket    string `gorm:"column:formal_bucket"`
	ObjectKey string `gorm:"column:formal_object_key"`
}

type VideoAccessScope int

const (
	VideoAccessOwner VideoAccessScope = iota
	VideoAccessAudit
	VideoAccessPublic
)
