package main

import "github.com/Shopify/sarama"

func ShipmentPendingMessageValidate(message *sarama.ConsumerMessage) (*sarama.ConsumerMessage, error) {
	return message, nil
}

func ShipmentPendingAction(message *sarama.ConsumerMessage) error {

	err := ShipmentPendingProduce("", []byte{})
	if err != nil {
		return err
	}
	return nil
}

func ShipmentPendingProduce(topic string, payload []byte) error {
	return nil
}
