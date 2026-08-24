package kafka

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Standard Kafka Topics matching the architecture diagram
const (
	TopicBookingConfirmed  = "BookingConfirmed"
	TopicBookingExpired    = "BookingExpired"
	TopicBookingCancelled  = "BookingCancelled"
	TopicChargingStarted   = "ChargingStarted"
	TopicChargingCompleted = "ChargingCompleted"
	TopicPaymentCreated    = "PaymentCreated"
	TopicPaymentCompleted  = "PaymentCompleted"
	TopicPaymentFailed     = "PaymentFailed"
)

type Event struct {
	ID          string                 `json:"id"`
	Topic       string                 `json:"topic"`
	Source      string                 `json:"source"`       // "Station DB", "Booking DB", "Session DB", "Billing DB"
	AggregateID string                 `json:"aggregate_id"` // e.g. bookingId, sessionId, invoiceId
	Payload     map[string]interface{} `json:"payload"`
	Timestamp   time.Time              `json:"timestamp"`
}

type EventBroker struct {
	mu        sync.RWMutex
	events    []Event
	listeners map[string][]func(event Event)
}

var globalBroker *EventBroker
var once sync.Once

func GetBroker() *EventBroker {
	once.Do(func() {
		globalBroker = &EventBroker{
			events:    make([]Event, 0),
			listeners: make(map[string][]func(event Event)),
		}
	})
	return globalBroker
}

func (b *EventBroker) Subscribe(topic string, handler func(event Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[topic] = append(b.listeners[topic], handler)
	log.Printf("[KAFKA MESSAGE BROKER] Subscribed handler to Topic: [%s]", topic)
}

func (b *EventBroker) Publish(topic, source, aggregateID string, payload map[string]interface{}) Event {
	evt := Event{
		ID:          "evt-" + uuid.New().String()[:8],
		Topic:       topic,
		Source:      source,
		AggregateID: aggregateID,
		Payload:     payload,
		Timestamp:   time.Now(),
	}

	b.mu.Lock()
	b.events = append(b.events, evt)
	if len(b.events) > 500 {
		b.events = b.events[len(b.events)-500:]
	}
	handlers := append([]func(event Event){}, b.listeners[topic]...)
	b.mu.Unlock()

	log.Printf("[KAFKA MESSAGE BROKER] Published Event to Topic: [%s] | Source: %s | AggregateID: %s", topic, source, aggregateID)

	for _, h := range handlers {
		go h(evt)
	}

	return evt
}

func (b *EventBroker) GetEvents() []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	copied := make([]Event, len(b.events))
	copy(copied, b.events)
	return copied
}
