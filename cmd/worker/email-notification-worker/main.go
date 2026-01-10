package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	messaging "github.com/RehanAthallahAzhar/tokohobby-messaging-go"
	orderMsg "github.com/RehanAthallahAzhar/tokohobby-orders/internal/messaging"
)

func main() {
	log.Println("🚀 Starting Order Email Notification Worker...")

	// Initialize RabbitMQ
	rmqConfig := messaging.DefaultConfig()
	rmq, err := messaging.NewRabbitMQ(rmqConfig)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Setup queue
	if err := messaging.SetupOrderEmailQueue(rmq); err != nil {
		log.Fatalf("Failed to setup queue: %v", err)
	}

	emailService := NewEmailService()

	// Message handler
	handler := func(ctx context.Context, body []byte) error {
		var event orderMsg.OrderStatusChangedEvent

		if err := messaging.UnmarshalMessage(body, &event); err != nil {
			return fmt.Errorf("failed to unmarshal: %w", err)
		}

		log.Printf("📧 Processing email notification for order %s: %s → %s",
			event.OrderID, event.OldStatus, event.NewStatus)

		// Send email based on status
		return emailService.SendOrderStatusEmail(event)
	}

	// Create consumer - consume ALL status changes
	consumerOpts := messaging.ConsumerOptions{
		QueueName:   "order.email.notifications",
		WorkerCount: 5, // 5 concurrent workers
		AutoAck:     false,
	}
	consumer := messaging.NewConsumer(rmq, consumerOpts, handler)

	// Start consuming
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down worker...")
	cancel()
}

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendOrderStatusEmail(event orderMsg.OrderStatusChangedEvent) error {
	// Different email template based on status
	templates := map[orderMsg.OrderStatus]string{
		orderMsg.OrderStatusPaid:      "✅ Payment received!",
		orderMsg.OrderStatusShipped:   "📦 Your order has been shipped!",
		orderMsg.OrderStatusDelivered: "🎉 Order delivered successfully!",
		orderMsg.OrderStatusCancelled: "❌ Order cancelled",
	}

	template, exists := templates[event.NewStatus]
	if !exists {
		template = "Order status updated"
	}

	log.Printf("✅ [MOCK] Email sent to %s: %s (Order: %s)",
		event.Email, template, event.OrderID)
	return nil
}
