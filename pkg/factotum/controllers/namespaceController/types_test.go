package namespacecontroller

import (
	"testing"
)

func TestNotifications(t *testing.T) {
	// Create a new NamespaceController
	c := &NamespaceController{
		MsgChan: make(chan Msg),
	}

	// Create a new NcMsg
	msg := Msg{
		Namespace: nil,
		Header:    "Test Header",
	}

	// Send the message
	go func() {
		c.MsgChan <- msg
	}()

	// Receive the message
	receivedMsg := <-c.MsgChan

	if receivedMsg.Header != msg.Header {
		t.Errorf("Expected %+v, but got %+v", msg, receivedMsg)
	}
}
