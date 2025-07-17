package main

import (
	"context"
	"flag"
	"log/slog"
	"time"

	"github.com/Dekr0/wwise-teller/automation"
	"github.com/Dekr0/wwise-teller/ui"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/waapi"
	"github.com/Dekr0/wwise-teller/wwise"
)

func main() {
	wwise.InitTranslation()

	proc := flag.String("proc", "", "Filepath to sound bank processor pipelines specification")
	procDeadline := flag.Uint64("deadline", 16, "Deadline in seconds of running sound bank processor pipelines")

	flag.Parse()

	if *proc != "" {
		defer utils.CleanTmp()
		utils.InitTmp()
		if *procDeadline == 0 {
			automation.Process(context.Background(), *proc)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second * time.Duration(*procDeadline))
			defer cancel()
			automation.Process(ctx, *proc)
		}
		return
	}

	if err := ui.Run(); err != nil {
		slog.Error("Failed to launch GUI", "error", err)
	}
	slog.Info("Application Closed")
	utils.CleanTmp()
	waapi.CleanWEMCache()
}
