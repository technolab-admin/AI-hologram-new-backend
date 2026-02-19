
package unit_test


import (
    "testing"
	"time"
	"fmt"
    "reflect"
    "github.com/gorilla/websocket"
	"AI-HOLOGRAM-NEW-BACKEND/internal/websockets"
)


////////////////
// Mockfunctions
////////////////


type mockConn struct {
    written [][]byte
    read    chan []byte
    closed  bool
}

func (m *mockConn) ReadMessage() (int, []byte, error) {
    msg, ok := <-m.read
    if !ok {
        return 0, nil, fmt.Errorf("Closed")
    }
    return websocket.TextMessage, msg, nil
}

func (m *mockConn) WriteMessage(_ int, msg []byte) error {
    m.written = append(m.written, msg)
    return nil
}

func (m *mockConn) Close() error {
    m.closed = true
    return nil
}


var exampleGoodMsgRaw = []byte(
    `{
        "from": "TEST-A",
        "target": "TEST-B",
        "event": "request_new_model",
        "data": "A golden knight"
    }`)

var exampleGoodMsgJSON = map[string]string{
        "from": "TEST-A",
        "target": "TEST-B",
        "event": "request_new_model",
        "data": "A golden knight",
	}


////////
// Tests
////////


func TestMessageRouting(t *testing.T) {
    s := websockets.NewServer(":0")

    c1Conn := &mockConn{read: make(chan []byte, 1)}
    c2Conn := &mockConn{read: make(chan []byte, 1)}

    c1 := websockets.NewClient("TEST-A", c1Conn)
    c2 := websockets.NewClient("TEST-B", c2Conn)

    s.AddClient(c1)
    s.AddClient(c2)

    go s.Listener("TEST-A")
    go s.MessageHandler("TEST-A")

    c1Conn.read <- exampleGoodMsgRaw

    time.Sleep(20 * time.Millisecond)

    if len(c2Conn.written) != 1 {
        t.Fatalf("Error: Expected message to be forwarded")
    }
}


func TestMarshalSuccess(t *testing.T) {

	raw, err := websockets.MarshalMessage(exampleGoodMsgJSON)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(raw) == 0 {
		t.Fatalf("Error: JSON output is empty")
	}
}


func TestUnmarshalSuccess(t *testing.T) {

	msg, err := websockets.UnmarshalMessage(exampleGoodMsgRaw)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if !reflect.DeepEqual(msg, exampleGoodMsgJSON) {
		t.Fatalf("Error: expected %v, got %v", exampleGoodMsgJSON, msg)
	}
}


func TestUnmarshalFail(t *testing.T) {

	badRaw1 := []byte(`{ invalid json }`)
	badRaw2 := []byte(`{ "from": 123 }`)
	badRaw3 := []byte(`{ "from": true }`)

	_, err := websockets.UnmarshalMessage(badRaw1)
	if err == nil {t.Fatalf("Error, expected error is nil")}

    _, err = websockets.UnmarshalMessage(badRaw2)
    if err == nil {t.Fatalf("Error, expected error is nil")}

    _, err = websockets.UnmarshalMessage(badRaw3)
    if err == nil {t.Fatalf("Error, expected error is nil")}
}


func TestValidationSuccess(t *testing.T) {

    err := websockets.HasAllFields(exampleGoodMsgJSON)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    err = websockets.HasNoExtraFields(exampleGoodMsgJSON)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    
    isKnown := websockets.IsKnownClient("TEST-A")
    if !isKnown {
        t.Fatalf("Error: client is unknown")
    }

    err = websockets.VallidateClients(exampleGoodMsgJSON)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
}


func TestValidationFail(t *testing.T) {

    msgMissingField := map[string]string{
        "from": "TEST-A",
        "target": "TEST-B",
        "event": "request_new_model",
	}

    msgEmptyField := map[string]string{
        "from": "TEST-A",
        "target": "TEST-B",
        "event": "request_new_model",
        "data": "",
	}

    msgExtraField := map[string]string{
        "from": "TEST-A",
        "target": "TEST-B",
        "event": "request_new_model",
        "data": "A golden knight",
        "fail": "fail",
	}

    msgTypo := map[string]string{
        "from": "TEST-A",
        "Target": "TEST-B",
        "event": "request_new_model",
        "data": "A golden knight",
	}
    
    msgBadFrom := map[string]string{
        "from": "TEST-C",
        "target": "TEST-B",
        "event": "request_new_model",
        "data": "A golden knight",
	}

    msgBadTarget := map[string]string{
        "from": "TEST-A",
        "target": "TEST-C",
        "event": "request_new_model",
        "data": "A golden knight",
	}

    err := websockets.HasAllFields(msgMissingField)
    if err == nil {t.Fatalf("Error: expected missing-field-error, got nil")}
    
    err = websockets.HasAllFields(msgEmptyField)
    if err == nil {t.Fatalf("Error: expected missing-field-error, got nil")}
    
    err = websockets.HasNoExtraFields(msgExtraField)
    if err == nil {t.Fatalf("Error: expected extra-field-error, got nil")}
    
    err = websockets.HasAllFields(msgTypo)
    if err == nil {t.Fatalf("Error: expected missing-field-error by typo, got nil")}

    err = websockets.VallidateClients(msgBadFrom)
    if err == nil {t.Fatalf("Error: expected unknown-from-error, got nil")}

    err = websockets.VallidateClients(msgBadTarget)
    if err == nil {t.Fatalf("Error: expected unknown-target-error, got nil")}
}


