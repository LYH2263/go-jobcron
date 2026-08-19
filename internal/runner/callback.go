package runner

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/LYH2263/go-jobcron/internal/httpx"
	"github.com/LYH2263/go-jobcron/internal/store"
)

// HTTPCallback 执行 HTTP POST 回调。
func HTTPCallback(ctx context.Context, client *httpx.Client, rec store.Record, ua string) (int, error) {
	if client == nil {
		return 0, ErrNilClient
	}
	url := rec.URL
	if url == "" {
		return 0, Wrap(ErrHTTP, "empty url")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rec.Payload))
	if err != nil {
		return 0, WrapErr(ErrHTTP, err)
	}
	req.Header.Set("Content-Type", "application/json")
	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	req.Header.Set("X-Job-Name", rec.Name)
	req.Header.Set("X-Job-ID", rec.ID)
	resp, err := client.Do(req)
	if err != nil {
		return 0, WrapErr(ErrHTTP, err)
	}
	defer DrainAndClose(resp.Body)
	code := resp.StatusCode
	if code < 200 || code >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return code, Wrap(ErrHTTP, string(body))
	}
	return code, nil
}
