package state_01

import (
	"context"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/app"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	system_action "gitlab.faza.io/order-project/order-service/domain/actions/system"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
	"gitlab.faza.io/order-project/order-service/infrastructure/future"
	stock_service "gitlab.faza.io/order-project/order-service/infrastructure/services/stock"
	"gitlab.faza.io/order-project/order-service/infrastructure/utils"
	"time"
)

const (
	stepName  string = "New_Order"
	stepIndex int    = 1
)

type newOrderState struct {
	*states.BaseStateImpl
}

func New(childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &newOrderState{states.NewBaseStep(stepName, stepIndex, childes, parents, actionStateMap)}
}

func NewOf(name string, index int, childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &newOrderState{states.NewBaseStep(name, index, childes, parents, actionStateMap)}
}

func NewFrom(base *states.BaseStateImpl) states.IState {
	return &newOrderState{base}
}

func NewValueOf(base *states.BaseStateImpl, params ...interface{}) states.IState {
	panic("implementation required")
}

func (state newOrderState) Process(ctx context.Context, iFrame frame.IFrame) {
	//var errStr string
	//logger.Audit("New Order Received . . .")

	order := iFrame.Header().Value(string(frame.HeaderOrder)).(*entities.Order)
	action := &entities.Action{
		Name:      state.Actions()[0].ActionEnum().ActionName(),
		Type:      "",
		UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
		UTP:       state.Actions()[0].ActionType().ActionName(),
		Perm:      "",
		Priv:      "",
		Policy:    "",
		Result:    string(states.ActionSuccess),
		Reasons:   nil,
		Data:      nil,
		CreatedAt: time.Now().UTC(),
		Extended:  nil,
	}

	state.UpdateOrderAllStatus(ctx, order, states.OrderNewStatus, states.PackageNewStatus, action)
	newOrder, err := app.Globals.OrderRepository.Save(ctx, *order)
	if err != nil {
		logger.Err("OrderRepository.Save in %s state failed, order: %v, error: %s", state.Name(), order, err)
		state.releasedStock(ctx, newOrder)
		future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
			SetError(future.ErrorCode(err.Code()), err.Message(), err.Reason()).
			Send()

	} else {
		newFrame := frame.Factory().
			SetFuture(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
			SetOrderId(newOrder.OrderId).SetBody(newOrder).Build()

		state.StatesMap()[state.Actions()[0]].Process(ctx, newFrame)
	}
}

func (state newOrderState) releasedStock(ctx context.Context, order *entities.Order) {

	for i := 0; i < len(order.Packages); i++ {
		for j := 0; j < len(order.Packages[i].Subpackages); j++ {
			result := true
			stockActionDataList := make([]entities.StockActionData, 0, 32)
			for z := 0; z < len(order.Packages[i].Subpackages[j].Items); z++ {
				item := order.Packages[i].Subpackages[j].Items[z]
				requestStock := stock_service.RequestStock{
					InventoryId: item.InventoryId,
					Count:       int(item.Quantity),
				}

				iFuture := app.Globals.StockService.SingleStockAction(ctx, requestStock, order.OrderId,
					system_action.New(system_action.StockRelease))

				futureData := iFuture.Get()
				if futureData.Error() != nil {
					result = false
					if futureData.Data() != nil {
						response := futureData.Data().(stock_service.ResponseStock)
						actionData := entities.StockActionData{
							InventoryId: response.InventoryId,
							Quantity:    response.Count,
							Result:      response.Result,
						}
						stockActionDataList = append(stockActionDataList, actionData)
						logger.Err("releasedStock() => Released stock from stockService failed, state: %s, orderId: %d, response: %v, error: %s", state.Name(), order.OrderId, response, futureData.Error())
					} else {
						actionData := entities.StockActionData{
							InventoryId: requestStock.InventoryId,
							Quantity:    requestStock.Count,
							Result:      false,
						}
						stockActionDataList = append(stockActionDataList, actionData)
						logger.Err("releasedStock() => Released stock from stockService failed, state: %s, orderId: %d, error: %s", state.Name(), order.OrderId, futureData.Error())
					}
				} else {
					response := futureData.Data().(stock_service.ResponseStock)
					actionData := entities.StockActionData{
						InventoryId: response.InventoryId,
						Quantity:    response.Count,
						Result:      response.Result,
					}
					stockActionDataList = append(stockActionDataList, actionData)
					logger.Audit("Release stock success, state: %s, orderId: %d", state.Name(), order.OrderId)
				}
			}
			var stockAction *entities.Action
			if !result {
				stockAction = &entities.Action{
					Name:      system_action.StockRelease.ActionName(),
					Type:      "",
					UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
					UTP:       actions.System.ActionName(),
					Perm:      "",
					Priv:      "",
					Policy:    "",
					Result:    string(states.ActionFail),
					Reasons:   nil,
					Data:      map[string]interface{}{"stockActionData": stockActionDataList},
					CreatedAt: time.Now().UTC(),
					Extended:  nil,
				}
			} else {
				stockAction = &entities.Action{
					Name:      system_action.StockRelease.ActionName(),
					Type:      "",
					UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
					UTP:       actions.System.ActionName(),
					Perm:      "",
					Priv:      "",
					Policy:    "",
					Result:    string(states.ActionSuccess),
					Reasons:   nil,
					Data:      map[string]interface{}{"stockActionData": stockActionDataList},
					CreatedAt: time.Now().UTC(),
					Extended:  nil,
				}
			}

			state.UpdateSubPackage(ctx, order.Packages[i].Subpackages[j], stockAction)
		}
	}
}
