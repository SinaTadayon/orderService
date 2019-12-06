package state_80

import (
	"context"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
)

const (
	stepName  string = "Pay_To_Buyer"
	stepIndex int    = 80
	//Canceled         = "CANCELED"
)

type payToBuyerState struct {
	*states.BaseStateImpl
}

func New(childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &payToBuyerState{states.NewBaseStep(stepName, stepIndex, childes, parents, actionStateMap)}
}

func NewOf(name string, index int, childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &payToBuyerState{states.NewBaseStep(name, index, childes, parents, actionStateMap)}
}

func NewFrom(base *states.BaseStateImpl) states.IState {
	return &payToBuyerState{base}
}

func NewValueOf(base *states.BaseStateImpl, params ...interface{}) states.IState {
	panic("implementation required")
}

func (state payToBuyerState) Process(ctx context.Context, iFrame frame.IFrame) {
	panic("implementation required")
}

//func (payToBuyer payToBuyerState) ProcessOrder(ctx context.Context, order entities.Order, itemsId []uint64, param interface{}) future.IFuture {
//
//	logger.Audit("Pay to Buyer step, orderId: %d", order.OrderId)
//
//	if len(order.Items) == len(itemsId) {
//		payToBuyer.UpdateAllOrderStatus(ctx, &order, itemsId, states.OrderClosedStatus, false)
//	} else {
//		payToBuyer.UpdateAllOrderStatus(ctx, &order, itemsId, states.OrderInProgressStatus, false)
//	}
//
//	payToBuyer.updateOrderItemsProgress(ctx, &order, itemsId, Canceled, true, states.OrderClosedStatus)
//	if err := payToBuyer.persistOrder(ctx, &order); err != nil {
//		returnChannel := make(chan future.IDataFuture, 1)
//		defer close(returnChannel)
//		returnChannel <- future.IDataFuture{Data: nil, Ex: future.FutureError{Code: future.InternalError, Reason: "Unknown Error"}}
//		return future.NewFuture(returnChannel, 1, 1)
//	}
//	returnChannel := make(chan future.IDataFuture, 1)
//	defer close(returnChannel)
//	returnChannel <- future.IDataFuture{Data: nil, Ex: nil}
//	return future.NewFuture(returnChannel, 1, 1)
//}
//
//func (payToBuyer payToBuyerState) persistOrder(ctx context.Context, order *entities.Order) error {
//	_, err := global.Singletons.OrderRepository.Save(*order)
//	if err != nil {
//		logger.Err("OrderRepository.Save in %s step failed, order: %v, error: %s", payToBuyer.Name(), order, err.Error())
//	}
//	return err
//}
//
//func (payToBuyer payToBuyerState) updateOrderItemsProgress(ctx context.Context, order *entities.Order, itemsId []uint64, action string, result bool, itemStatus string) {
//
//	findFlag := false
//	if itemsId != nil && len(itemsId) > 0 {
//		for _, id := range itemsId {
//			findFlag = false
//			for i := 0; i < len(order.Items); i++ {
//				if order.Items[i].ItemId == id {
//					payToBuyer.doUpdateOrderItemsProgress(ctx, order, i, action, result, itemStatus)
//					findFlag = true
//					break
//				}
//			}
//
//			if findFlag == false {
//				logger.Err("%s received itemId %d not exist in order, orderId: %d", payToBuyer.Name(), id, order.OrderId)
//			}
//		}
//	} else {
//		for i := 0; i < len(order.Items); i++ {
//			payToBuyer.doUpdateOrderItemsProgress(ctx, order, i, action, result, itemStatus)
//		}
//	}
//}
//
//func (payToBuyer payToBuyerState) doUpdateOrderItemsProgress(ctx context.Context, order *entities.Order, index int,
//	actionName string, result bool, itemStatus string) {
//
//	order.Items[index].Status = itemStatus
//	order.Items[index].UpdatedAt = time.Now().UTC()
//
//	length := len(order.Items[index].Progress.StepsHistory) - 1
//
//	if order.Items[index].Progress.StepsHistory[length].ActionHistory == nil || len(order.Items[index].Progress.StepsHistory[length].ActionHistory) == 0 {
//		order.Items[index].Progress.StepsHistory[length].ActionHistory = make([]entities.Action, 0, 5)
//	}
//
//	action := entities.Action{
//		Name:      actionName,
//		Result:    result,
//		CreatedAt: order.Items[index].UpdatedAt,
//	}
//
//	order.Items[index].Progress.StepsHistory[length].ActionHistory = append(order.Items[index].Progress.StepsHistory[length].ActionHistory, action)
//}
//
////
////import (
////	"github.com/Shopify/sarama"
////	"gitlab.faza.io/order-project/order-service"
////)
////
////// TODO: must be implemented
////func PayToBuyerMessageValidate(message *sarama.ConsumerMessage) (*sarama.ConsumerMessage, error) {
////	mess, err := main.CheckOrderKafkaAndMongoStatus(message, main.PayToBuyer)
////	if err != nil {
////		return mess, err
////	}
////	return message, nil
////}
////
////func PayToBuyerAction(message *sarama.ConsumerMessage) error {
////
////	err := PayToBuyerProduce("", []byte{})
////	if err != nil {
////		return err
////	}
////	return nil
////}
////
////func PayToBuyerProduce(topic string, payload []byte) error {
////	return nil
////}
