package return_shipment_delivery_delayed_step

//
//import (
//	"context"
//	"gitlab.faza.io/order-project/order-service"
//	"log"
//	"os"
//	"os/signal"
//	"strings"
//	"sync"
//	"syscall"
//
//	"gitlab.faza.io/go-framework/logger"
//
//	"github.com/Shopify/sarama"
//)
//
//// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
//func (c *ReturnShipmentDeliveryDelayedConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//	for m := range claim.Messages() {
//		logger.Audit("ReturnShipmentDeliveryDelayed: value = %s, OffSet = %v, topic = %s",
//			string(m.Value), m.Offset, m.Topic)
//		// Validate message
//		message, err := ReturnShipmentDeliveryDelayedMessageValidate(m)
//		if err != nil {
//			logger.Err("ReturnShipmentDeliveryDelayed validation failed: %s", err)
//			continue
//		}
//		// Message action
//		err = ReturnShipmentDeliveryDelayedAction(message)
//		if err != nil {
//			logger.Err("ReturnShipmentDeliveryDelayed action failed: %s", err)
//			continue
//		}
//		session.MarkMessage(message, "")
//	}
//	return nil
//}
//
//func startReturnShipmentDeliveryDelayed(Version string, topics string) {
//	logger.Audit("starting ReturnShipmentDeliveryDelayed consumers...")
//	if main.App.Config.Kafka.Brokers == "" {
//		log.Fatal("Cant start ReturnShipmentDeliveryDelayed consumer, No brokers defined")
//	}
//	if main.App.Config.Kafka.ConsumerTopic == "" {
//		log.Fatal("Cant start ReturnShipmentDeliveryDelayed consumer, No topic defined")
//	}
//	if main.App.Config.Kafka.Version == "" {
//		log.Fatal("Cant start ReturnShipmentDeliveryDelayed consumer, No kafka version defined")
//	}
//
//	brokers := strings.Split(main.App.Config.Kafka.Brokers, ",")
//	version, err := sarama.ParseKafkaVersion(Version)
//	if err != nil {
//		log.Panicf("Error parsing Kafka version: %v", err)
//	}
//
//	config := sarama.NewConfig()
//	config.Version = version
//	//Config.Consumer.Retry.Backoff
//	config.Consumer.Offsets.Initial = sarama.OffsetOldest
//
//	consumer := ReturnShipmentDeliveryDelayedConsumer{
//		ready: make(chan bool, 0),
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	client, err := sarama.NewConsumerGroup(brokers, main.App.Config.Kafka.ConsumerGroup, config)
//	if err != nil {
//		log.Panicf("Error creating ReturnShipmentDeliveryDelayed consumer group client: %v", err)
//	}
//
//	wg := &sync.WaitGroup{}
//	go func() {
//		wg.Add(1)
//		defer wg.Done()
//		for {
//			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
//				log.Panicf("Error from ReturnShipmentDeliveryDelayed consumer: %v", err)
//			}
//			if ctx.Err() != nil {
//				return
//			}
//			consumer.ready = make(chan bool, 0)
//		}
//	}()
//
//	<-consumer.ready // Await till the consumer has been set up
//	log.Println("ReturnShipmentDeliveryDelayed consumer up and running!...")
//
//	sigterm := make(chan os.Signal, 1)
//	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
//	select {
//	case <-ctx.Done():
//		log.Println("terminating: context cancelled")
//	case <-sigterm:
//		log.Println("terminating: via signal")
//	}
//	cancel()
//	wg.Wait()
//	if err = client.Close(); err != nil {
//		log.Panicf("Error closing client: %v", err)
//	}
//}
//
//// Consumer represents a Sarama consumer group consumer
//type ReturnShipmentDeliveryDelayedConsumer struct {
//	ready chan bool
//}
//
//// Setup is run at the beginning of a new session, before ConsumeClaim
//func (c *ReturnShipmentDeliveryDelayedConsumer) Setup(sarama.ConsumerGroupSession) error {
//	close(c.ready)
//	return nil
//}
//
//// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
//func (c *ReturnShipmentDeliveryDelayedConsumer) Cleanup(sarama.ConsumerGroupSession) error {
//	return nil
//}
