package check

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sinspired/subs-check-pro/v3/config"
)

// latencyWithinLimit 判断一次 HTTP 响应耗时是否在配置阈值内。
func latencyWithinLimit(elapsed time.Duration, maxLatencyMS int) bool {
	return maxLatencyMS > 0 && elapsed >= 0 && elapsed <= time.Duration(maxLatencyMS)*time.Millisecond
}

// filterResultsByLatency 在原有检测流水线完成后，使用每个节点自己的代理连接测量 HTTP 响应时间。
// 此处测得的是 HTTP 请求响应时间（含代理连接建立、TLS 和服务端响应），不是 ICMP Ping。
func (pc *ProxyChecker) filterResultsByLatency() {
	cfg := config.GlobalConfig
	if !cfg.LatencyFilter {
		return
	}
	if cfg.MaxLatency <= 0 {
		slog.Warn("已启用延迟筛选，但 max-latency 必须大于 0；本轮跳过筛选")
		return
	}
	if len(pc.results) == 0 {
		slog.Info("延迟筛选：没有可筛选的可用节点")
		return
	}

	testURL := cfg.LatencyTestURL
	if testURL == "" {
		testURL = "https://www.gstatic.com/generate_204"
	}
	parsedURL, err := url.ParseRequestURI(testURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		slog.Warn("延迟测试地址无效，本轮跳过筛选", "url", testURL)
		return
	}

	timeoutMS := cfg.LatencyTestTimeout
	if timeoutMS <= 0 {
		timeoutMS = 5000
	}
	concurrency := cfg.LatencyConcurrent
	if concurrency <= 0 {
		concurrency = 20
	}
	concurrency = min(concurrency, min(len(pc.results), 64))

	requestTimeout := time.Duration(timeoutMS) * time.Millisecond
	total := len(pc.results)
	keep := make([]bool, total)

	var next atomic.Int64
	var completed atomic.Int32
	var passed atomic.Int32
	var cancelled atomic.Bool
	filterCtx, cancel := context.WithCancel(context.Background())
	monitorDone := make(chan struct{})

	CurrentStepName.Store("延迟筛选")
	ProxyCount.Store(uint32(total))
	Progress.Store(0)
	Available.Store(0)

	// 用户手动停止检测时，取消延迟请求并保留原始检测结果，避免只筛选一半。
	go func() {
		defer close(monitorDone)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-filterCtx.Done():
				return
			case <-ticker.C:
				if ForceClose.Load() {
					cancelled.Store(true)
					cancel()
					return
				}
			}
		}
	}()

	slog.Info("开始独立延迟筛选",
		"节点数", total,
		"最大延迟", cfg.MaxLatency,
		"测试地址", testURL,
		"超时毫秒", timeoutMS,
		"并发数", concurrency,
	)

	var wg sync.WaitGroup
	for range concurrency {
		wg.Go(func() {
			for {
				if filterCtx.Err() != nil {
					return
				}
				index := int(next.Add(1) - 1)
				if index >= total {
					return
				}

				proxyClient := CreateClient(pc.results[index].Proxy)
				if proxyClient != nil {
					proxyClient.Timeout = requestTimeout
					requestCtx, requestCancel := context.WithTimeout(filterCtx, requestTimeout)
					req, reqErr := http.NewRequestWithContext(requestCtx, http.MethodGet, testURL, nil)
					if reqErr == nil {
						req.Header.Set("User-Agent", "clash.meta")
						started := time.Now()
						resp, requestErr := proxyClient.Do(req)
						elapsed := time.Since(started)
						if resp != nil {
							_ = resp.Body.Close()
						}
						keep[index] = requestErr == nil &&
							resp != nil &&
							resp.StatusCode >= http.StatusOK &&
							resp.StatusCode < http.StatusMultipleChoices &&
							latencyWithinLimit(elapsed, cfg.MaxLatency)
					}
					requestCancel()
					proxyClient.Close()
				}

				done := completed.Add(1)
				if keep[index] {
					passed.Add(1)
				}
				Progress.Store(uint32(done))
				Available.Store(uint32(passed.Load()))
			}
		})
	}
	wg.Wait()
	cancel()
	<-monitorDone

	if cancelled.Load() || ForceClose.Load() {
		slog.Warn("用户中止延迟筛选，本轮保留未筛选前的全部可用节点")
		return
	}

	filtered := make([]Result, 0, int(passed.Load()))
	for i, result := range pc.results {
		if keep[i] {
			filtered = append(filtered, result)
		}
	}
	dropped := total - len(filtered)
	pc.results = filtered

	slog.Info("延迟筛选完成",
		"筛选前", total,
		"保留", len(filtered),
		"剔除", dropped,
		"最大延迟毫秒", cfg.MaxLatency,
	)
}
