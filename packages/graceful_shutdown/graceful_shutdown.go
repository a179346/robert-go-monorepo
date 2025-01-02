package graceful_shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

func OnShutdown(f func(singal os.Signal)) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	f(sig)
}
