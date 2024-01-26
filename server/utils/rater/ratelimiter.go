package rater

import (
	"errors"
	"sync"
	"time"
)

type RateLimter interface {
	TryAcquire() bool
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	limit        int           // 窗口请求上限
	window       int64         // 窗口时间大小
	smallWindow  int64         // 小窗口时间大小
	smallWindows int64         // 小窗口数量
	counters     map[int64]int // 小窗口计数器
	mutex        sync.Mutex    // 避免并发问题
}

var _ RateLimter = &SlidingWindowLimiter{}

func NewSlidingWindowLimiter(limit int, window, smallWindow time.Duration) (*SlidingWindowLimiter, error) {
	// 窗口时间必须能够被小窗口时间整除
	if window%smallWindow != 0 {
		return nil, errors.New("window cannot be split by integers")
	}

	return &SlidingWindowLimiter{
		limit:        limit,
		window:       int64(window),
		smallWindow:  int64(smallWindow),
		smallWindows: int64(window / smallWindow),
		counters:     make(map[int64]int),
	}, nil
}

func (l *SlidingWindowLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currentSmallWindow := time.Now().UnixNano() / l.smallWindow * l.smallWindow
	startSmallWindow := currentSmallWindow - l.smallWindow*(l.smallWindows-1)

	// count
	var count int
	for smallWindow, counter := range l.counters {
		if smallWindow < startSmallWindow {
			delete(l.counters, smallWindow)
		} else {
			count += counter
		}
	}

	// 若到达窗口请求上限，请求失败
	if count >= l.limit {
		return false
	}
	// 若没到窗口请求上限，当前小窗口计数器+1，请求成功
	l.counters[currentSmallWindow]++
	return true
}

func (l *SlidingWindowLimiter) StartTimer(sig chan struct{}) {
	ticker := time.NewTicker(time.Duration(l.smallWindow / 2))
	for {
		select {
		case <-ticker.C:
			if l.TryAcquire() {
				sig <- struct{}{}
			}
		}
	}
}
