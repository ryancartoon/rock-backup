package agentd

import (
	"context"
	pb "rockbackup/proto"

	"google.golang.org/grpc"
)

const (
	PORT = ":50001"
)

type Agent struct {
	Host string
	Port uint
	Conn pb.AgentClient
}

type Agentd struct {
	Agents []Agent
}

func (a Agentd) GetAgent(host string) (Agent, error) {
	return Agent{}, nil
}

func (a *Agent) RunCmd(ctx context.Context, name string, env string) (int, string, string, error) {
	resp, err := a.Conn.RunCmd(ctx, &pb.CmdRequest{Cmd: name, Env: env})

	// Contact the server and print out its response.
	if err != nil {
		logger.Errorf("could not run command: %v", err)
		return 0, "", "", err
	}

	return int(resp.ReturnCode), resp.Stdout, resp.Stderr, nil
}

func (a *Agent) Connect() error {
	conn, err := grpc.Dial("localhost:50001", grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	a.Conn = pb.NewAgentClient(conn)

	return nil
}
