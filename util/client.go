package util

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"log"
	"regexp"
	"time"
)

type Client struct {
	Client *clientv3.Client
	Services []*ServiceInfo
}
type ServiceInfo struct {
	ServiceID string
	ServiceName string
	ServiceAddr string
}
func NewClient() *Client  {
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:            []string{"106.53.5.146:23791", "106.53.5.146:23792"},
		DialKeepAliveTime:    10 * time.Second,
	})
	if err != nil {
		// handle errors
		log.Fatal(err)
	}
	return &Client{Client: cli}
}
// 服务发现
func (t *Client)GetService()  {
	kv := clientv3.NewKV(t.Client)
	ctx := context.Background()
	keyPrefix := "/services"
	// 取出以 /services 为前缀的所有key
	getResponse , err := kv.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, item := range getResponse.Kvs {
		//fmt.Println(string(item.Key))
		t.parseService(item.Key,item.Value)
	}
}

func(t *Client) parseService(key []byte,value []byte)  {
	// 正则取出id， name
	reg:=regexp.MustCompile("/services/(\\w+)/(\\w+)")
	if reg.Match(key){
		idandname:=reg.FindSubmatch(key)
		sid:=idandname[1]
		sname:=idandname[2]
		t.Services = append(t.Services,&ServiceInfo{ServiceID:string(sid),
			ServiceName:string(sname),ServiceAddr:string(value)})
	}
}