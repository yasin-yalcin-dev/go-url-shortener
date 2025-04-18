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
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// setupMockRedis creates a mock Redis server for testing
func setupMockRedis() (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

// TestRecordURLAccess tests recording URL access
// TestRecordURLAccess tests recording URL access
func TestRecordURLAccess(t *testing.T) {
	// Setup mock Redis
	mr, client := setupMockRedis()
	defer mr.Close()
	defer client.Close()

	// Create analytics store
	store := NewAnalyticsStore(client)

	testCases := []struct {
		name           string
		shortID        string
		ipAddresses    []string
		expectedUnique int
	}{
		{
			name:           "Single Access",
			shortID:        "test-url-1",
			ipAddresses:    []string{"192.168.1.1"},
			expectedUnique: 1,
		},
		{
			name:           "Multiple Unique IPs",
			shortID:        "test-url-2",
			ipAddresses:    []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"},
			expectedUnique: 3,
		},
		{
			name:           "Repeated IP",
			shortID:        "test-url-3",
			ipAddresses:    []string{"192.168.1.1", "192.168.1.1", "192.168.1.1"},
			expectedUnique: 1,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Record accesses
			for _, ip := range tc.ipAddresses {
				err := store.RecordURLAccess(ctx, tc.shortID, ip)
				if err != nil {
					t.Fatalf("RecordURLAccess failed: %v", err)
				}
			}

			// Get analytics
			analytics, err := store.GetURLAnalytics(ctx, tc.shortID)
			if err != nil {
				t.Fatalf("GetURLAnalytics failed: %v", err)
			}

			// Verify total clicks
			expectedTotalClicks := len(tc.ipAddresses)

			if int(analytics.TotalClicks) != expectedTotalClicks {
				t.Errorf("Expected total clicks %d, got %d", expectedTotalClicks, analytics.TotalClicks)
			}

			if int(analytics.UniqueVisits) != tc.expectedUnique {
				t.Errorf("Expected unique visits %d, got %d", tc.expectedUnique, analytics.UniqueVisits)
			}

			// Check timestamps
			if analytics.FirstAccessed.IsZero() {
				t.Error("FirstAccessed should not be zero")
			}
			if analytics.LastAccessed.IsZero() {
				t.Error("LastAccessed should not be zero")
			}
		})
	}
}

// uniqueIPs returns unique IP addresses
func uniqueIPs(ips []string) []string {
	unique := make(map[string]bool)
	var result []string
	for _, ip := range ips {
		if !unique[ip] {
			unique[ip] = true
			result = append(result, ip)
		}
	}
	return result
}

// Benchmark analytics store performance
func BenchmarkAnalyticsStore(b *testing.B) {
	mr, client := setupMockRedis()
	defer mr.Close()
	defer client.Close()

	store := NewAnalyticsStore(client)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shortID := "bench-url"
		ip := "192.168.1.1"

		err := store.RecordURLAccess(ctx, shortID, ip)
		if err != nil {
			b.Fatalf("RecordURLAccess failed: %v", err)
		}
	}
}
