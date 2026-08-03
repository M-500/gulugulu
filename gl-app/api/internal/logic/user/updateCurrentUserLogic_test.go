package user

import (
	"testing"
	"time"

	"gl-app/api/internal/types"
)

func TestBuildProfileUpdateIncludesZeroValues(t *testing.T) {
	req := &types.UpdateUserProfileReq{
		NickName: "  测试用户  ",
		Bio:      "",
		Sex:      0,
		BothDay:  "",
	}

	update, err := buildProfileUpdate(req, time.Date(2026, 8, 3, 12, 0, 0, 0, time.Local))
	if err != nil {
		t.Fatalf("buildProfileUpdate() error = %v", err)
	}
	if update["nickname"] != "测试用户" || update["bio"] != "" || update["sex"] != int64(0) {
		t.Fatalf("零值字段没有完整保留: %#v", update)
	}
	if update["both_day"] != (*time.Time)(nil) {
		t.Fatalf("清空生日时 both_day = %#v, want nil", update["both_day"])
	}
}

func TestBuildProfileUpdateParsesBirthday(t *testing.T) {
	update, err := buildProfileUpdate(&types.UpdateUserProfileReq{
		NickName: "测试用户",
		Bio:      "保持热爱",
		Sex:      2,
		BothDay:  "2000-02-29",
	}, time.Date(2026, 8, 3, 12, 0, 0, 0, time.Local))
	if err != nil {
		t.Fatalf("buildProfileUpdate() error = %v", err)
	}
	birthday, ok := update["both_day"].(*time.Time)
	if !ok || birthday.Format(time.DateOnly) != "2000-02-29" {
		t.Fatalf("both_day = %#v, want 2000-02-29", update["both_day"])
	}
}

func TestBuildProfileUpdateRejectsInvalidFields(t *testing.T) {
	now := time.Date(2026, 8, 3, 12, 0, 0, 0, time.Local)
	tests := []struct {
		name string
		req  types.UpdateUserProfileReq
	}{
		{name: "invalid sex", req: types.UpdateUserProfileReq{NickName: "测试用户", Sex: 3}},
		{name: "invalid birthday", req: types.UpdateUserProfileReq{NickName: "测试用户", BothDay: "2026/01/01"}},
		{name: "future birthday", req: types.UpdateUserProfileReq{NickName: "测试用户", BothDay: "2026-08-04"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := buildProfileUpdate(&tt.req, now); err == nil {
				t.Fatal("buildProfileUpdate() error = nil, want validation error")
			}
		})
	}
}
