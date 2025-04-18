/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AnalyticsStoreInterface interface {
	RecordURLAccess(ctx context.Context, shortID, ipAddress string) error
	GetURLAnalytics(ctx context.Context, shortID string) (*URLAnalytics, error)
}

// URLAnalytics stores analytics information for the URL
type URLAnalytics struct {
	TotalClicks   int64     `json:"total_clicks"`
	FirstAccessed time.Time `json:"first_accessed"`
	LastAccessed  time.Time `json:"last_accessed"`
	UniqueVisits  int64     `json:"unique_visits"`
}

// AnalyticsStore manages analytics on Redis
type AnalyticsStore struct {
	client *redis.Client
}

// New Analytics Store creates a new analytics store
func NewAnalyticsStore(client *redis.Client) *AnalyticsStore {
	return &AnalyticsStore{client: client}
}

// RecordURLAccess records a URL access
func (a *AnalyticsStore) RecordURLAccess(
	ctx context.Context,
	shortID,
	ipAddress string,
) error {
	// Analytics keys
	totalClicksKey := fmt.Sprintf("analytics:%s:total_clicks", shortID)
	uniqueVisitsKey := fmt.Sprintf("analytics:%s:unique_visits", shortID)
	lastAccessedKey := fmt.Sprintf("analytics:%s:last_accessed", shortID)
	firstAccessedKey := fmt.Sprintf("analytics:%s:first_accessed", shortID)
	uniqueIPKey := fmt.Sprintf("analytics:%s:unique_ips", shortID)

	// Pipeline
	pipe := a.client.Pipeline()

	// Increase total clicks
	pipe.Incr(ctx, totalClicksKey)

	// Record first access time (does not change if already exists)
	pipe.SetNX(ctx, firstAccessedKey, time.Now().Format(time.RFC3339), 0)

	// Update last access time
	pipe.Set(ctx, lastAccessedKey, time.Now().Format(time.RFC3339), 0)

	//Unique IP control
	uniqueVisit, err := pipe.SAdd(ctx, uniqueIPKey, ipAddress).Result()
	if err != nil {
		return err
	}

	// Increase unique visits
	if uniqueVisit > 0 {
		pipe.Incr(ctx, uniqueVisitsKey)
	}

	// Run pipeline
	_, err = pipe.Exec(ctx)
	return err
}

// GetURLAnalytics retrieves analytics for a URL
func (a *AnalyticsStore) GetURLAnalytics(
	ctx context.Context,
	shortID string,
) (*URLAnalytics, error) {
	// Analytics keys
	totalClicksKey := fmt.Sprintf("analytics:%s:total_clicks", shortID)
	uniqueVisitsKey := fmt.Sprintf("analytics:%s:unique_visits", shortID)
	lastAccessedKey := fmt.Sprintf("analytics:%s:last_accessed", shortID)
	firstAccessedKey := fmt.Sprintf("analytics:%s:first_accessed", shortID)
	uniqueIPKey := fmt.Sprintf("analytics:%s:unique_ips", shortID)

	//Collecting data with Pipeline
	pipe := a.client.Pipeline()
	totalClicksCmd := pipe.Get(ctx, totalClicksKey)
	uniqueVisitsCmd := pipe.Get(ctx, uniqueVisitsKey)
	lastAccessedCmd := pipe.Get(ctx, lastAccessedKey)
	firstAccessedCmd := pipe.Get(ctx, firstAccessedKey)
	uniqueIPsCmd := pipe.SCard(ctx, uniqueIPKey)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// Create URLAnalytics struct
	analytics := &URLAnalytics{}

	// Total clicks
	if totalClicks, err := totalClicksCmd.Int64(); err == nil {
		analytics.TotalClicks = totalClicks
	}

	if uniqueIPCount, err := uniqueIPsCmd.Result(); err == nil {
		analytics.UniqueVisits = uniqueIPCount
	}

	// Number of unique visits
	if uniqueVisits, err := uniqueVisitsCmd.Int64(); err == nil {
		analytics.UniqueVisits = uniqueVisits
	}

	// Last access time
	if lastAccessed, err := lastAccessedCmd.Result(); err == nil {
		analytics.LastAccessed, _ = time.Parse(time.RFC3339, lastAccessed)
	}

	// First access time
	if firstAccessed, err := firstAccessedCmd.Result(); err == nil {
		analytics.FirstAccessed, _ = time.Parse(time.RFC3339, firstAccessed)
	}

	return analytics, nil
}
