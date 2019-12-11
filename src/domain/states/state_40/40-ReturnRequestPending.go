package state_40

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/app"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	"gitlab.faza.io/order-project/order-service/domain/events"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
	"gitlab.faza.io/order-project/order-service/infrastructure/future"
	"time"
)

const (
	stepName  string = "Return_Request_Pending"
	stepIndex int    = 40
)

type returnRequestPendingState struct {
	*states.BaseStateImpl
}

func New(childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &returnRequestPendingState{states.NewBaseStep(stepName, stepIndex, childes, parents, actionStateMap)}
}

func NewOf(name string, index int, childes, parents []states.IState, actionStateMap map[actions.IAction]states.IState) states.IState {
	return &returnRequestPendingState{states.NewBaseStep(name, index, childes, parents, actionStateMap)}
}

func NewFrom(base *states.BaseStateImpl) states.IState {
	return &returnRequestPendingState{base}
}

func NewValueOf(base *states.BaseStateImpl, params ...interface{}) states.IState {
	panic("implementation required")
}

func (state returnRequestPendingState) Process(ctx context.Context, iFrame frame.IFrame) {
	if iFrame.Header().KeyExists(string(frame.HeaderSubpackage)) {
		subpkg, ok := iFrame.Header().Value(string(frame.HeaderSubpackage)).(*entities.Subpackage)
		if !ok {
			logger.Err("iFrame.Header() not a subpackage, frame: %v, %s state ", iFrame, state.Name())
			return
		}

		// TODO must be read from reids config
		expiredTime := time.Now().UTC().Add(time.Hour*
			time.Duration(72) +
			time.Minute*time.Duration(0) +
			time.Second*time.Duration(0))

		state.UpdateSubPackage(ctx, subpkg, nil)
		if subpkg.Tracking.State != nil {
			subpkg.Tracking.State.Data = map[string]interface{}{
				"expiredTime": expiredTime,
			}
			logger.Audit("Process() => set expiredTime: %s , orderId: %d, sellerId: %d, sid: %d, %s state ",
				expiredTime, subpkg.OrderId, subpkg.SellerId, subpkg.SId, state.Name())
		}

		_, err := app.Globals.SubPkgRepository.Update(ctx, *subpkg)
		if err != nil {
			logger.Err("Process() => SubPkgRepository.Update in %s state failed, orderId: %d, sellerId: %d, sid: %d, error: %s",
				state.Name(), subpkg.OrderId, subpkg.SellerId, subpkg.SId, err.Error())
		} else {
			logger.Audit("Process() => Status of subpackage update to %s state, orderId: %d, sellerId: %d, sid: %d",
				state.Name(), subpkg.OrderId, subpkg.SellerId, subpkg.SId)
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

			// TODO cleaning subpackage after merging subpackages
			var newSubPackage *entities.Subpackage
			var nextActionState states.IState
			var returnRequestPendingAction *entities.Action

			// iterate subpackages
			for _, eventSubPkg := range actionData.SubPackages {
				for i := 0; i < len(pkgItem.Subpackages); i++ {
					if eventSubPkg.SId == pkgItem.Subpackages[i].SId && pkgItem.Subpackages[i].Status == state.Name() {
						var findAction = false
						for action, nextState := range state.StatesMap() {
							if action.ActionType().ActionName() == event.Action().ActionType().ActionName() &&
								action.ActionEnum().ActionName() == event.Action().ActionEnum().ActionName() {
								findAction = true

								//var newSubPkg *entities.Subpackage
								var newPkgItems []entities.Item

								// iterate items
								for _, actionItem := range eventSubPkg.Items {
									for j := 0; j < len(pkgItem.Subpackages[i].Items); j++ {
										if actionItem.InventoryId == pkgItem.Subpackages[i].Items[j].InventoryId {
											nextActionState = nextState

											if actionItem.Quantity != pkgItem.Subpackages[i].Items[j].Quantity {
												if newSubPackage == nil {
													newSubPackage = pkgItem.Subpackages[i].DeepCopy()
													newSubPackage.SId = 0
													newSubPackage.Items = make([]entities.Item, 0, len(eventSubPkg.Items))

													returnRequestPendingAction = &entities.Action{
														Name:      action.ActionEnum().ActionName(),
														Type:      action.ActionType().ActionName(),
														Result:    string(states.ActionSuccess),
														Reasons:   actionItem.Reasons,
														CreatedAt: time.Now().UTC(),
													}
												}

												if newPkgItems == nil {
													newPkgItems = make([]entities.Item, 0, len(pkgItem.Subpackages[i].Items))
												}

												pkgItem.Subpackages[i].Items[j].Quantity -= actionItem.Quantity
												pkgItem.Subpackages[i].Items[j].Invoice.Total = pkgItem.Subpackages[i].Items[j].Invoice.Unit *
													uint64(pkgItem.Subpackages[i].Items[j].Quantity)
												newPkgItem := pkgItem.Subpackages[i].Items[j].DeepCopy()
												newPkgItems = append(newPkgItems, *newPkgItem)

												newItem := pkgItem.Subpackages[i].Items[j].DeepCopy()
												newItem.Quantity = actionItem.Quantity
												newItem.Reasons = actionItem.Reasons
												newItem.Invoice.Total = newItem.Invoice.Unit * uint64(newItem.Quantity)
												if newSubPackage != nil {
													newSubPackage.Items = append(newSubPackage.Items, *newItem)
												}
											} else {
												// action contain item with all quantity
												newItem := pkgItem.Subpackages[i].Items[j].DeepCopy()
												newItem.Reasons = actionItem.Reasons
												if newSubPackage == nil {
													newSubPackage = pkgItem.Subpackages[i].DeepCopy()
													newSubPackage.SId = 0
													newSubPackage.Items = make([]entities.Item, 0, len(eventSubPkg.Items))

													returnRequestPendingAction = &entities.Action{
														Name:      action.ActionEnum().ActionName(),
														Type:      action.ActionType().ActionName(),
														Result:    string(states.ActionSuccess),
														Reasons:   actionItem.Reasons,
														CreatedAt: time.Now().UTC(),
													}
												}
												newSubPackage.Items = append(newSubPackage.Items, *newItem)
											}
										}
									}
								}

								// create diff packages
								if newPkgItems != nil {
									pkgItem.Subpackages[i].Items = newPkgItems
								} else {
									if newSubPackage != nil &&
										len(newSubPackage.Items) == len(pkgItem.Subpackages[i].Items) {
										pkgItem.Subpackages[i].Items = nil
									}
								}
							}
						}

						if !findAction {
							logger.Err("Process() => received action not acceptable, state: %s, event: %v", state.String(), event)
							future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
								SetError(future.NotAccepted, "Action Not Accepted", errors.New("Action Not Accepted")).Send()
							return
						}
					}
				}
			}

			if newSubPackage != nil {
				// remove subpackage with zero of items
				var subpackages = make([]entities.Subpackage, 0, len(pkgItem.Subpackages))
				for i := 0; i < len(pkgItem.Subpackages); i++ {
					if len(pkgItem.Subpackages[i].Items) > 0 {
						subpackages = append(subpackages, pkgItem.Subpackages[i])
					}
				}

				if len(pkgItem.Subpackages) != len(subpackages) {
					pkgItem.Subpackages = subpackages
					pkgItemUpdated, err := app.Globals.PkgItemRepository.Update(ctx, *pkgItem)
					if err != nil {
						logger.Err("Process() => PkgItemRepository.Update in %s state failed, orderId: %d, sellerId: %d, event: %v, error: %s", state.Name(),
							pkgItem.OrderId, pkgItem.PId, event, err.Error())
						// TODO must distinct system error from update version error
						future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
							SetError(future.InternalError, "Unknown Err", err).Send()
						return
					}
					pkgItem = pkgItemUpdated
				}

				state.UpdateSubPackage(ctx, newSubPackage, returnRequestPendingAction)
				err := app.Globals.SubPkgRepository.Save(ctx, newSubPackage)
				if err != nil {
					logger.Err("Process() => SubPkgRepository.Save in %s state failed, orderId: %d, sellerId: %d, event: %v, error: %s", state.Name(),
						newSubPackage.OrderId, newSubPackage.SellerId, event, err.Error())
					// TODO must distinct system error from update version error
					future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
						SetError(future.InternalError, "Unknown Err", err).Send()
					return
				} else {
					logger.Audit("Process() => Status of new subpackage update to %v event, orderId: %d, sellerId: %d, sid: %d",
						event, newSubPackage.OrderId, newSubPackage.SellerId, newSubPackage.SId)
				}

				if nextActionState != nil {
					pkgItemUpdated, err := app.Globals.PkgItemRepository.Update(ctx, *pkgItem)
					if err != nil {
						logger.Err("Process() => PkgItemRepository.Update in %s state failed, orderId: %d, sellerId: %d, event: %v, error: %s", state.Name(),
							pkgItem.OrderId, pkgItem.PId, event, err.Error())
					} else {
						pkgItem = pkgItemUpdated
					}

					response := events.ActionResponse{
						OrderId: newSubPackage.OrderId,
						SIds:    nil,
					}

					future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
						SetData(response).Send()
					nextActionState.Process(ctx, frame.Factory().SetSubpackage(newSubPackage).SetBody(pkgItem).Build())
				}
			} else {
				logger.Err("Process() => result of event invalid, state: %s, event: %v, orderId: %d, sellerId: %d",
					state.String(), event, pkgItem.OrderId, pkgItem.PId)
				future.FactoryOf(iFrame.Header().Value(string(frame.HeaderFuture)).(future.IFuture)).
					SetError(future.InternalError, "Unknown Err", errors.New("event type invalid")).Send()
				return
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
