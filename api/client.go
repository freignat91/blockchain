package api

import (
	"fmt"
	"github.com/freignat91/blockchain/server/gnode"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

type gnodeClient struct {
	api      *BchainAPI
	id       string
	client   gnode.GNodeServiceClient
	nodeName string
	nodeHost string
	ctx      context.Context
	stream   gnode.GNodeService_GetClientStreamClient
	recvChan chan *gnode.AntMes
	lock     sync.RWMutex
	conn     *grpc.ClientConn
	//nbNode      int
	nbDuplicate int
}

func (g *gnodeClient) init(api *BchainAPI) error {
	g.api = api
	g.ctx = context.Background()
	g.recvChan = make(chan *gnode.AntMes)
	if err := g.connectServer(); err != nil {
		return err
	}
	if err := g.startServerReader(); err != nil {
		return err
	}
	api.info("Client %s connected to node %s (%s)\n", g.id, g.nodeName, g.nodeHost)
	return nil
}

func (g *gnodeClient) connectServer() error {
	cn, err := grpc.Dial(g.api.getNextServerAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*20))
	if err != nil {
		return err
	}
	g.conn = cn
	g.client = gnode.NewGNodeServiceClient(g.conn)
	ret, errp := g.client.Ping(g.ctx, &gnode.AntMes{})
	if errp != nil {
		return errp
	}
	g.nodeName = ret.Name
	g.nodeHost = ret.Host
	//g.nbNode = int(ret.NbNode)
	g.nbDuplicate = int(ret.NbDuplicate)
	return nil
}

func (g *gnodeClient) startServerReader() error {
	stream, err := g.client.GetClientStream(g.ctx)
	if err != nil {
		return err
	}
	g.stream = stream
	ack, err2 := g.stream.Recv()
	if err2 != nil {
		g.api.info("Client register EOF\n")
		close(g.recvChan)
		return fmt.Errorf("Client register error: %v\n", err2)
	}
	g.id = ack.FromClient
	g.api.info("Client register: %s\n", g.id)
	go func() {
		for {
			mes, err := g.stream.Recv()
			if err == io.EOF {
				g.api.debug("Server stream EOF\n")
				close(g.recvChan)
				return
			}
			if err != nil {
				g.api.debug("Server stream error: %v\n", err)
				return
			}
			if mes.NoBlocking {
				select {
				case g.recvChan <- mes:
					//fmt.Printf("receive mes noBlocking: %v\n", mes)
				default:
					//fmt.Printf("receive mes noBlocking (wipeout): %v\n", mes)
				}
			} else {
				//fmt.Printf("receive mes Blocking: %v\n", mes)
				g.recvChan <- mes
			}
			g.api.debug("Receive answer: %v\n", mes)
		}
	}()
	return nil
}

func (g *gnodeClient) createMessage(target string, returnAnswer bool, functionName string, args ...string) *gnode.AntMes {
	mes := gnode.NewAntMes(target, returnAnswer, functionName, args...)
	mes.UserName = g.api.userName
	return mes
}

func (g *gnodeClient) createSignedMessage(target string, returnAnswer bool, functionName string, payload []byte, args ...string) (*gnode.AntMes, error) {
	mes := gnode.NewAntMes(target, returnAnswer, functionName, args...)
	mes.UserName = g.api.userName
	mes.Data = payload
	dataToSign := g.getDataToSign(payload, args)
	key, err := g.api.sign(dataToSign)
	if err != nil {
		return nil, fmt.Errorf("Signature error: %v", err)
	}
	mes.Key = key
	return mes, nil
}

func (g *gnodeClient) getDataToSign(payload []byte, args []string) []byte {
	size := len(payload)
	for _, arg := range args {
		size += len(arg)
	}
	dataToSign := make([]byte, size, size)
	nn := g.appendData(dataToSign, 0, payload)
	for _, arg := range args {
		nn = g.appendData(dataToSign, nn, []byte(arg))
	}
	return dataToSign
}

func (g *gnodeClient) appendData(buffer []byte, nn int, item []byte) int {
	for i := 0; i < len(item); i++ {
		buffer[nn+i] = item[i]
	}
	return nn + len(item)
}

func (g *gnodeClient) createSendMessageNoAnswer(target string, functionName string, args ...string) error {
	mes := gnode.NewAntMes(target, false, functionName, args...)
	mes.UserName = g.api.userName
	_, err := g.sendMessage(mes, true)
	return err
}

func (g *gnodeClient) createSendMessage(target string, waitForAnswer bool, functionName string, args ...string) (*gnode.AntMes, error) {
	mes := gnode.NewAntMes(target, true, functionName, args...)
	mes.UserName = g.api.userName
	return g.sendMessage(mes, waitForAnswer)
}

func (g *gnodeClient) sendMessage(mes *gnode.AntMes, wait bool) (*gnode.AntMes, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	mes.FromClient = g.id
	mes.UserName = g.api.userName
	//fmt.Printf("Order: %d size: %d\n", mes.Order, len(mes.Data))
	err := g.stream.Send(mes)
	if err != nil {
		return nil, err
	}
	//g.printf(Info, "Message sent: %v\n", mes)
	if wait {
		ret := <-g.recvChan
		if ret.ErrorMes != "" {
			return nil, fmt.Errorf("%s", ret.ErrorMes)
		}
		return ret, nil
	}
	return nil, nil
}

func (g *gnodeClient) getNextAnswer(timeout int) (*gnode.AntMes, error) {
	if timeout > 0 {
		timer := time.AfterFunc(time.Millisecond*time.Duration(timeout), func() {
			g.recvChan <- &gnode.AntMes{ErrorMes: "timeout"}
		})
		mes := <-g.recvChan
		timer.Stop()
		if mes == nil {
			return nil, fmt.Errorf("Receive nil")
		}
		if mes.ErrorMes != "" {
			return nil, fmt.Errorf("Error: %s", mes.ErrorMes)
		}
		return mes, nil
	}
	mes := <-g.recvChan
	if mes.ErrorMes != "" {
		return nil, fmt.Errorf("Error: %s", mes.ErrorMes)
	}
	return mes, nil
}

func (g *gnodeClient) close() {
	if g.conn != nil {
		g.conn.Close()
	}
}
