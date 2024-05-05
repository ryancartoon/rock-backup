package main

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	pb "rockbackup/proto"
)

var (
	AppHome  = "."
	Config   *viper.Viper
	GrpcPort int = 50001
	logger   *logrus.Logger
	LogPath  = filepath.Join(AppHome, "logs")
)

var cmdRoot = &cobra.Command{
	Use:   "rock_agent",
	Short: "rock agent",
	Long:  "rock agent",
}

func init() {
	initLogger()
	initCmd()
}

func initCmd() {
	cmdRoot.AddCommand(cmdAgent)
}

func initLogger() {
	logger = logrus.New()
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(LogPath, "rock.log"),
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	logger.SetOutput(io.MultiWriter(logFile, os.Stdout))
	logger.SetLevel(logrus.DebugLevel)
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

	agent := &Agent{}
	agent.Serve(ctx)
}

type Agent struct {
	pb.UnimplementedAgentServer
}

func (a *Agent) Serve(ctx context.Context) {
	addr := fmt.Sprintf(":%d", GrpcPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("listen error: %v", err)
	}

	logger.Infof("agnet listen at %s", addr)

	s := grpc.NewServer()

	go func() {
		logger.Info("registering agent server")
		pb.RegisterAgentServer(s, a)
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

func (a *Agent) RunCmd(ctx context.Context, req *pb.RunCmdRequest) (*pb.RunCmdReply, error) {
	cmd := exec.Command(req.Name, req.Args...)
	cmd.Env = append(cmd.Env, req.Envs...)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	logger.Infof("%v", string(stdout))

	return &pb.RunCmdReply{
		ReturnCode: int32(cmd.ProcessState.ExitCode()),
		Stdout:     stdout,
	}, nil
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
