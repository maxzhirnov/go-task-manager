package analytics

import (
	"context"
	"log"

	"github.com/mixpanel/mixpanel-go"
)

type Tracker interface {
	Track(ctx context.Context, eventName string, distinctID string, properties map[string]any) error
	SetUserProfile(ctx context.Context, distinctID string, properties map[string]any) error
}

type Mixpanel struct {
	client *mixpanel.ApiClient
}

// Ensure Mixpanel implements Tracker interface
var _ Tracker = (*Mixpanel)(nil)

func NewMixpanel(token string) *Mixpanel {
	log.Printf("Initializing Mixpanel with token: %s", token)
	return &Mixpanel{
		client: mixpanel.NewApiClient(token),
	}
}

func (a *Mixpanel) Track(ctx context.Context, eventName string, distinctID string, properties map[string]any) error {
	log.Printf("Tracking event: %s, distinctID: %s, properties: %+v", eventName, distinctID, properties)
	event := a.client.NewEvent(eventName, distinctID, properties)
	err := a.client.Track(ctx, []*mixpanel.Event{event})
	if err != nil {
		log.Printf("Failed to track event: %s, distinctID: %s, properties: %+v, error: %v", eventName, distinctID, properties, err)
		return err
	}
	return nil
}

func (a *Mixpanel) SetUserProfile(ctx context.Context, distinctID string, properties map[string]any) error {
	log.Printf("Setting user profile: %s, properties: %+v", distinctID, properties)
	people := []*mixpanel.PeopleProperties{
		mixpanel.NewPeopleProperties(distinctID, properties),
	}
	err := a.client.PeopleSet(ctx, people)
	if err != nil {
		log.Printf("Failed to set user profile: %s, properties: %+v, error: %v", distinctID, properties, err)
		return err
	}
	return nil
}
