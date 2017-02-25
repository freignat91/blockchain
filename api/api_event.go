package abcapi

import (
	"fmt"
	"github.com/freignat91/agrid/server/gnode"
	"reflect"
)

type TransferEvent struct {
	EventType  string
	TransferId string
	EventDate  string
	UserName   string
	FileType   string
	FileName   string
	State      string
	Metadata   []string
}

func (api *BchainAPI) FileSetTransferEventCallback(fileType string, callbackFunction interface{}) error {
	f := reflect.ValueOf(callbackFunction)
	if f.Type().String() != "func(*agridapi.TransferEvent) error" {
		return fmt.Errorf("the callback should be a function type: func(*AgridAPI.TransferEvent) error")
	}
	client, err := api.getClient()
	if err != nil {
		return err
	}
	client.sendMessage(&gnode.AntMes{
		Function:  "setEventListener",
		UserName:  api.userName,
		UserToken: api.userToken,
		Args:      []string{"TransferEvent", fileType},
	}, false)
	for {
		mes, err := client.getNextAnswer(0)
		if err != nil {
			api.info("received error%v\n", mes)
			return err
		}
		if mes.Function == "sendBackEvent" && mes.Args[0] == "TransferEvent" && (fileType == "" || mes.Args[4] == fileType) {
			event := &TransferEvent{
				EventType:  mes.Args[0],
				EventDate:  mes.Args[1],
				TransferId: mes.Args[2],
				State:      mes.Args[3],
				UserName:   mes.UserName,
				FileName:   mes.TargetedPath,
				FileType:   mes.Args[4],
				Metadata:   mes.Args[5:],
			}
			ret := f.Call([]reflect.Value{reflect.ValueOf(event)})
			if ret[0].Interface() != nil {
				return ret[0].Interface().(error)
			}
		}
	}
}
