package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errAddrMisMatch = errors.New("consul resolver: invalied uri")

	errEndsWithColon = errors.New("consul resolver: missing port after port-separator colon")

	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")
)

func NewInit() {
	fmt.Printf("calling consul init\n")
	resolver.Register(NewBuilder())
}

type consulBuilder struct {
}

type consulResolver struct {
	address              string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	tag                  string
	disableServiceConfig bool
	lastIndex            uint64
}

func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	name := ""
	tag := ""
	temp := strings.Split(target.Endpoint(), "/")
	if len(temp) > 1 {
		name = temp[0]
		tag = temp[1]
	} else {
		name = temp[0]
	}
	cr := &consulResolver{
		address:              target.URL.Host,
		name:                 name,
		cc:                   cc,
		tag:                  tag,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	cr.wg.Add(1)
	go cr.watcher()
	return cr, nil

}

func (cr *consulResolver) watcher() {
	config := api.DefaultConfig()
	config.Address = cr.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}

	retryLimit := 5 // 设置重试次数限制
	retryCount := 0 // 初始化重试次数

	for {
		if retryCount >= retryLimit {
			fmt.Println("reached retry limit, exiting watcher")
			return // 达到重试次数限制，退出函数
		}

		services, metainfo, err := client.Health().Service(cr.name, cr.tag, true, &api.QueryOptions{
			WaitIndex: cr.lastIndex,
			WaitTime:  time.Minute,
		})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v\n", err)
			retryCount++                 // 错误时增加重试次数
			time.Sleep(time.Second * 10) // 在下次重试前等待
			continue                     // 继续下一轮重试
		}

		// 成功获取到服务信息时重置重试次数
		retryCount = 0

		cr.lastIndex = metainfo.LastIndex
		var newAddrs []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			newAddrs = append(newAddrs, resolver.Address{Addr: addr})
		}
		cr.cc.NewAddress(newAddrs)
		cr.cc.NewServiceConfig(cr.name)
	}
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
}

func (cr *consulResolver) Close() {
}
