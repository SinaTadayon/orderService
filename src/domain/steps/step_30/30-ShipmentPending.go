package shipment_pending_step

import (
	"context"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/domain/steps"
	"gitlab.faza.io/order-project/order-service/infrastructure/global"
	"gitlab.faza.io/order-project/order-service/infrastructure/promise"
	message "gitlab.faza.io/protos/order"
	"time"
)

const (
	stepName string 	= "Shipment_Pending"
	stepIndex int		= 30
	Shipped				= "Shipped"
	SellerShipmentPending = "SellerShipmentPending"
	StockReleased		= "StockReleased"
)

type shipmentPendingStep struct {
	*steps.BaseStepImpl
}

func New(childes, parents []steps.IStep, states ...states.IState) steps.IStep {
	return &shipmentPendingStep{steps.NewBaseStep(stepName, stepIndex, childes, parents, states)}
}

func NewOf(name string, index int, childes, parents []steps.IStep, states ...states.IState) steps.IStep {
	return &shipmentPendingStep{steps.NewBaseStep(name, index, childes, parents, states)}
}

func NewFrom(base *steps.BaseStepImpl) steps.IStep {
	return &shipmentPendingStep{base}
}

func NewValueOf(base *steps.BaseStepImpl, params ...interface{}) steps.IStep {
	panic("implementation required")
}

func (shipmentPending shipmentPendingStep) ProcessMessage(ctx context.Context, request *message.MessageRequest) promise.IPromise {
	panic("implementation required")
}

