package agentd

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
