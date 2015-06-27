package retry

import (
	"fmt"
	"testing"
	"time"
)

func TestRetrySuccess(t *testing.T) {
	cnt := 0
	err := Retry(3, 1*time.Nanosecond, func() error {
		if cnt == 0 {
			cnt++
			return fmt.Errorf("retry")
		}
		return nil
	})

	if err != nil {
		t.Errorf("error should be nil, %s", err)
	}

	if cnt != 1 {
		t.Errorf("cnt should be 1, %d", cnt)
	}
}

func TestRetryFail(t *testing.T) {
	cnt := 0
	err := Retry(4, 1*time.Nanosecond, func() error {
		cnt++
		return fmt.Errorf("retry")
	})

	if err == nil || err.Error() != "retry" {
		t.Errorf("error should be occured")
	}

	if cnt != 4 {
		t.Errorf("cnt should be 4, but %d", cnt)
	}
}
