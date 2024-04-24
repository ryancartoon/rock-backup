package agentd

import (
	"google.golang.org/grpc"
	"net"
	"os"
)

const (
	PORT = ":50001"
)

type Agent struct {
	Host string
	Port uint
}

type Agentd struct {
	Agents []Agent
}

func (a Agentd) GetAgent(host string) (Agent, error) {
	return Agent{}, nil
}

func (a *Agent) RunCmd() (returnCode int, stdout []byte, stderr []byte) {
	return
}

func (a *Agent) Connect() error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRemoteServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.RunCmd(context.Background(), &pb.CommandRequest{Command: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
