# 任务调度

注册一次性延迟任务与 cron 表达式任务，时钟推进触发执行（HTTP 回调或本地 Handler），失败重试、历史记录；管理页创建/暂停/查看。


## 环境

- 镜像：`benzhi.Dockerfile` 基于 `golang:1.22`（官方多架构）
- `go.mod` 语言版本：go 1.22
- 容器内使用镜像自带工具链即可

## 标准命令

```bash
go build ./...
go test ./... -count=1
go vet ./...
```

## 构建评测镜像（须双架构）

验证请用 `bash -c`（勿用 `bash -lc`）。

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-jobcron linux/amd64
docker run --platform linux/amd64 --rm go-jobcron:latest bash -c 'go build ./...'

./build_benzhi_docker.sh go-jobcron linux/arm64
docker run --platform linux/arm64 --rm go-jobcron:latest bash -c 'go build ./...'
```

构建阶段已 `go mod download`；容器内编译不应再出现 `downloading ...`。
