package bucket

import (
	"context"
	"io"
	"sync"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/textileio/go-threads/core/thread"
	bucketsClient "github.com/textileio/textile/api/buckets/client"
	bucketsproto "github.com/textileio/textile/api/buckets/pb"
)

type BucketData struct {
	Key       string `json:"_id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	DNSRecord string `json:"dns_record,omitempty"`
	//Archives  Archives `json:"archives"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type DirEntries bucketsproto.ListPathReply

type BucketsClient interface {
	PushPath(ctx context.Context, key, pth string, reader io.Reader, opts ...bucketsClient.Option) (result path.Resolved, root path.Resolved, err error)
	PullPath(ctx context.Context, key, pth string, writer io.Writer, opts ...bucketsClient.Option) error
	ListPath(ctx context.Context, key, pth string) (*bucketsproto.ListPathReply, error)
	RemovePath(ctx context.Context, key, pth string, opts ...bucketsClient.Option) (path.Resolved, error)
}

// NOTE: all write operations should use the lock for the bucket to keep consistency
// TODO: Maybe read operations dont need a lock, needs testing
// struct for implementing bucket interface
type Bucket struct {
	lock          sync.RWMutex
	root          *bucketsproto.Root
	ctx           context.Context
	bucketsClient BucketsClient
	threadID      *thread.ID
}

func (b *Bucket) Slug() string {
	return b.GetData().Name
}

func New(root *bucketsproto.Root, ctx context.Context, bucketsClient BucketsClient) *Bucket {
	return &Bucket{
		root:          root,
		bucketsClient: bucketsClient,
		ctx:           ctx,
	}
}

func (b *Bucket) Key() string {
	return b.GetData().Key
}

func (b *Bucket) GetData() BucketData {
	return BucketData{
		Key:       b.root.Key,
		Name:      b.root.Name,
		Path:      b.root.Path,
		DNSRecord: "",
		CreatedAt: b.root.CreatedAt,
		UpdatedAt: b.root.UpdatedAt,
	}
}

func (b *Bucket) getContext() (context.Context, *thread.ID) {
	return b.ctx, b.threadID
}