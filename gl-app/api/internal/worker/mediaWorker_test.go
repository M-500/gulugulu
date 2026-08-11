package worker

import (
	"slices"
	"testing"
)

func TestBuildTranscodeArgsForMpegTS(t *testing.T) {
	args := buildTranscodeArgs("/tmp/input.ts", "/tmp/segment_%05d.ts", "/tmp/index.m3u8", 6)

	assertArgPair(t, args, "-fflags", "+genpts")
	assertArgPair(t, args, "-map", "0:v:0")
	assertArgPair(t, args, "-avoid_negative_ts", "make_zero")
	if !slices.Contains(args, "0:a:0?") {
		t.Fatal("音频轨应为可选，纯视频TS也必须能够转码")
	}
}

func TestBuildTranscodeArgsDoesNotForceGeneratedTimestampsForMP4(t *testing.T) {
	args := buildTranscodeArgs("/tmp/input.mp4", "/tmp/segment_%05d.ts", "/tmp/index.m3u8", 6)
	if slices.Contains(args, "+genpts") {
		t.Fatal("普通MP4不应强制生成时间戳")
	}
}

func assertArgPair(t *testing.T, args []string, key, value string) {
	t.Helper()
	for index := 0; index+1 < len(args); index++ {
		if args[index] == key && args[index+1] == value {
			return
		}
	}
	t.Fatalf("转码参数缺少 %s %s: %#v", key, value, args)
}
