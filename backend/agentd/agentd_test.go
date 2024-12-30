package agentd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAgentRunCmd(t *testing.T) {
	agent := &Agent{}
	ctx := context.Background()
	returnCode, stdout, stderr, err := agent.RunCmd(ctx, "ls", []string{}, []string{})
	assert.Equal(t, 0, returnCode)
	fmt.Println(returnCode)
	fmt.Println(stdout)
	fmt.Println(stderr)
	fmt.Println(err)
	t.Log(stdout)
}
