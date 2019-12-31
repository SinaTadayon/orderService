package state_80

import (
	"bytes"
	"context"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/app"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	operator_action "gitlab.faza.io/order-project/order-service/domain/actions/operator"
	seller_action "gitlab.faza.io/order-project/order-service/domain/actions/seller"
	system_action "gitlab.faza.io/order-project/order-service/domain/actions/system"
	"gitlab.faza.io/order-project/order-service/domain/events"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
	notify_service "gitlab.faza.io/order-project/order-service/infrastructure/services/notification"
	"gitlab.faza.io/order-project/order-service/infrastructure/utils"
	"text/template"
	"time"
)

const (
	stepName  string = "Pay_To_Buyer"
	stepIndex int    = 80
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
	if iFrame.Header().KeyExists(string(frame.HeaderSubpackages)) {
		subpackages, ok := iFrame.Header().Value(string(frame.HeaderSubpackages)).([]*entities.Subpackage)
		if !ok {
			logger.Err("iFrame.Header() not a subpackages, frame: %v, %s state ", iFrame, state.Name())
			return
		}

		sids, ok := iFrame.Header().Value(string(frame.HeaderSIds)).([]uint64)
		if !ok {
			logger.Err("iFrame.Header() not a sids, frame: %v, %s state ", iFrame, state.Name())
			return
		}

		event, ok := iFrame.Header().Value(string(frame.HeaderEvent)).(events.IEvent)
		if !ok {
			logger.Err("Process() => received frame doesn't have a event, state: %s, frame: %v", state.String(), iFrame)
			return
		}

		pkgItem, ok := iFrame.Body().Content().(*entities.PackageItem)
		if !ok {
			logger.Err("Process() => iFrame.Body().Content() is nil, orderId: %d, pid: %d, sid: %d, %s state ",
				subpackages[0].OrderId, subpackages[0].PId, subpackages[0].SId, state.Name())
			return
		}

		var buyerNotificationAction *entities.Action = nil
		if event.Action().ActionEnum() == seller_action.Accept ||
			event.Action().ActionEnum() == operator_action.Accept {

			var message string
			if event.Action().ActionEnum() == seller_action.Accept {
				message = app.Globals.SMSTemplate.OrderNotifyBuyerReturnDeliveredToPayToBuyerState
			} else if event.Action().ActionEnum() == operator_action.Accept {
				message = app.Globals.SMSTemplate.OrderNotifyBuyerReturnRejectedToPayToBuyerState
			}

			var templateData struct {
				OrderId  uint64
				ShopName string
			}

			templateData.OrderId = pkgItem.OrderId
			templateData.ShopName = pkgItem.ShopName

			smsTemplate, err := template.New("SMS").Parse(message)
			if err != nil {
				logger.Err("Process() => smsTemplate.Parse failed, state: %s, orderId: %d, message: %s, err: %s",
					state.Name(), pkgItem.OrderId, message, err)
			} else {
				var buf bytes.Buffer
				err = smsTemplate.Execute(&buf, templateData)
				if err != nil {
					logger.Err("Process() => smsTemplate.Execute failed, state: %s, orderId: %d, message: %s, err: %s",
						state.Name(), pkgItem.OrderId, message, err)
				} else {
					buyerNotify := notify_service.SMSRequest{
						Phone: pkgItem.ShippingAddress.Mobile,
						Body:  buf.String(),
					}

					buyerFutureData := app.Globals.NotifyService.NotifyBySMS(ctx, buyerNotify).Get()
					if buyerFutureData.Error() != nil {
						logger.Err("Process() => NotifyService.NotifyBySMS failed, request: %v, state: %s, orderId: %d, pid: %d, sids: %v, error: %s",
							buyerNotify, state.Name(), pkgItem.OrderId, pkgItem.PId, sids, buyerFutureData.Error().Reason())
						buyerNotificationAction = &entities.Action{
							Name:      system_action.BuyerNotification.ActionName(),
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
						logger.Audit("Process() => NotifyService.NotifyBySMS success, state: %s, orderId: %d, pid: %d, sids: %v",
							state.Name(), pkgItem.OrderId, pkgItem.PId, sids)
						buyerNotificationAction = &entities.Action{
							Name:      system_action.BuyerNotification.ActionName(),
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
				}
			}
		}

		var releaseStockAction *entities.Action
		if err := state.releasedStock(ctx, subpackages); err != nil {
			releaseStockAction = &entities.Action{
				Name:      system_action.StockRelease.ActionName(),
				Type:      "",
				UId:       0,
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
			releaseStockAction = &entities.Action{
				Name:      system_action.StockRelease.ActionName(),
				Type:      "",
				UId:       0,
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

		for _, subpackage := range subpackages {
			if buyerNotificationAction != nil {
				state.UpdateSubPackage(ctx, subpackage, buyerNotificationAction)
			}
			state.UpdateSubPackage(ctx, subpackage, releaseStockAction)
			_, err := app.Globals.SubPkgRepository.Update(ctx, *subpackage)
			if err != nil {
				logger.Err("SubPkgRepository.Update in %s state failed, orderId: %d, pid: %d, sid: %d, error: %s",
					state.Name(), subpackage.OrderId, subpackage.PId, subpackage.SId, err.Error())
				return
			} else {
				logger.Audit("%s state success, orderId: %d, pid: %d, sid: %d", state.Name(), subpackage.OrderId, subpackage.PId, subpackage.SId)
			}
		}
		order, err := app.Globals.OrderRepository.FindById(ctx, subpackages[0].OrderId)
		if err != nil {
			logger.Err("OrderRepository.FindById in %s state failed, orderId: %d, pid: %d, sid: %d, error: %s",
				state.Name(), subpackages[0].OrderId, subpackages[0].PId, subpackages[0].SId, err.Error())
			return
		}

		var findFlag = true
		for i := 0; i < len(order.Packages); i++ {
			findFlag = true
			for j := 0; j < len(order.Packages[i].Subpackages); j++ {
				if order.Packages[i].Subpackages[j].Status != states.PayToBuyer.StateName() &&
					order.Packages[i].Subpackages[j].Status != states.PayToSeller.StateName() {
					findFlag = false
					break
				}
			}

			if findFlag {
				state.SetPkgStatus(ctx, &order.Packages[i], states.PackageClosedStatus)
				_, err := app.Globals.PkgItemRepository.Update(ctx, order.Packages[i])
				if err != nil {
					logger.Err("update pkgItem status to closed failed, state: %s, orderId: %d, pid: %d, error: %s",
						state.Name(), order.Packages[i].OrderId, order.Packages[i].PId, err.Error())
				} else {
					logger.Audit("update pkgItem status to closed success, state: %s, orderId: %d, pid: %d",
						state.Name(), order.Packages[i].OrderId, order.Packages[i].PId)
				}
			}
		}

		findFlag = true
		for i := 0; i < len(order.Packages); i++ {
			if order.Packages[i].Status != string(states.PackageClosedStatus) {
				findFlag = false
				break
			}
		}

		if findFlag {
			state.SetOrderStatus(ctx, order, states.OrderClosedStatus)
			_, err := app.Globals.OrderRepository.Save(ctx, *order)
			if err != nil {
				logger.Err("update order status to closed failed, state: %s, orderId: %d, error: %s",
					state.Name(), order.OrderId, err.Error())
			} else {
				logger.Audit("update order status to closed success, state: %s, orderId: %d", state.Name(), order.OrderId)
			}
		}
	} else {
		logger.Err("iFrame.Header() not a subpackage or pid not found, state: %s iframe: %v", state.Name(), iFrame)
	}
}

func (state payToBuyerState) releasedStock(ctx context.Context, subpackages []*entities.Subpackage) error {

	var sids = make([]uint64, 0, len(subpackages))
	var inventories = make(map[string]int, 32)
	for _, subpackage := range subpackages {
		for z := 0; z < len(subpackage.Items); z++ {
			item := subpackage.Items[z]
			inventories[item.InventoryId] = int(item.Quantity)
		}
		sids = append(sids, subpackage.SId)
	}

	iFuture := app.Globals.StockService.BatchStockActions(ctx, inventories, system_action.New(system_action.StockRelease))
	futureData := iFuture.Get()
	if futureData.Error() != nil {
		logger.Err("Reserved stock from stockService failed, state: %s, orderId: %d, pid: %d, sids: %v, error: %s",
			state.Name(), subpackages[0].OrderId, subpackages[0].PId, sids, futureData.Error())
		return futureData.Error().Reason()
	}

	logger.Audit("Release stock success, state: %s, orderId: %d, pid: %d, sids: %v",
		state.Name(), subpackages[0].OrderId, subpackages[0].PId, sids)
	return nil
}