func (shipmentPending shipmentPendingStep) ProcessOrder(ctx context.Context, order entities.Order, itemsId []string, param interface{}) promise.IPromise {

	if param == nil {
		shipmentPending.UpdateAllOrderStatus(ctx, &order, itemsId, steps.InProgressStatus, false)
		shipmentPending.updateOrderItemsProgress(ctx, &order, itemsId, SellerShipmentPending, true, "", nil, true, steps.InProgressStatus)
		if err := shipmentPending.persistOrder(ctx, &order); err != nil {
			returnChannel := make(chan promise.FutureData, 1)
			defer close(returnChannel)
			returnChannel <- promise.FutureData{Data:nil, Ex:promise.FutureError{Code: promise.InternalError, Reason:"Unknown Error"}}
			return promise.NewPromise(returnChannel, 1, 1)
		}
		returnChannel := make(chan promise.FutureData, 1)
		defer close(returnChannel)
		returnChannel <- promise.FutureData{Data: nil, Ex: nil}
		return promise.NewPromise(returnChannel, 1, 1)
	} else {
		req, ok := param.(*message.RequestSellerOrderAction)
		if ok != true {
			if param == "actionExpired" {
				iPromise := global.Singletons.StockService.BatchStockActions(ctx, order, itemsId, StockReleased)
				futureData := iPromise.Data()
				if futureData == nil {
					if err := shipmentPending.persistOrder(ctx, &order); err != nil {}
					logger.Err("StockService promise channel has been closed, order: %s", order.OrderId)
				} else if futureData.Ex != nil {
					if err := shipmentPending.persistOrder(ctx, &order); err != nil {}
					logger.Err("released stock from stockService failed, error: %s, orderId: %s", futureData.Ex.Error(), order.OrderId)
					returnChannel := make(chan promise.FutureData, 1)
					defer close(returnChannel)
					returnChannel <- promise.FutureData{Data: nil, Ex: promise.FutureError{Code: promise.InternalError, Reason: "Unknown Error"}}
					return promise.NewPromise(returnChannel, 1, 1)
				}

				if len(order.Items) == len(itemsId) {
					shipmentPending.UpdateAllOrderStatus(ctx, &order, itemsId, steps.ClosedStatus, false)
				} else {
					shipmentPending.UpdateAllOrderStatus(ctx, &order, itemsId, steps.InProgressStatus, false)
				}

				shipmentPending.updateOrderItemsProgress(ctx, &order, itemsId, Shipped, false, "Action Expired", nil, false, steps.ClosedStatus)
				if err := shipmentPending.persistOrder(ctx, &order); err != nil {
					returnChannel := make(chan promise.FutureData, 1)
					defer close(returnChannel)
					returnChannel <- promise.FutureData{Data:nil, Ex:promise.FutureError{Code: promise.InternalError, Reason:"Unknown Error"}}
					return promise.NewPromise(returnChannel, 1, 1)
				}

				return shipmentPending.Childes()[1].ProcessOrder(ctx, order, itemsId, nil)
			} else {
				logger.Err("param not a message.RequestSellerOrderAction type , order: %v", order)
				returnChannel := make(chan promise.FutureData, 1)
				defer close(returnChannel)
				returnChannel <- promise.FutureData{Data: nil, Ex: promise.FutureError{Code: promise.InternalError, Reason: "Unknown Error"}}
				return promise.NewPromise(returnChannel, 1, 1)
			}
		}

		if !shipmentPending.validateAction(ctx, &order, itemsId) {
			logger.Err("%s step received invalid action, order: %v, action: %s", shipmentPending.Name(), order, req.Action)
			returnChannel := make(chan promise.FutureData, 1)
			defer close(returnChannel)
			returnChannel <- promise.FutureData{Data: nil, Ex:promise.FutureError{Code:promise.NotAccepted, Reason:"Action Expired"}}
			return promise.NewPromise(returnChannel, 1, 1)
		}

		if req.Data == nil {
			returnChannel := make(chan promise.FutureData, 1)
			defer close(returnChannel)
			returnChannel <- promise.FutureData{Data:nil, Ex:promise.FutureError{Code: promise.BadRequest, Reason:"Reason Data Required"}}
			return promise.NewPromise(returnChannel, 1, 1)
		}

		if req.Action == "success" {
			actionData, ok := req.Data.(*message.RequestSellerOrderAction_Success)
			if ok != true {
				logger.Err("request data not a message.RequestSellerOrderAction_Success type , order: %v", order)
				returnChannel := make(chan promise.FutureData, 1)
				defer close(returnChannel)
				returnChannel <- promise.FutureData{Data: nil, Ex: promise.FutureError{Code: promise.InternalError, Reason: "Unknown Error"}}
				return promise.NewPromise(returnChannel, 1, 1)
			}

			shipmentPending.UpdateAllOrderStatus(ctx, &order, itemsId, steps.InProgressStatus, false)
			shipmentPending.updateOrderItemsProgress(ctx, &order, itemsId, Shipped, true, "", actionData, false, steps.InProgressStatus)
			if err := shipmentPending.persistOrder(ctx, &order); err != nil {
				returnChannel := make(chan promise.FutureData, 1)
				defer close(returnChannel)
				returnChannel <- promise.FutureData{Data:nil, Ex:promise.FutureError{Code: promise.InternalError, Reason:"Unknown Error"}}
				return promise.NewPromise(returnChannel, 1, 1)
			}

			return shipmentPending.Childes()[0].ProcessOrder(ctx, order, itemsId, nil)
		} else if req.Action == "failed" {
			actionData, ok := req.Data.(*message.RequestSellerOrderAction_Failed)
			if ok != true {
				logger.Err("request data not a message.RequestSellerOrderAction_Failed type , order: %v", order)
				returnChannel := make(chan promise.FutureData, 1)
				defer close(returnChannel)
				returnChannel <- promise.FutureData{Data: nil, Ex: promise.FutureError{Code: promise.InternalError, Reason: "Unknown Error"}}
				return promise.NewPromise(returnChannel, 1, 1)
			}

			iPromise := global.Singletons.StockService.BatchStockActions(ctx, order, itemsId, StockReleased)
			futureData := iPromise.Data()
			if futureData == nil {
				if err := shipmentPending.persistOrder(ctx, &order); err != nil {}
				logger.Err("StockService promise channel has been closed, order: %s", order.OrderId)
			} else if futureData.Ex != nil {
				if err := shipmentPending.persistOrder(ctx, &order); err != nil {}
				logger.Err("released stock from stockService failed, error: %s, orderId: %s", futureData.Ex.Error(), order.OrderId)
				returnChannel := make(chan promise.FutureData, 1)
				defer close(returnChannel)
				returnChannel <- promise.FutureData{Data: nil, Ex: promise.FutureError{Code: promise.InternalError, Reason: "Unknown Error"}}
				return promise.NewPromise(returnChannel, 1, 1)
			}

			if len(order.Items) == len(itemsId) {
				shipmentPending.UpdateAllOrderStatus(ctx, &order, itemsId, steps.ClosedStatus, false)
			} else {
				shipmentPending.UpdateAllOrderStatus(ctx, &order, itemsId, steps.InProgressStatus, false)
			}

			shipmentPending.updateOrderItemsProgress(ctx, &order, itemsId, Shipped, false, actionData.Failed.Reason, nil, false, steps.ClosedStatus)
			if err := shipmentPending.persistOrder(ctx, &order); err != nil {
				returnChannel := make(chan promise.FutureData, 1)
				defer close(returnChannel)
				returnChannel <- promise.FutureData{Data:nil, Ex:promise.FutureError{Code: promise.InternalError, Reason:"Unknown Error"}}
				return promise.NewPromise(returnChannel, 1, 1)
			}

			return shipmentPending.Childes()[1].ProcessOrder(ctx, order, itemsId, nil)
		}

		logger.Err("%s step received invalid action, order: %v, action: %s", shipmentPending.Name(), order, req.Action)
		returnChannel := make(chan promise.FutureData, 1)
		defer close(returnChannel)
		returnChannel <- promise.FutureData{Data: nil, Ex: promise.FutureError{Code: promise.InternalError, Reason: "Unknown Error"}}
		return promise.NewPromise(returnChannel, 1, 1)
	}
}

