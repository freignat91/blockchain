package gnode

import (
	"fmt"
)

func (g *GNode) addEntry(mes *AntMes) error {
	payload := mes.Data
	fmt.Printf("payload: %s\n", string(payload))
	if err := g.key.verifyUserSignature(mes.UserName, mes.Key, payload); err != nil {
		return fmt.Errorf("User %s not authenticated: %v", mes.UserName, err)
	}
	answer := g.createAnswer(mes, true)
	g.senderManager.sendMessage(answer)
	return nil
}
