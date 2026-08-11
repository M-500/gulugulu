package media

import "testing"

func TestMpegTSUploadIsAcceptedAndNormalized(t *testing.T) {
	_, ext, err := normalizeFileName("camera-recording.TS")
	if err != nil {
		t.Fatalf("TS文件名应合法: %v", err)
	}
	if ext != ".ts" {
		t.Fatalf("扩展名未规范化: %q", ext)
	}
	if !isAllowedExt("video", ext) {
		t.Fatal("TS应在视频上传白名单中")
	}

	for _, inputType := range []string{"", "application/octet-stream", "video/mp2t"} {
		if got := normalizeUploadContentType("video", ext, inputType); got != "video/mp2t" {
			t.Fatalf("TS Content-Type规范化错误: input=%q got=%q", inputType, got)
		}
	}
}
