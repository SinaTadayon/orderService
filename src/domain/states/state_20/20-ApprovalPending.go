package state_20

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/app"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	scheduler_action "gitlab.faza.io/order-project/order-service/domain/actions/scheduler"
	seller_action "gitlab.faza.io/order-project/order-service/domain/actions/seller"
	system_action "gitlab.faza.io/order-project/order-service/domain/actions/system"
	"gitlab.faza.io/order-project/order-service/domain/events"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
	"gitlab.faza.io/order-project/order-service/infrastructure/future"
	notify_service "gitlab.faza.io/order-project/order-service/infrastructure/services/notification"
	"gitlab.faza.io/order-project/order-service/infrastructure/utils"
	"strconv"
	"text/template"
	"time"
)

const (
	stepName  string = "Approval_Pending"
	stepIndex int    = 20
)

type approvalPendingState struct {
	*states.BaseStateImpl
}

func New(childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &approvalPendingState{states.NewBaseStep(stepName, stepIndex, childes, parents, actionStateMap)}
}

func NewOf(name string, index int, childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &approvalPendingState{states.NewBaseStep(name, index, childes, parents, actionStateMap)}
}

func NewFrom(base *states.BaseStateImpl) states.IState {
	return &approvalPendingState{base}
}

func NewValueOf(base *states.BaseStateImpl, params ...interface{}) states.IState {
	panic("implementation required")
}

