package main

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"

	"time"
)


type JobClient struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}


func InitJobClient() *JobClient {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
		err error
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2389"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		return nil
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	return &JobClient{
		client: client,
		kv:     kv,
		lease:  lease,
	}
}

// 保存和更新数据
func (jobClient *JobClient) SaveJob(job *Job) (oldJob *Job, err error) {
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj Job
	)

	jobKey = "/cron/jobs/" + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	// 保存到etcd
	if putResp, err = jobClient.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	// 如果更新, 那么返回旧址
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}





func main()  {
	var (
		job Job
		oldJob *Job
		err error
	)
	c := InitJobClient()


	job = Job{
		Name: "test",
		CronExpr: "*/5 * * * * *",
		Command: "echo hello word",
	}
	if oldJob, err = c.SaveJob(&job); err != nil {
		return
	}
	fmt.Println("----->", oldJob)

}
