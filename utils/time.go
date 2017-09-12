package utils

import (
	"fmt"
	"time"
)

// LoopRunUntilSuccessOrTimeout to run until OK or timeout
func LoopRunUntilSuccessOrTimeout(fn func() bool, timeout, interval time.Duration) error {
	expiredTime := time.Now().Add(timeout)
	tk := time.NewTicker(interval)
	defer tk.Stop()

	for {
		if fn() {
			return nil
		}

		now := <-tk.C
		if now.After(expiredTime) {
			return fmt.Errorf("timeout after %v", timeout)
		}
	}
}
