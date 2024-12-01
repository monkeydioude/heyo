package async

import (
	"context"
	"time"
)

func Timeout(dur time.Duration, trial, fail func() error) error {
	errChan := make(chan error)
	stop := make(chan bool)
	ctx, cancelFunc := context.WithTimeout(context.TODO(), dur)
	defer cancelFunc()

	go func() {
		if err := trial(); err != nil {
			errChan <- err
			return
		}
		stop <- true
	}()
	select {
	case <-ctx.Done():
		return fail()
	case <-stop:
		return nil
	case err := <-errChan:
		return err
	}
}
