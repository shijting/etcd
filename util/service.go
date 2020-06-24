package util

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

type Service struct {
	Client *clientv3.Client
}

const (
	ttl int64 = 60
)

func NewService() *Service {
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:            []string{"106.53.5.146:23791", "106.53.5.146:23792"},
		DialKeepAliveTime:    10 * time.Second,
	})
	if err == context.DeadlineExceeded {
		// handle errors
		log.Fatal(err)
	}
	return &Service{Client: cli}
}
// 服务注册
func (t *Service)RegService(id string, name string, addr string) (err error) {
	kv := clientv3.NewKV(t.Client)
	keyPrefix := "/services/"
	key := keyPrefix + id + "/" +name
	ctx := context.Background()
	// 设置租约
	lease := clientv3.NewLease(t.Client)
	leaseResp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return
	}
	// 存入etcd
	_, err = kv.Put(ctx, key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return
	}
	// 定时续租
	keepaliveRes, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return
	}
	go lisKeepAlive(keepaliveRes)
	return
}
// 自动续租
func lisKeepAlive(keepaliveRes <- chan *clientv3.LeaseKeepAliveResponse)  {
	for  {
		select {
		case ret := <- keepaliveRes:
			if ret !=nil {
				log.Println("续租成功")
			} else {
				log.Println("续租失败")
			}
		}
	}
}
// 反服务注册（删除)
func (t *Service)UnRegService(id string) (err error) {
	kv := clientv3.NewKV(t.Client)
	keyPrefix := "/services/"
	key := keyPrefix + id
	// 删除 keyPrefix + id 为前缀的所有key
	_, err = kv.Delete(context.Background(), key, clientv3.WithPrefix())
	log.Println("反注册...")
	return
}