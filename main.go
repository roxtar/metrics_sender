package main
import (
"flag"
"github.com/cloudfoundry/sonde-go/events"
"github.com/gogo/protobuf/proto"
	"log"
	"fmt"
"net"
	"time"
)

const DEFAULT_METRON_PORT = 3457
const DEFAULT_ORIGIN = "message_sender"

func main() {
	messageType := flag.String("type", "log", `Type of the message. Can be one of "log", "value", "counter"`)
	logMessage := flag.String("message", "log message", "Message bytes for a log message")
	appID := flag.String("appid", "default appid", "Application ID for a message")
	value := flag.Float64("value", 0.0, "Value for a value metric")
	unit := flag.String("unit", "unit", "Unit for a value metric")
	name := flag.String("name", "metricName", "Name for a value metric or a counter event")
	delta := flag.Uint64("delta", 0, "Delta for the counter event")

	flag.Parse()

	switch *messageType {
	case "log":
		sendLog(*logMessage, *appID)
	case "value":
		sendValue(*name, *value, *unit)
	case "counter":
		sendCounter(*name, *delta)
	default:
		log.Fatalf("Unknown or unsupported event type: %s", *messageType)
	}

}

func sendLog(message, appID string) {
	envelope := &events.Envelope {
		EventType: events.Envelope_LogMessage.Enum(),
		Timestamp: proto.Int64(time.Now().UnixNano()),
		Origin: proto.String(DEFAULT_ORIGIN),
		LogMessage: &events.LogMessage {
			MessageType: events.LogMessage_OUT.Enum(),
			Message: []byte(message),
			Timestamp: proto.Int64(time.Now().UnixNano()),
			AppId: proto.String(appID),
		},
	}
	sendEnvelope(envelope)
}

func sendValue(name string, value float64, unit string) {
	envelope := &events.Envelope {
		EventType: events.Envelope_ValueMetric.Enum(),
		Timestamp: proto.Int64(time.Now().UnixNano()),
		Origin: proto.String(DEFAULT_ORIGIN),
		ValueMetric: &events.ValueMetric {
			Name: proto.String(name),
			Value: proto.Float64(value),
			Unit: proto.String(unit),
		},
	}
	sendEnvelope(envelope)
}

func sendCounter(name string, delta uint64) {
	envelope := &events.Envelope {
		EventType: events.Envelope_CounterEvent.Enum(),
		Timestamp: proto.Int64(time.Now().UnixNano()),
		Origin: proto.String(DEFAULT_ORIGIN),
		CounterEvent: &events.CounterEvent {
			Name: proto.String(name),
			Delta: proto.Uint64(delta),
		},
	}
	sendEnvelope(envelope)
}

func sendEnvelope(envelope *events.Envelope) {
	bytes, err := proto.Marshal(envelope)
	if err != nil {
		log.Fatalf("Error marshalling envelope %s", err.Error())
	}
	address := fmt.Sprintf("127.0.0.1:%d", DEFAULT_METRON_PORT)
	conn, err := net.Dial("udp4", address)
	if err != nil {
		log.Fatalf("Error dialing address [%s]: %s", address, err.Error())
	}
	_, err = conn.Write(bytes)
	if err != nil {
		log.Fatalf("Error writing bytes: %s", err.Error())
	}
}
