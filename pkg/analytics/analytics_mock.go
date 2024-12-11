package analytics

import (
	"context"
	"log"
	"time"
)

type Mock struct {
	isDebug bool
}

func NewMock(token string, debug bool) *Mock {
	return &Mock{
		isDebug: debug,
	}
}

var _ Tracker = (*Mock)(nil)

func (a *Mock) Track(ctx context.Context, eventName string, distinctID string, properties map[string]any) error {
	if properties == nil {
		properties = make(map[string]any)
	}
	properties["timestamp"] = time.Now()

	if a.isDebug {
		log.Printf("EVENT: %s, USER: %s, PROPERTIES: %+v\n", eventName, distinctID, properties)
	}
	return nil
}

func (a *Mock) SetUserProfile(ctx context.Context, distinctID string, properties map[string]any) error {
	if a.isDebug {
		log.Printf("SET USER PROFILE: %s, PROPERTIES: %+v\n", distinctID, properties)
	}
	return nil
}
