package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"

	"github.com/dailydotdev/debezium_test/cmd/cdc_listener"
	_ "github.com/dailydotdev/platform-go-common/ext/gcp"
	"github.com/dailydotdev/platform-go-common/util/logging"
	"github.com/dailydotdev/platform-go-common/util/metrics"
	"github.com/dailydotdev/platform-go-common/util/migrations"
	"github.com/dailydotdev/platform-go-common/util/tracing"
)

var (
	rootCmd = &cobra.Command{
		Use:   "debezium_test",
		Short: "debezium test project",
	}
)

func main() {

	if err := logging.Configure(); err != nil {
		log.Fatal().Err(err).Msg("init: logging.Configure()")
	}

	if err, finalizer := tracing.Configure(); err != nil {
		log.Fatal().Err(err).Msg("init: tracing.Configure()")
	} else if finalizer != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		defer finalizer(ctx)
	}

	if err, finalizer := metrics.Configure(rootCmd.Use); err != nil {
		log.Fatal().Err(err).Msg("init: metrics.Configure()")
	} else if finalizer != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		defer finalizer(ctx)
	}

	rootCmd.AddCommand(cdc_listener.GetCmd())
	rootCmd.AddCommand(migrations.AddPostgresMigration())

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("rootCmd.Execute: %w", err))
	}
}
