package nodecontroller

import (
	"testing"
)

func TestNotifications(t *testing.T) {
	// Create a new NodeController
	nc := &NodeController{
		MsgChan: make(chan NcMsg),
	}

	// Create a new NcMsg
	msg := NcMsg{
		Node:   nil,
		Header: "Test Header",
	}

	// Send the message
	go func() {
		nc.MsgChan <- msg
	}()

	// Receive the message
	receivedMsg := <-nc.MsgChan

	if receivedMsg.Header != msg.Header {
		t.Errorf("Expected %+v, but got %+v", msg, receivedMsg)
	}
}
