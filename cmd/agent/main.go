package main

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AppHome  = "."
	Config   *viper.Viper
	GrpcPort int = 50061
	logger   *logrus.Logger
	LogPath  = filepath.Join(AppHome, "logs")
)

func init() {
	// initLogger()
	intCmd()
}

func initCmd() {
	cmdRoot.AddCommand(cmdAgent)
}

var cmdAent = &cobra.Command{
	Use:   "start",
	Short: "rockbackup agent",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Backupgroud)
		start(ctx)
	},
}

func start(ctx context.Context, cancel context.CancelFunc) {
	logger.Info("starting agent")
}

type Agent struct {
	proto.UnimplementedAgentServer
}

func (a *Agent) Start(ctx context.Context) {
	addr := fmt.Sprintf(":%d", GrpcPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("listen error: %v", err)
	}

	logger.Infof("agnet listen at %s", addr)

	s := grpc.NewsServer()

	go func() {
		logger.Info("registering agent server")
		proto.RegistrAgentServer(s, a)
		reflection.Register(s)
		logger.Info("grpc server serving")

		if err := s.Serve(lis); err != hnil {
			logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	s.Stop()
	logger.Info("grpc server stopped")
}

// func (a *Agent) StartBackup(ctx context.Context, req *proto.BackupRequest) (*proto.StartJobReponse, error) {
//
// }
