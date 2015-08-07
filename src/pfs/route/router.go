package route

import (
	"fmt"
	"math/rand"

	"github.com/pachyderm/pachyderm/src/pkg/grpcutil"
	"google.golang.org/grpc"
)

type router struct {
	addresser    Addresser
	dialer       grpcutil.Dialer
	localAddress string
}

func newRouter(
	addresser Addresser,
	dialer grpcutil.Dialer,
	localAddress string,
) *router {
	return &router{
		addresser,
		dialer,
		localAddress,
	}
}

func (r *router) GetMasterShards() (map[int]bool, error) {
	return r.addresser.GetMasterShards(r.localAddress)
}

func (r *router) GetSlaveShards() (map[int]bool, error) {
	return r.addresser.GetSlaveShards(r.localAddress)
}

func (r *router) GetMasterClientConn(shard int) (*grpc.ClientConn, error) {
	address, err := r.getAddress(shard, r.addresser.GetMasterShards)
	if err != nil {
		return nil, err
	}
	if address == "" {
		return nil, fmt.Errorf("no master found for %d", shard)
	}
	return r.dialer.Dial(address)
}

func (r *router) GetMasterOrSlaveClientConn(shard int) (*grpc.ClientConn, error) {
	address, err := r.getAddress(shard, r.addresser.GetSlaveShards)
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = r.getAddress(shard, r.addresser.GetMasterShards)
		if err != nil {
			return nil, err
		}
		if address == "" {
			return nil, fmt.Errorf("no slave or master found for %d", shard)
		}
	}
	return r.dialer.Dial(address)
}

func (r *router) GetAllSlaveClientConns(shard int) ([]*grpc.ClientConn, error) {
	addresses, err := r.getAddresses(shard, r.addresser.GetSlaveShards)
	if err != nil {
		return nil, err
	}
	var result []*grpc.ClientConn
	for _, address := range addresses {
		conn, err := r.dialer.Dial(address)
		if err != nil {
			return nil, err
		}
		result = append(result, conn)
	}
	return result, nil
}

func (r *router) getAddresses(shard int, testFunc func(string) (map[int]bool, error)) ([]string, error) {
	addresses, err := r.addresser.GetAllAddresses()
	if err != nil {
		return nil, err
	}
	var foundAddresses []string
	for _, address := range addresses {
		shards, err := testFunc(address)
		if err != nil {
			return nil, err
		}
		if _, ok := shards[shard]; ok {
			foundAddresses = append(foundAddresses, address)
		}
	}
	return foundAddresses, nil
}

func (r *router) getAddress(shard int, testFunc func(string) (map[int]bool, error)) (string, error) {
	addresses, err := r.getAddresses(shard, testFunc)
	if err != nil {
		return "", err
	}
	if len(addresses) == 0 {
		return "", nil
	}
	return addresses[int(rand.Uint32())%len(addresses)], nil
}

func (r *router) GetAllClientConns() ([]*grpc.ClientConn, error) {
	addresses, err := r.addresser.GetAllAddresses()
	if err != nil {
		return nil, err
	}
	clientConns := make([]*grpc.ClientConn, len(addresses)-1)
	j := 0
	for _, address := range addresses {
		// TODO(pedge): huge race, this whole thing is bad
		if address != r.localAddress {
			clientConn, err := r.dialer.Dial(address)
			if err != nil {
				return nil, err
			}
			clientConns[j] = clientConn
			j++
		}
	}
	return clientConns, nil
}
