package media

import "testing"

func TestNormalizeCreatorWorkPage(t *testing.T) {
	tests := []struct {
		name         string
		page         int64
		pageSize     int64
		wantPage     int64
		wantPageSize int64
	}{
		{name: "默认分页", page: 0, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "正常分页", page: 2, pageSize: 10, wantPage: 2, wantPageSize: 10},
		{name: "限制最大条数", page: 1, pageSize: 100, wantPage: 1, wantPageSize: 50},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			page, pageSize := normalizeCreatorWorkPage(test.page, test.pageSize)
			if page != test.wantPage || pageSize != test.wantPageSize {
				t.Fatalf("got (%d,%d), want (%d,%d)", page, pageSize, test.wantPage, test.wantPageSize)
			}
		})
	}
}
