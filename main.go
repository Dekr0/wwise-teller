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
)

func main() {
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

	utils.InitTmp()
	waapi.InitWEMCache()
	defer utils.CleanTmp()
	defer waapi.CleanWEMCache()
	if err := ui.Run(); err != nil {
		slog.Error("Failed to launch GUI", "error", err)
	}
}
