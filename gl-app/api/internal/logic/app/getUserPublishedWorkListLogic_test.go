package app

import "testing"

func TestNormalizeUserPublishedWorkPage(t *testing.T) {
	tests := []struct {
		name         string
		page         int64
		pageSize     int64
		wantPage     int64
		wantPageSize int64
	}{
		{name: "使用默认分页", page: 0, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "保留合法分页", page: 2, pageSize: 12, wantPage: 2, wantPageSize: 12},
		{name: "限制单页最大数量", page: 1, pageSize: 100, wantPage: 1, wantPageSize: 50},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			page, pageSize := normalizeRecommendPage(test.page, test.pageSize)
			if page != test.wantPage || pageSize != test.wantPageSize {
				t.Fatalf("分页结果为(%d,%d)，期望(%d,%d)", page, pageSize, test.wantPage, test.wantPageSize)
			}
		})
	}
}
