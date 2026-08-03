package user_repo

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type fakeUserDAO struct {
	user  *User
	err   error
	steps *[]string
}

func (f *fakeUserDAO) FindOneByEmail(context.Context, string) (*User, error) { return f.user, f.err }
func (f *fakeUserDAO) FindOneByID(context.Context, int64) (*User, error) {
	if f.steps != nil {
		*f.steps = append(*f.steps, "dao.find")
	}
	return f.user, f.err
}
func (f *fakeUserDAO) Create(context.Context, *User) error { return f.err }
func (f *fakeUserDAO) Update(context.Context, *User) error { return f.err }
func (f *fakeUserDAO) UpdateByMap(context.Context, int64, map[string]any) error {
	if f.steps != nil {
		*f.steps = append(*f.steps, "dao.update")
	}
	return f.err
}
func (f *fakeUserDAO) Delete(context.Context, int64) error { return f.err }

type fakeUserCache struct {
	user  *User
	err   error
	steps *[]string
}

func (f *fakeUserCache) FindByID(context.Context, int64) (*User, error) {
	if f.steps != nil {
		*f.steps = append(*f.steps, "cache.find")
	}
	return f.user, f.err
}
func (f *fakeUserCache) SetByID(context.Context, int64, *User, time.Duration) error {
	if f.steps != nil {
		*f.steps = append(*f.steps, "cache.set")
	}
	return f.err
}
func (f *fakeUserCache) DeleteByID(context.Context, int64) error {
	if f.steps != nil {
		*f.steps = append(*f.steps, "cache.delete")
	}
	return f.err
}

func TestFindOneByIDUsesCacheFirst(t *testing.T) {
	steps := make([]string, 0)
	want := &User{ID: 7, Email: "cache@example.com"}
	repo := NewUserRepo(
		&fakeUserDAO{user: &User{ID: 99}, steps: &steps},
		&fakeUserCache{user: want, steps: &steps},
	)
	got, err := repo.FindOneByID(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("FindOneByID() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindOneByID() = %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(steps, []string{"cache.find"}) {
		t.Fatalf("调用顺序 = %v, want [cache.find]", steps)
	}
}

func TestFindOneByIDFallsBackToDAOAndBackfillsCache(t *testing.T) {
	steps := make([]string, 0)
	want := &User{ID: 8, Email: "db@example.com"}
	repo := NewUserRepo(
		&fakeUserDAO{user: want, steps: &steps},
		&fakeUserCache{err: redis.Nil, steps: &steps},
	)
	got, err := repo.FindOneByID(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("FindOneByID() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindOneByID() = %#v, want %#v", got, want)
	}
	wantSteps := []string{"cache.find", "dao.find", "cache.set"}
	if !reflect.DeepEqual(steps, wantSteps) {
		t.Fatalf("调用顺序 = %v, want %v", steps, wantSteps)
	}
}

func TestUpdateByMapInvalidatesCacheAfterDatabaseSuccess(t *testing.T) {
	steps := make([]string, 0)
	repo := NewUserRepo(&fakeUserDAO{steps: &steps}, &fakeUserCache{steps: &steps})
	if err := repo.UpdateByMap(context.Background(), 9, map[string]any{"nickname": "新昵称"}); err != nil {
		t.Fatalf("UpdateByMap() error = %v", err)
	}
	want := []string{"dao.update", "cache.delete"}
	if !reflect.DeepEqual(steps, want) {
		t.Fatalf("调用顺序 = %v, want %v", steps, want)
	}
}

func TestUpdateByMapDoesNotDeleteCacheWhenDatabaseFails(t *testing.T) {
	dbErr := errors.New("database unavailable")
	steps := make([]string, 0)
	repo := NewUserRepo(&fakeUserDAO{err: dbErr, steps: &steps}, &fakeUserCache{steps: &steps})
	err := repo.UpdateByMap(context.Background(), 10, map[string]any{"nickname": "新昵称"})
	if !errors.Is(err, dbErr) {
		t.Fatalf("UpdateByMap() error = %v, want %v", err, dbErr)
	}
	if !reflect.DeepEqual(steps, []string{"dao.update"}) {
		t.Fatalf("调用顺序 = %v, want [dao.update]", steps)
	}
}
