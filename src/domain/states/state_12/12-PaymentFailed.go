package state_12

import (
	"context"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/app"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	system_action "gitlab.faza.io/order-project/order-service/domain/actions/system"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
	"gitlab.faza.io/order-project/order-service/infrastructure/utils"
	"time"
)

const (
	stepName  string = "Payment_Failed"
	stepIndex int    = 12
)

type paymentFailedState struct {
	*states.BaseStateImpl
}

func New(childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &paymentFailedState{states.NewBaseStep(stepName, stepIndex, childes, parents, actionStateMap)}
}

func NewOf(name string, index int, childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &paymentFailedState{states.NewBaseStep(name, index, childes, parents, actionStateMap)}
}

func NewFrom(base *states.BaseStateImpl) states.IState {
	return &paymentFailedState{base}
}

func NewValueOf(base *states.BaseStateImpl, params ...interface{}) states.IState {
	panic("implementation required")
}

func (state paymentFailedState) Process(ctx context.Context, iFrame frame.IFrame) {

	if iFrame.Header().KeyExists(string(frame.HeaderOrderId)) && iFrame.Body().Content() != nil {
		order, ok := iFrame.Body().Content().(*entities.Order)
		if !ok {
			logger.Err("iFrame.Body().Content() not a order, orderId: %d, %s state ",
				iFrame.Header().Value(string(frame.HeaderOrderId)), state.Name())
			return
		}

		var stockAction *entities.Action
		if err := state.releasedStock(ctx, order); err != nil {
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
				Data:      nil,
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
				Data:      nil,
				CreatedAt: time.Now().UTC(),
				Extended:  nil,
			}
		}

		state.UpdateOrderAllStatus(ctx, order, states.OrderClosedStatus, states.PackageClosedStatus, stockAction)
		_, err := app.Globals.OrderRepository.Save(ctx, *order)
		if err != nil {
			logger.Err("OrderRepository.Save in %s state failed, orderId: %d, error: %s", state.Name(), order.OrderId, err.Error())
		}
		logger.Audit("Order System Failed, orderId: %d", order.OrderId)
	} else {
		logger.Err("HeaderOrderId of iFrame.Header not found and content of iFrame.Body() not set, state: %s iframe: %v", state.Name(), iFrame)
	}
}

func (state paymentFailedState) releasedStock(ctx context.Context, order *entities.Order) error {

	var inventories = make(map[string]int, 32)
	for i := 0; i < len(order.Packages); i++ {
		for j := 0; j < len(order.Packages[i].Subpackages); j++ {
			for z := 0; z < len(order.Packages[i].Subpackages[j].Items); z++ {
				item := order.Packages[i].Subpackages[j].Items[z]
				inventories[item.InventoryId] = int(item.Quantity)
			}
		}
	}

	iFuture := app.Globals.StockService.BatchStockActions(ctx, inventories,
		system_action.New(system_action.StockRelease))
	futureData := iFuture.Get()
	if futureData.Error() != nil {
		logger.Err("Reserved stock from stockService failed, state: %s, orderId: %d, error: %s", state.Name(), order.OrderId, futureData.Error())
		return futureData.Error().Reason()
	}

	logger.Audit("Release stock success, state: %s, orderId: %d", state.Name(), order.OrderId)
	return nil
}
