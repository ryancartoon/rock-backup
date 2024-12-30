package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	pb "rockbackup/proto"

	"rockbackup/cmd/agent/scan"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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

	logger.Infof("received request %v", req)

	cmd := exec.Command(req.Name, req.Args...)
	cmd.Env = append(cmd.Env, req.Envs...)

	logger.Infof("running command: %s", cmd.String())
	stdout, err := cmd.Output()

	if err != nil {
		logger.Errorf("running command error: %v", err)
		return nil, err
	}

	logger.Infof("%v", string(stdout))

	return &pb.RunCmdReply{
		ReturnCode: int32(cmd.ProcessState.ExitCode()),
		Stdout:     stdout,
	}, nil
}

func (a *Agent) Scan(ctx context.Context, req *pb.ScanRequest) (*pb.ScanReply, error) {
	logger.Infof("received scan request %v", req)
	var (
		res *pb.ScanReply
	)
	// res.FileMetas = []*pb.FileMeta{}
	scanner := scan.NewLogScaner()

	metas, _ := scanner.Scan(ctx, req.Path, req.StartTime.AsTime())

	for _, meta := range metas {
		res.FileMetas = append(res.FileMetas, &pb.FileMeta{
			Path: meta.Path,
			Size: meta.Size,
		})
	}

	return res, nil
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
