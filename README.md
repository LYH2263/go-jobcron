# go-jobcron

Go 实现的延迟与 Cron 任务调度器：Enqueue 入队，Tick 推进到期队列，RunDue 执行 handler 或 HTTP 回调，可选 JSON 快照与管理页。

## 运行

```bash
go test ./... -count=1
go run ./cmd/jobd -addr :8093 -web web -url http://127.0.0.1:9999/hook
```

浏览器打开 `http://localhost:8093/`：创建任务、查看列表、暂停/恢复、触发执行。

## 库面

```go
s := jobcron.New(jobcron.WithDefaultURL("https://example.com/hook"))
id, err := s.Enqueue("notify", []byte(`{"u":1}`), time.Now().Add(time.Minute))
results, err := s.RunDueContext(ctx, 10)
_ = s.ListJobs(20)
_ = s.Close()
```
