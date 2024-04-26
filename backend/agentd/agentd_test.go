package agentd

import (
	"context"
	"fmt"
	"testing"
)

func TestAgentRunCmd(t *testing.T) {
	agent := &Agent{}
	agent.Connect()
	ctx := context.Background()
	returnCode, stdout, stderr, err := agent.RunCmd(ctx, "ls", "")
	fmt.Println(returnCode)
	fmt.Println(stdout)
	fmt.Println(stderr)
	fmt.Println(err)
	t.Errorf(err.Error())

}