func (state approvalPendingState) Process(ctx context.Context, iFrame frame.IFrame) {

	if iFrame.Header().KeyExists(string(frame.HeaderOrderId)) {
		if iFrame.Body().Content() == nil {
			logger.Err("Process() => iFrame.Body().Content() is nil, orderId: %d, %s state ",
				iFrame.Header().Value(string(frame.HeaderOrderId)), state.Name())
			return
		}

		order, ok := iFrame.Body().Content().(*entities.Order)
		if !ok {
			logger.Err("Process() => iFrame.Body().Content() not a order, orderId: %d, %s state ",
				iFrame.Header().Value(string(frame.HeaderOrderId)), state.Name())
			return
		}

		var sellerNotificationAction = &entities.Action{
			Name:      system_action.SellerNotification.ActionName(),
			Type:      "",
			UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
			UTP:       actions.System.ActionName(),
			Perm:      "",
			Priv:      "",
			Policy:    "",
			Result:    string(states.ActionSuccess),
			Reasons:   nil,
			Note:      "",
			Data:      nil,
			CreatedAt: time.Now().UTC(),
			Extended:  nil,
		}

		for i := 0; i < len(order.Packages); i++ {
			futureData := app.Globals.UserService.GetSellerProfile(ctx, strconv.Itoa(int(order.Packages[i].PId))).Get()
			if futureData.Error() != nil {
				logger.Err("Process() => UserService.GetSellerProfile failed, send sms message failed, state: %s, orderId: %d, pid: %d, error: %s",
					state.Name(), order.Packages[i].OrderId, order.Packages[i].PId, futureData.Error().Reason())
			} else {
				if futureData.Data() != nil {
					sellerProfile := futureData.Data().(*entities.SellerProfile)
					smsTemplate, err := template.New("SMS").Parse(app.Globals.SMSTemplate.OrderNotifySellerApprovalPendingState)
					if err != nil {
						logger.Err("Process() => smsTemplate.Parse failed, state: %s, orderId: %d, message: %s, err: %s",
							state.Name(), order.OrderId, app.Globals.SMSTemplate.OrderNotifySellerApprovalPendingState, err)
					} else {
						var buf bytes.Buffer
						err = smsTemplate.Execute(&buf, order.OrderId)
						newBuf := bytes.NewBuffer(bytes.Replace(buf.Bytes(), []byte("\\n"), []byte{10}, -1))
						if err != nil {
							logger.Err("Process() => smsTemplate.Execute failed, state: %s, orderId: %d, message: %s, err: %s",
								state.Name(), order.OrderId, app.Globals.SMSTemplate.OrderNotifySellerApprovalPendingState, err)
						} else {
							sellerNotify := notify_service.SMSRequest{
								Phone: sellerProfile.GeneralInfo.MobilePhone,
								Body:  newBuf.String(),
							}
							sellerFutureData := app.Globals.NotifyService.NotifyBySMS(ctx, sellerNotify).Get()
							if sellerFutureData.Error() != nil {
								logger.Err("Process() => NotifyService.NotifyBySMS failed, request: %v, state: %s, orderId: %d, pid: %d, error: %s",
									sellerNotify, state.Name(), order.Packages[i].OrderId, order.Packages[i].PId, sellerFutureData.Error().Reason())
								sellerNotificationAction = &entities.Action{
									Name:      system_action.SellerNotification.ActionName(),
									Type:      "",
									UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
									UTP:       actions.System.ActionName(),
									Perm:      "",
									Priv:      "",
									Policy:    "",
									Result:    string(states.ActionFail),
									Reasons:   nil,
									Note:      "",
									Data:      nil,
									CreatedAt: time.Now().UTC(),
									Extended:  nil,
								}
							} else {
								logger.Audit("Process() => NotifyService.NotifyBySMS success, sellerNotify: %v, state: %s, orderId: %d, pid: %d",
									sellerNotify, state.Name(), order.Packages[i].OrderId, order.Packages[i].PId)
								sellerNotificationAction = &entities.Action{
									Name:      system_action.SellerNotification.ActionName(),
									Type:      "",
									UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
									UTP:       actions.System.ActionName(),
									Perm:      "",
									Priv:      "",
									Policy:    "",
									Result:    string(states.ActionSuccess),
									Reasons:   nil,
									Note:      "",
									Data:      nil,
									CreatedAt: time.Now().UTC(),
									Extended:  nil,
								}
							}
							//update subpackage and package
							state.UpdatePackageAllSubPkg(ctx, order.Packages[i], sellerNotificationAction)
						}
					}
				} else {
					logger.Err("Process() => UserService.GetSellerProfile futureData.Data() is nil, send sms message failed, state: %s, orderId: %d, pid: %d",
						state.Name(), order.Packages[i].OrderId, order.Packages[i].PId)
				}
			}
		}

		var expireTime time.Time
		timeUnit := app.Globals.FlowManagerConfig[app.FlowManagerSchedulerStateTimeUintConfig].(string)
		if timeUnit == app.DurationTimeUnit {
			value := app.Globals.FlowManagerConfig[app.FlowManagerSchedulerApprovalPendingStateConfig].(time.Duration)
			expireTime = time.Now().UTC().Add(value)
		} else {
			value := app.Globals.FlowManagerConfig[app.FlowManagerSchedulerApprovalPendingStateConfig].(int)
			if timeUnit == app.HourTimeUnit {
				expireTime = time.Now().UTC().Add(
					time.Hour*time.Duration(value) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0))
			} else {
				expireTime = time.Now().UTC().Add(
					time.Hour*time.Duration(0) +
						time.Minute*time.Duration(value) +
						time.Second*time.Duration(0))
			}
		}

		order.UpdatedAt = time.Now().UTC()
		for i := 0; i < len(order.Packages); i++ {
			order.Packages[i].UpdatedAt = time.Now().UTC()
			for j := 0; j < len(order.Packages[i].Subpackages); j++ {
				order.Packages[i].Subpackages[j].Tracking.State.Schedulers = []*entities.SchedulerData{
					{
						order.OrderId,
						order.Packages[i].PId,
						order.Packages[i].Subpackages[j].SId,
						state.Name(),
						state.Index(),
						states.SchedulerJobName,
						states.SchedulerGroupName,
						scheduler_action.Cancel.ActionName(),
						0,
						0,
						"",
						nil,
						nil,
						string(states.SchedulerSubpackageStateExpire),
						"",
						nil,
						true,
						expireTime,
						time.Now().UTC(),
						time.Now().UTC(),
						nil,
						nil,
					},
				}
				state.UpdateSubPackage(ctx, order.Packages[i].Subpackages[j], sellerNotificationAction)
			}
		}

		_, err := app.Globals.OrderRepository.Save(ctx, *order)
		if err != nil {
			logger.Err("Process() => OrderRepository.Save in %s state failed, orderId: %d, error: %v", state.Name(), order.OrderId, err)
		} else {
			logger.Audit("Process() => Order state of all subpackages update to %s state, orderId: %d", state.Name(), order.OrderId)
		}
	} else if iFrame.Header().KeyExists(string(frame.HeaderEvent)) {
		event, ok := iFrame.Header().Value(string(frame.HeaderEvent)).(events.IEvent)
		if !ok {
			logger.Err("Process() => received frame doesn't have a event, state: %s, frame: %v", state.String(), iFrame)
			future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
				SetError(future.InternalError, "Unknown Err", nil).Send()
			return
		}

		if event.EventType() == events.Action {
			pkgItem, ok := iFrame.Body().Content().(*entities.PackageItem)
			if !ok {
				logger.Err("Process() => received frame body not a PackageItem, state: %s, event: %v, frame: %v", state.String(), event, iFrame)
				future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
					SetError(future.InternalError, "Unknown Err", errors.New("frame body invalid")).Send()
				return
			}

			actionData, ok := event.Data().(events.ActionData)
			if !ok {
				logger.Err("Process() => received action event data invalid, state: %s, event: %v", state.String(), event)
				future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
					SetError(future.InternalError, "Unknown Err", errors.New("Action Data event invalid")).Send()
				return
			}

			var newSubPackages []*entities.Subpackage
			var requestAction *entities.Action
			var newSubPkg *entities.Subpackage
			var fullItems []*entities.Item
			var nextActionState states.IState
			var actionState actions.IAction

			for action, nextState := range state.StatesMap() {
				if action.ActionType().ActionName() == event.Action().ActionType().ActionName() &&
					action.ActionEnum().ActionName() == event.Action().ActionEnum().ActionName() {
					nextActionState = nextState
					actionState = action
					break
				}
			}

			if nextActionState == nil || actionState == nil {
				logger.Err("Process() => received action not acceptable, state: %s, event: %v", state.String(), event)
				future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
					SetError(future.NotAccepted, "Action Not Accepted", errors.New("Action Not Accepted")).Send()
				return
			}

			// iterate subpackages
			for _, eventSubPkg := range actionData.SubPackages {
				for i := 0; i < len(pkgItem.Subpackages); i++ {
					if eventSubPkg.SId == pkgItem.Subpackages[i].SId && pkgItem.Subpackages[i].Status == state.Name() {
						newSubPkg = nil
						fullItems = nil
						var findItem = false

						// iterate items
						for _, actionItem := range eventSubPkg.Items {
							findItem = false
							for j := 0; j < len(pkgItem.Subpackages[i].Items); j++ {
								if actionItem.InventoryId == pkgItem.Subpackages[i].Items[j].InventoryId {
									findItem = true

									// create new subpackages which contains new items along
									// with new quantity and recalculated related invoice
									if actionItem.Quantity < pkgItem.Subpackages[i].Items[j].Quantity {
										if newSubPkg == nil {
											newSubPkg = pkgItem.Subpackages[i].DeepCopy()
											newSubPkg.SId = 0
											newSubPkg.Items = make([]*entities.Item, 0, len(eventSubPkg.Items))
											newSubPkg.CreatedAt = time.Now().UTC()
											newSubPkg.UpdatedAt = time.Now().UTC()

											requestAction = &entities.Action{
												Name:      actionState.ActionEnum().ActionName(),
												Type:      "",
												UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
												UTP:       actionState.ActionType().ActionName(),
												Perm:      "",
												Priv:      "",
												Policy:    "",
												Result:    string(states.ActionSuccess),
												Reasons:   actionItem.Reasons,
												Note:      "",
												Data:      nil,
												CreatedAt: time.Now().UTC(),
												Extended:  nil,
											}
										}

										unit, err := decimal.NewFromString(pkgItem.Subpackages[i].Items[j].Invoice.Unit.Amount)
										if err != nil {
											logger.Err("Process() => decimal.NewFromString failed, Unit.Amount invalid, unit: %s, orderId: %d, pid: %d, sid: %d, state: %s, event: %v",
												pkgItem.Subpackages[i].Items[j].Invoice.Unit.Amount, pkgItem.Subpackages[i].OrderId, pkgItem.Subpackages[i].PId, pkgItem.Subpackages[i].SId, state.Name(), event)
											future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
												SetError(future.InternalError, "Unknown Err", errors.New("Subpackage Unit invalid")).Send()
											return
										}

										pkgItem.Subpackages[i].Items[j].Quantity -= actionItem.Quantity
										pkgItem.Subpackages[i].Items[j].Invoice.Total.Amount = strconv.Itoa(int(unit.IntPart() * int64(pkgItem.Subpackages[i].Items[j].Quantity)))

										// create new item from requested action item
										newItem := pkgItem.Subpackages[i].Items[j].DeepCopy()
										newItem.Quantity = actionItem.Quantity
										newItem.Reasons = actionItem.Reasons
										newItem.Invoice.Total.Amount = strconv.Itoa(int(unit.IntPart() * int64(newItem.Quantity)))
										newSubPkg.Items = append(newSubPkg.Items, newItem)

									} else if actionItem.Quantity > pkgItem.Subpackages[i].Items[j].Quantity {
										logger.Err("Process() => received action not acceptable, Requested quantity greater than item quantity, state: %s, event: %v", state.String(), event)
										future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
											SetError(future.NotAccepted, "Requested quantity greater than item quantity", errors.New("Action Not Accepted")).Send()
										return

									} else {
										if fullItems == nil {
											fullItems = make([]*entities.Item, 0, len(pkgItem.Subpackages[i].Items))
											requestAction = &entities.Action{
												Name:      actionState.ActionEnum().ActionName(),
												Type:      "",
												UId:       ctx.Value(string(utils.CtxUserID)).(uint64),
												UTP:       actionState.ActionType().ActionName(),
												Perm:      "",
												Priv:      "",
												Policy:    "",
												Result:    string(states.ActionSuccess),
												Reasons:   actionItem.Reasons,
												Data:      nil,
												CreatedAt: time.Now().UTC(),
												Extended:  nil,
											}
										}
										fullItems = append(fullItems, pkgItem.Subpackages[i].Items[j])
										pkgItem.Subpackages[i].Items[len(pkgItem.Subpackages[i].Items)-1], pkgItem.Subpackages[i].Items[j] =
											pkgItem.Subpackages[i].Items[j], pkgItem.Subpackages[i].Items[len(pkgItem.Subpackages[i].Items)-1]
										pkgItem.Subpackages[i].Items = pkgItem.Subpackages[i].Items[:len(pkgItem.Subpackages[i].Items)-1]
									}
								}
							}
							if !findItem {
								logger.Err("Process() => received action item inventory not found, Requested action item inventory not found in requested subpackage, inventoryId: %s, state: %s, event: %v", actionItem.InventoryId, state.String(), event)
								future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
									SetError(future.NotFound, "Request action item not found", errors.New("Action Item Not Found")).Send()
								return
							}
						}

						newSubPackages = make([]*entities.Subpackage, 0, len(actionData.SubPackages))
						if newSubPkg != nil {
							if fullItems != nil {
								for z := 0; z < len(fullItems); z++ {
									newSubPkg.Items = append(newSubPkg.Items, fullItems[z])
								}
							}
							newSubPackages = append(newSubPackages, newSubPkg)
						} else {
							for z := 0; z < len(fullItems); z++ {
								pkgItem.Subpackages[i].Items = append(pkgItem.Subpackages[i].Items, fullItems[z])
							}
							newSubPackages = append(newSubPackages, pkgItem.Subpackages[i])
						}
					}
				}
			}

			if newSubPackages != nil {
				var sids = make([]uint64, 0, 32)
				for i := 0; i < len(newSubPackages); i++ {
					if newSubPackages[i].SId == 0 {
						newSid, err := app.Globals.SubPkgRepository.GenerateUniqSid(ctx, pkgItem.OrderId)
						if err != nil {
							logger.Err("Process() => SubPkgRepository.GenerateUniqSid failed, orderId: %d, pid: %d, sid: %d, state: %s, event: %v",
								newSubPackages[i].OrderId, newSubPackages[i].PId, newSubPackages[i].SId, state.Name(), event)
							future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
								SetError(future.ErrorCode(err.Code()), err.Message(), err.Reason()).Send()
							return
						}
						newSubPackages[i].SId = newSid
						pkgItem.Subpackages = append(pkgItem.Subpackages, newSubPackages[i])
					}
					sids = append(sids, newSubPackages[i].SId)
					state.UpdateSubPackage(ctx, newSubPackages[i], requestAction)
				}

				if event.Action().ActionEnum() != seller_action.Approve {
					var rejectedSubtotal int64 = 0
					var rejectedDiscount int64 = 0

					for i := 0; i < len(newSubPackages); i++ {
						for j := 0; j < len(newSubPackages[i].Items); j++ {
							amount, err := decimal.NewFromString(newSubPackages[i].Items[j].Invoice.Total.Amount)
							if err != nil {
								logger.Err("Process() => decimal.NewFromString failed, Total.Amount invalid, total: %s, orderId: %d, pid: %d, sid: %d, state: %s, event: %v",
									newSubPackages[i].Items[j].Invoice.Total.Amount, newSubPackages[i].OrderId, newSubPackages[i].PId, newSubPackages[i].SId, state.Name(), event)
								future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
									SetError(future.InternalError, "Unknown Error", errors.New("Subpackage Total Invalid")).Send()
								return
							}

							discount, err := decimal.NewFromString(newSubPackages[i].Items[j].Invoice.Discount.Amount)
							if err != nil {
								logger.Err("Process() => decimal.NewFromString failed, Invoice.Discount invalid, discount: %s, orderId: %d, pid: %d, sid: %d, state: %s, event: %v",
									newSubPackages[i].Items[j].Invoice.Discount, newSubPackages[i].OrderId, newSubPackages[i].PId, newSubPackages[i].SId, state.Name(), event)
								future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
									SetError(future.InternalError, "Unknown Error", errors.New("Subpackage Discount Invalid")).Send()
								return
							}

							rejectedSubtotal += amount.IntPart()
							rejectedDiscount += discount.IntPart()
						}
					}

					subtotal, err := decimal.NewFromString(pkgItem.Invoice.Subtotal.Amount)
					if err != nil {
						logger.Err("Process() => decimal.NewFromString failed, Subtotal.Amount invalid, subtotal: %s, orderId: %d, pid: %d, state: %s, event: %v",
							pkgItem.Invoice.Subtotal.Amount, pkgItem.OrderId, pkgItem.PId, state.Name(), event)
						future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
							SetError(future.InternalError, "Unknown Error", errors.New("Package Invoice Invalid")).Send()
						return
					}

					pkgDiscount, err := decimal.NewFromString(pkgItem.Invoice.Discount.Amount)
					if err != nil {
						logger.Err("Process() => decimal.NewFromString failed, Pkg Discount.Amount invalid, pkg discount: %s, orderId: %d, pid: %d, state: %s, event: %v",
							pkgItem.Invoice.Discount.Amount, pkgItem.OrderId, pkgItem.PId, state.Name(), event)
						future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
							SetError(future.InternalError, "Unknown Error", errors.New("Package Invoice Invalid")).Send()
						return
					}

					if rejectedSubtotal < subtotal.IntPart() && rejectedDiscount < pkgDiscount.IntPart() {
						pkgItem.Invoice.Subtotal.Amount = strconv.Itoa(int(subtotal.IntPart() - rejectedSubtotal))
						pkgItem.Invoice.Discount.Amount = strconv.Itoa(int(pkgDiscount.IntPart() - rejectedDiscount))
						logger.Audit("Process() => calculate package invoice success, orderId: %d, pid:%d, action: %s, subtotal: %s, discount: %s",
							pkgItem.OrderId, pkgItem.PId, event.Action().ActionEnum().ActionName(), pkgItem.Invoice.Subtotal.Amount, pkgItem.Invoice.Discount.Amount)

					} else if rejectedSubtotal > subtotal.IntPart() || rejectedDiscount > pkgDiscount.IntPart() {
						logger.Err("Process() => calculate package invoice failed, orderId: %d, pid:%d, action: %s, subtotal: %s, discount: %s",
							pkgItem.OrderId, pkgItem.PId, event.Action().ActionEnum().ActionName(), pkgItem.Invoice.Subtotal.Amount, pkgItem.Invoice.Discount.Amount)
						future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
							SetError(future.InternalError, "Unknown Error", errors.New("Package Invoice Invalid")).Send()
						return
					}
				}

				//pkgItemUpdated, newSids, err := app.Globals.PkgItemRepository.UpdateWithUpsert(ctx, *pkgItem)
				//if err != nil {
				//	logger.Err("Process() => PkgItemRepository.Update failed, state: %s, orderId: %d, pid: %d, sids: %v, event: %v, error: %v", state.Name(),
				//		pkgItem.OrderId, pkgItem.PId, sids, event, err)
				//	future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
				//		SetError(future.ErrorCode(err.Code()), err.Message(), err.Reason()).Send()
				//	return
				//}
				//sids = append(sids, newSids...)
				//pkgItem = pkgItemUpdated

				//response := events.ActionResponse{
				//	OrderId: pkgItem.OrderId,
				//	SIds:    sids,
				//}

				logger.Audit("Process() => Status of subpackages update success, state: %s, action: %s, orderId: %d, pid: %d, sids: %d",
					state.Name(), event.Action().ActionEnum().ActionName(), pkgItem.OrderId, pkgItem.PId, sids)

				//future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).SetData(response).Send()
				nextActionState.Process(ctx, frame.FactoryOf(iFrame).SetSIds(sids).SetBody(pkgItem).Build())
			} else {
				logger.Err("Process() => event action data invalid, state: %s, event: %v, frame: %v", state.String(), event, iFrame)
				future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
					SetError(future.BadRequest, "Event Action Data Invalid", errors.New("event action data invalid")).Send()
			}
		} else {
			logger.Err("Process() => event type not supported, state: %s, event: %v, frame: %v", state.String(), event, iFrame)
			future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
				SetError(future.InternalError, "Unknown Err", errors.New("event type invalid")).Send()
			return
		}
	} else {
		logger.Err("HeaderOrderId or HeaderEvent of iFrame.Header not found, state: %s iframe: %v", state.Name(), iFrame)
	}
}
