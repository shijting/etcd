package main

import "github.com/shijting/etcd/util"

func main()  {
	cli := util.NewClient()
	cli.GetService()
}
