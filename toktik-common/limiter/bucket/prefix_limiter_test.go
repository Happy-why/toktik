package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPrefixLimiter(t *testing.T) {
	p := NewPrefixLimiter()
	p.AddBuckets(BucketRule{
		Key:          "/post/test1",
		FillInterval: time.Second,
		Cap:          100,
		Quantum:      100,
	}, BucketRule{
		Key:          "/post/test2",
		FillInterval: time.Second,
		Cap:          10,
		Quantum:      10,
	}, BucketRule{
		Key:          "/user/test3",
		FillInterval: time.Second,
		Cap:          20,
		Quantum:      20,
	})
	key := p.testKey("/post/test1")
	bucket, ok := p.GetBucket(key)
	require.True(t, ok)
	require.NotEmpty(t, bucket)
	key = p.testKey("/post/test2")
	bucket, ok = p.GetBucket(key)
	require.True(t, ok)
	require.NotEmpty(t, bucket)
	key = p.testKey("/post/test3")
	bucket, ok = p.GetBucket(key)
	require.False(t, ok)
	require.Empty(t, bucket)
}
