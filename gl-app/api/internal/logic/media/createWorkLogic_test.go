package media

import (
	"testing"
	"time"
)

func TestValidateCreateWorkPayload(t *testing.T) {
	payload := createWorkPayload{
		IdempotencyKey: "publish-1",
		Type:           "video",
		Title:          "测试视频",
	}
	payload.Visibility.Type = "public"
	payload.Assets = append(payload.Assets, struct {
		MediaID int64 `json:"mediaId"`
		Sort    int   `json:"sort"`
	}{MediaID: 1})
	if err := validateCreateWorkPayload(&payload); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}

	payload.Assets = append(payload.Assets, struct {
		MediaID int64 `json:"mediaId"`
		Sort    int   `json:"sort"`
	}{MediaID: 2})
	if err := validateCreateWorkPayload(&payload); err == nil {
		t.Fatal("video payload with two assets should be rejected")
	}
}

func TestValidateScheduledAt(t *testing.T) {
	payload := createWorkPayload{
		IdempotencyKey: "publish-2",
		Type:           "image",
		Title:          "定时图片",
		ScheduledAt:    time.Now().Add(10 * time.Minute).Format(time.RFC3339),
	}
	payload.Visibility.Type = "private"
	payload.Assets = append(payload.Assets, struct {
		MediaID int64 `json:"mediaId"`
		Sort    int   `json:"sort"`
	}{MediaID: 1})
	if err := validateCreateWorkPayload(&payload); err != nil {
		t.Fatalf("valid scheduled payload rejected: %v", err)
	}
}

func TestValidateImageAssetLimit(t *testing.T) {
	newPayload := func(assetCount int) createWorkPayload {
		payload := createWorkPayload{IdempotencyKey: "image-limit", Type: "image", Title: "多图作品"}
		payload.Visibility.Type = "public"
		for index := 0; index < assetCount; index++ {
			payload.Assets = append(payload.Assets, struct {
				MediaID int64 `json:"mediaId"`
				Sort    int   `json:"sort"`
			}{MediaID: int64(index + 1), Sort: index})
		}
		return payload
	}

	allowed := newPayload(maxImageAssets)
	if err := validateCreateWorkPayload(&allowed); err != nil {
		t.Fatalf("%d张图片应允许发布: %v", maxImageAssets, err)
	}
	rejected := newPayload(maxImageAssets + 1)
	if err := validateCreateWorkPayload(&rejected); err == nil {
		t.Fatalf("%d张图片应被拒绝", maxImageAssets+1)
	}
}
