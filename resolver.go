package resolver

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	Scheme = "consul"
)

type ConsulResolver struct {
	lock          sync.RWMutex
	target        resolver.Target
	cc            resolver.ClientConn
	consul        *api.Client
	state         chan resolver.State
	done          chan struct{}
	watchInterval time.Duration
}

func (r *ConsulResolver) ResolveNow(resolver.ResolveNowOption) {
	r.resolve()
}

func (r *ConsulResolver) Close() {
	close(r.done)
}

func (r *ConsulResolver) updater() {
	for {
		select {
		case state := <-r.state:
			r.cc.UpdateState(state)
		case <-r.done:
			return
		}
	}
}

func (r *ConsulResolver) watcher() {
	ticker := time.NewTicker(r.watchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.resolve()
		case <-r.done:
			return
		}
	}
}

func (r *ConsulResolver) resolve() {
	r.lock.Lock()
	defer r.lock.Unlock()

	state := resolver.State{}

	switch r.target.Authority {
	case "service":
		parts := strings.SplitN(r.target.Endpoint, "/", 2)

		var tag string
		if len(parts) == 2 {
			tag = parts[1]
		}

		services, _, err := r.consul.Catalog().Service(parts[0], tag, nil)
		if err != nil {
			return
		}

		addresses := make([]resolver.Address, 0, len(services))

		for _, s := range services {
			address := s.ServiceAddress
			port := s.ServicePort

			if address == "" {
				address = s.Address
			}

			addresses = append(addresses, resolver.Address{
				Addr:       address + ":" + strconv.Itoa(port),
				ServerName: r.target.Endpoint,
			})
		}

		state.Addresses = addresses
	case "query":
		queryResp, _, err := r.consul.PreparedQuery().Execute(r.target.Endpoint, nil)
		if err != nil {
			return
		}

		addresses := make([]resolver.Address, 0, len(queryResp.Nodes))

		for _, s := range queryResp.Nodes {
			address := s.Service.Address
			port := s.Service.Port

			if address == "" {
				address = s.Node.Address
			}

			addresses = append(addresses, resolver.Address{
				Addr:       address + ":" + strconv.Itoa(port),
				ServerName: r.target.Endpoint,
			})
		}

		state.Addresses = addresses
	default:
		return
	}

	r.state <- state
}
