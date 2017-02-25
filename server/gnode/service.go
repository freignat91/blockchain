package gnode

import (
	"fmt"
	"golang.org/x/net/context"
	"net"
)

// Send send message to grid
func (g *GNode) ExecuteFunction(ctx context.Context, mes *AntMes) (*AntRet, error) {
	if mes.Id == "" {
		mes.Id = g.getNewId(false)
		mes.Origin = g.name
		//logf.debugMes(mes, "Received message from client: %v\n", mes.Id)
	} else {
		if ok := g.idMap.Exists(mes.Id); ok {
			//logf.info("execute store bloc ack doublon id=%s order=%d\n", mes.Id, mes.Order)
			return &AntRet{Ack: true}, nil
		}
		g.idMap.Add(mes.Id)
	}
	if !g.receiverManager.receiveMessage(mes) {
		return &AntRet{Ack: false}, nil
	}
	return &AntRet{Ack: true}, nil
}

func (g *GNode) Ping(ctx context.Context, mes *AntMes) (*PingRet, error) {
	return &PingRet{
		Name:         g.name,
		Host:         g.host,
		NbNode:       int32(g.nbNode),
		NbDuplicate:  int32(config.nbDuplicate),
		ClientNumber: int32(g.clientMap.len()),
	}, nil
}

func (g *GNode) CheckReceiver(ctx context.Context, req *HealthRequest) (*AntRet, error) {
	return &AntRet{Ack: g.receiverManager.buffer.isAvailable()}, nil
}

func (g *GNode) Healthcheck(ctx context.Context, req *HealthRequest) (*AntRet, error) {
	if !g.healthy {
		return nil, fmt.Errorf("Not ready")
	}
	return &AntRet{Ack: true}, nil
}

func (g *GNode) GetClientStream(stream GNodeService_GetClientStreamServer) error {
	if !g.healthy {
		return fmt.Errorf("Node %s not yet ready", g.name)
	}
	g.receiverManager.startClientReader(stream)
	return nil
}

func (g *GNode) AskConnection(ctx context.Context, req *AskConnectionRequest) (*PingRet, error) {
	ret := &PingRet{
		Name: g.name,
		Host: g.host,
	}
	ip := net.ParseIP(req.Ip)
	if ip == nil {
		return nil, fmt.Errorf("IP addresse parse error %s", req.Ip)
	}
	if _, ok := g.targetMap[req.Name]; ok {
		return ret, nil
	}
	conn, err := g.startGRPCClient(ip)
	if err != nil {
		return nil, fmt.Errorf("Start GRPC error: %v", err)
	}
	client := NewGNodeServiceClient(conn)
	target := &gnodeTarget{
		name:         req.Name,
		host:         req.Host,
		ip:           req.Ip,
		conn:         conn,
		client:       client,
		updateNumber: g.updateNumber,
	}
	g.targetMap[req.Name] = target
	return ret, nil
}
