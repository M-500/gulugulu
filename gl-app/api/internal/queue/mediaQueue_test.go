package queue

import "testing"

func TestDecodeMessage(t *testing.T) {
	message, err := DecodeMessage(`{"type":"process_work","workId":1001,"attempt":2}`)
	if err != nil {
		t.Fatalf("valid message rejected: %v", err)
	}
	if message.WorkID != 1001 || message.Attempt != 2 {
		t.Fatalf("unexpected message: %#v", message)
	}
}

func TestDecodeMessageRejectsInvalidPayload(t *testing.T) {
	if _, err := DecodeMessage(`{"type":"unknown","workId":0,"attempt":-1}`); err == nil {
		t.Fatal("invalid message should be rejected")
	}
}
