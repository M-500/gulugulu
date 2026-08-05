package constants

type BizType string

const (
	WorkType    BizType = "NOTE_TYPE"
	CommentType BizType = "COMMENT_TYPE"
	AvatarType  BizType = "AVATAR_TYPE"

	// CommentTYpe 保留旧名称，避免已有调用方升级时中断。
	CommentTYpe = CommentType
)

func (b BizType) Valid() bool {
	return b == WorkType || b == CommentType || b == AvatarType
}