func (shipmentPending shipmentPendingStep) persistOrder(ctx context.Context, order *entities.Order) error {
	_ , err := global.Singletons.OrderRepository.Save(*order)
	if err != nil {
		logger.Err("OrderRepository.Save in %s step failed, order: %v, error: %s", shipmentPending.Name(), order, err.Error())
	}

	return err
}

func (shipmentPending shipmentPendingStep) validateAction(ctx context.Context, order *entities.Order,
	itemsId []string) bool {
	if itemsId != nil && len(itemsId) > 0 {
		for _, id := range itemsId {
			for i := 0; i < len(order.Items); i++ {
				length := len(order.Items[i].Progress.StepsHistory) - 1
				if order.Items[i].ItemId == id && order.Items[i].Progress.StepsHistory[length].Name != shipmentPending.Name() {
					return false
				}
			}
		}
	} else {
		for i := 0; i < len(order.Items); i++ {
			length := len(order.Items[i].Progress.StepsHistory) - 1
			if order.Items[i].Progress.StepsHistory[length].Name != shipmentPending.Name() {
				return false
			}
		}
	}

	return true
}

func (shipmentPending shipmentPendingStep) updateOrderItemsProgress(ctx context.Context, order *entities.Order, itemsId []string,
	action string, result bool, reason string, req *message.RequestSellerOrderAction_Success, isSetExpireTime bool, itemStatus string) {

	findFlag := false
	if itemsId != nil && len(itemsId) > 0 {
		for _, id := range itemsId {
			findFlag = false
			for i := 0; i < len(order.Items); i++ {
				if order.Items[i].ItemId == id {
					findFlag = true
					if req != nil {
						order.Items[i].ShipmentDetails.SellerShipmentDetail = entities.ShipmentDetail{
							TrackingNumber: req.Success.TrackingId,
							ShippingMethod: req.Success.ShipmentMethod,
						}
						break
					} else {
						shipmentPending.doUpdateOrderItemsProgress(ctx, order, i, action, result, reason, isSetExpireTime, itemStatus)
					}
				}
			}
			if !findFlag {
				logger.Err("%s received itemId %s not exist in order, orderId: %v", shipmentPending.Name(), id, order.OrderId)
			}
		}
	} else {
		for i := 0; i < len(order.Items); i++ {
			shipmentPending.doUpdateOrderItemsProgress(ctx, order, i, action, result, reason, isSetExpireTime, itemStatus)
		}
	}
}

func (shipmentPending shipmentPendingStep) doUpdateOrderItemsProgress(ctx context.Context, order *entities.Order, index int,
	actionName string, result bool, reason string, isSetExpireTime bool, itemStatus string) {

	order.Items[index].Status = itemStatus
	order.Items[index].UpdatedAt = time.Now().UTC()

	length := len(order.Items[index].Progress.StepsHistory) - 1

	if order.Items[index].Progress.StepsHistory[length].ActionHistory == nil || len(order.Items[index].Progress.StepsHistory[length].ActionHistory) == 0 {
		order.Items[index].Progress.StepsHistory[length].ActionHistory = make([]entities.Action, 0, 5)
	}

	var action entities.Action
	if isSetExpireTime {
		expiredTime := order.Items[index].UpdatedAt.Add(time.Hour *
			time.Duration(order.Items[index].ShipmentSpec.ReactionTime) +
			time.Minute * time.Duration(0) +
			time.Second * time.Duration(0))

		action = entities.Action{
			Name:      actionName,
			Result:    result,
			Reason:    reason,
			Data:		map[string]interface{}{
				"expiredTime": expiredTime,
			},
			CreatedAt: order.Items[index].UpdatedAt,
		}
	} else {
		action = entities.Action{
			Name:      actionName,
			Result:    result,
			Reason:    reason,
			CreatedAt: order.Items[index].UpdatedAt,
		}
	}

	order.Items[index].Progress.StepsHistory[length].ActionHistory = append(order.Items[index].Progress.StepsHistory[length].ActionHistory, action)
}


//
//import (
//	"gitlab.faza.io/order-project/order-service"
//	OrderService "gitlab.faza.io/protos/order"
//)
//
//func ShipmentPendingEnteredDetail(ppr PaymentPendingRequest, req *OrderService.ShipmentDetailRequest) error {
//	ppr.ShipmentDetail.ShipmentDetail.ShipmentProvider = req.ShipmentProvider
//	ppr.ShipmentDetail.ShipmentDetail.ShipmentTrackingNumber = req.ShipmentTrackingNumber
//	ppr.ShipmentDetail.ShipmentDetail.Description = req.GetDescription()
//	err := main.MoveOrderToNewState("seller", "", main.Shipped, "shipped", ppr)
//	if err != nil {
//		return err
//	}
//	return nil
//}
