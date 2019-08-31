package main

import "github.com/Shopify/sarama"

func PayToMarketSuccessMessageValidate(message *sarama.ConsumerMessage) (*sarama.ConsumerMessage, error) {
	return message, nil
}

func PayToMarketSuccessAction(message *sarama.ConsumerMessage) error {

	err := PayToMarketSuccessProduce("", []byte{})
	if err != nil {
		return err
	}
	return nil
}

func PayToMarketSuccessProduce(topic string, payload []byte) error {
	return nil
}
