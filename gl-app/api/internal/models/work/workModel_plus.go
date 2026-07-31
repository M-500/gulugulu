package work

import "context"

func (*defaultWorkModel) PageList(ctx context.Context, workIDs []int, userID int, visibility string, title string, page, size int) ([]*Work, error) {
	var res []*Work

	return res, nil
}
