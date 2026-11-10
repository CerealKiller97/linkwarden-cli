//go:build unix

package cmd

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Unix only: a test can't send its own process a console Ctrl+C on Windows.
func TestInterruptContextCancelsOnSignal(t *testing.T) {
	assert := require.New(t)

	ctx, stop := interruptContext()
	defer stop()

	assert.NoError(syscall.Kill(os.Getpid(), syscall.SIGINT))

	select {
	case <-ctx.Done():
		assert.ErrorIs(ctx.Err(), context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("context not cancelled by SIGINT")
	}
}
