package agentd

import (
	"context"
	"fmt"
	pb "rockbackup/proto"

	"google.golang.org/grpc"
)

const (
	PORT = ":50001"
)

type Agent struct {
	Host    string
	Port    uint
	gClient pb.AgentClient
	gConn   *grpc.ClientConn
}

type Agentd struct {
	Agents []*Agent
}

func (a Agentd) GetAgent(host string) (Agent, error) {
	return Agent{}, nil
}

func (a *Agent) RunCmd(ctx context.Context, name string, args []string, envs []string) (int, []byte, []byte, error) {
	err := a.Connect()
	// defer a.Close()

	if err != nil {
		logger.Errorf("connect error: %v", err)
		return 0, nil, nil, err
	}

	resp, err := a.gClient.RunCmd(ctx, &pb.RunCmdRequest{Name: name, Args: args, Envs: envs})

	// Contact the server and print out its response.
	if err != nil {
		logger.Errorf("could not run command: %v", err)
		return 0, nil, nil, err
	}

	return int(resp.ReturnCode), resp.Stdout, resp.Stderr, nil
}

func (a *Agent) Connect() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", a.Host, a.Port), grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	a.gClient = pb.NewAgentClient(conn)

	return nil
}

func (a *Agent) Close() error {
	return a.gConn.Close()
}
