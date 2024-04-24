package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	pb "rockbackup/proto"
)

var (
	AppHome  = "."
	Config   *viper.Viper
	GrpcPort int = 50061
	logger   *logrus.Logger
	LogPath  = filepath.Join(AppHome, "logs")
	cmdRoot  *cobra.Command
)

func init() {
	// initLogger()
	initCmd()
}

func initCmd() {
	cmdRoot.AddCommand(cmdAgent)
}

var cmdAgent = &cobra.Command{
	Use:   "start",
	Short: "rockbackup agent",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		start(ctx, cancel)
	},
}

func start(ctx context.Context, cancel context.CancelFunc) {
	logger.Info("starting agent")
}

type Agent struct {
	pb.UnimplementedAgentServer
}

func (a *Agent) RunCMD(ctx context.Context) {
}

func Serve(ctx context.Context, agent *Agent) {
	addr := fmt.Sprintf(":%d", GrpcPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("listen error: %v", err)
	}

	logger.Infof("agnet listen at %s", addr)

	s := grpc.NewsServer()

	go func() {
		logger.Info("registering agent server")
		s.RegistrAgentServer(s, agent)
		// reflection.Register(s)
		logger.Info("grpc server serving")

		if err := s.Serve(lis); err != nil {
			logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	s.Stop()
	logger.Info("grpc server stopped")
}

func (a *Agent) RunCmd(ctx context.Context, req *pb.CmdRequest) (*pb.CmdReply, error) {
	cmd := exec.Command(req.Cmd)
	cmd.Env = append(cmd.Env, req.Env)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return &pb.CmdReply{
		ReturnCode: int32(cmd.ProcessState.ExitCode()),
		Stdout:     string(stdout),
		Stderr:     "", // Assuming you want to leave stderr empty for now
	}, nil
}
