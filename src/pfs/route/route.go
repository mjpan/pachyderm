package route

import (
	"github.com/pachyderm/pachyderm/src/pfs"
	"github.com/pachyderm/pachyderm/src/pkg/discovery"
	"github.com/pachyderm/pachyderm/src/pkg/grpcutil"
	"google.golang.org/grpc"
)

type Sharder interface {
	NumShards() int
	GetShard(path *pfs.Path) (int, error)
}

func NewSharder(numShards int) Sharder {
	return newSharder(numShards)
}

type Addresser interface {
	GetMasterShards(address string) (map[int]bool, error)
	GetSlaveShards(address string) (map[int]bool, error)
	GetAllAddresses() ([]string, error)
}

func NewSingleAddresser(address string, numShards int) Addresser {
	return newSingleAddresser(address, numShards)
}

func NewDiscoveryAddresser(discoveryClient discovery.Client, baseKey string) Addresser {
	return newDiscoveryAddresser(discoveryClient, baseKey)
}

type Router interface {
	GetMasterShards() (map[int]bool, error)
	GetSlaveShards() (map[int]bool, error)
	GetMasterClientConn(shard int) (*grpc.ClientConn, error)
	GetMasterOrSlaveClientConn(shard int) (*grpc.ClientConn, error)
	GetAllSlaveClientConns(shard int) (*grpc.ClientConn, error)
	GetAllClientConns() ([]*grpc.ClientConn, error)
}

func NewRouter(
	addresser Addresser,
	dialer grpcutil.Dialer,
	localAddress string,
) Router {
	return newRouter(
		addresser,
		dialer,
		localAddress,
	)
}
