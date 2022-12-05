package cdc_listener

import (
	"context"
	"github.com/dailydotdev/platform-go-common/ext/gin_server"
	"github.com/dailydotdev/platform-go-common/ext/handler"
	gpubsubstream "github.com/dailydotdev/platform-go-common/ext/stream/gpubsub"
	"github.com/dailydotdev/platform-go-common/lib/cdc_handler"
	"github.com/dailydotdev/platform-go-common/util"
	"github.com/dailydotdev/platform-go-common/util/tracing"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

// GetCmd returns a Cobra cmd which runs feed server.
func GetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "cdc_listener",
		Short:   "Listens and processes cdc events",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {

			svc, err := gin_server.NewServer(cmd.Context(), gin_server.WithHealthCheck())
			if err != nil {
				log.Fatal().Msgf("gin_server.NewServer(): %s", err)
			}

			svc.Use(tracing.GetMiddleware(cmd.Use))

			cdcSubscriptionName, err := util.EnvGet[string]("CDC_SUBSCRIPTION_NAME")
			if err != nil {
				log.Fatal().Msgf("init: %s", err)
			}

			table3subscriptionName, err := util.EnvGet[string]("TABLE3_SUBSCRIPTION_NAME")
			if err != nil {
				log.Fatal().Msgf("init: %s", err)
			}

			gpubsub, err := gpubsubstream.NewProvider(log.Logger)
			if err != nil {
				log.Fatal().Msgf("gpubsub_stream.NewProvider(): %s", err)
			}

			sub, err := gpubsub.GetSubscriber(cmd.Context(), cdcSubscriptionName,
				cdc_handler.ProcessStream(
					cdc_handler.Register(
						regexp.MustCompile(`.*`),
						cdc_handler.Create,
						func(ctx context.Context, message cdc_handler.CDCMessage) error {
							log.Debug().Str("subscription", cdcSubscriptionName).Interface("message", message).Send()
							return nil
						}),
				),
				gpubsubstream.MethodName(handler.MethodStream),
			)
			if err != nil {
				log.Fatal().Msgf("gpubsub.Subscribe(): %s", err)
			}

			sub1, err := gpubsub.GetSubscriber(cmd.Context(), table3subscriptionName,
				cdc_handler.ProcessStream(
					cdc_handler.Register(
						regexp.MustCompile(`.*`),
						cdc_handler.Create,
						func(ctx context.Context, message cdc_handler.CDCMessage) error {
							log.Debug().Str("subscription", table3subscriptionName).Interface("message", message).Send()
							return nil
						}),
				),
				gpubsubstream.MethodName(handler.MethodStream),
			)
			if err != nil {
				log.Fatal().Msgf("gpubsub.Subscribe(): %s", err)
			}

			if err := svc.Start(); err != nil {
				log.Fatal().Msgf("svc.Start: %v", err)
			}

			// Wait for signal
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
			signal.NotifyContext(cmd.Context())
			<-signalCh
			sub.Stop()
			sub1.Stop()

			if err := svc.Stop(); err != nil {
				log.Fatal().Msgf("svc.Stop: %v", err)
			}

		},
	}

	return cmd
}
