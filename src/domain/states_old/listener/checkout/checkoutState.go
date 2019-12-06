package checkout_action_state

import (
	"context"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	checkout_action "gitlab.faza.io/order-project/order-service/domain/actions/checkout"
	"gitlab.faza.io/order-project/order-service/domain/events"
	actor_event "gitlab.faza.io/order-project/order-service/domain/events/actor"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states_old"
	launcher_state "gitlab.faza.io/order-project/order-service/domain/states_old/launcher"
	listener_state "gitlab.faza.io/order-project/order-service/domain/states_old/listener"
	"gitlab.faza.io/order-project/order-service/infrastructure/future"
	"gitlab.faza.io/order-project/order-service/infrastructure/global"
	"time"
	//message "gitlab.faza.io/protos/order"
)

const (
	actorType        = actions.Checkout
	stateName string = "Checkout_Action_State"
)

type checkoutActionListener struct {
	*listener_state.BaseListenerImpl
}

func New(index int, childes, parents []states_old.IState, actions actions.IAction) listener_state.IListenerState {
	return &checkoutActionListener{listener_state.NewBaseListener(stateName, index, childes, parents,
		actions, actorType)}
}

func NewOf(name string, index int, childes, parents []states_old.IState, actions actions.IAction) listener_state.IListenerState {
	return &checkoutActionListener{listener_state.NewBaseListener(name, index, childes, parents,
		actions, actorType)}
}

func NewFrom(base *listener_state.BaseListenerImpl) listener_state.IListenerState {
	return &checkoutActionListener{base}
}

func NewValueOf(base *listener_state.BaseListenerImpl, params ...interface{}) listener_state.IListenerState {
	panic("implementation required")
}

// TODO context handling
func (checkoutActionState checkoutActionListener) ActionListener(ctx context.Context, event events.IEvent, param interface{}) future.IFuture {

	if event == nil {
		logger.Err("Received Event is nil")
		returnChannel := make(chan future.IDataFuture, 1)
		returnChannel <- future.IDataFuture{Data: nil, Ex: future.FutureError{Code: future.InternalError, Reason: "Unknown Error"}}
		defer close(returnChannel)
		return future.NewFuture(returnChannel, 1, 1)
	}

	stockState, ok := checkoutActionState.Childes()[0].(launcher_state.ILauncherState)
	if ok != true {
		logger.Err("StockState isn't child of CheckoutState, event: %v", event)
		returnChannel := make(chan future.IDataFuture, 1)
		returnChannel <- future.IDataFuture{Data: nil, Ex: future.FutureError{Code: future.InternalError, Reason: "Unknown Error"}}
		defer close(returnChannel)
		return future.NewFuture(returnChannel, 1, 1)
	}

	// TODO checking type and result cast
	actorEvent := event.(actor_event.IActorEvent)

	if actorEvent.ActorType() != actions.Checkout {
		logger.Err("Received actorType of event is not Checkout, event: %v", event)
		returnChannel := make(chan future.IDataFuture, 1)
		returnChannel <- future.IDataFuture{Data: nil, Ex: future.FutureError{Code: future.InternalError, Reason: "Unknown Error"}}
		defer close(returnChannel)
		return future.NewFuture(returnChannel, 1, 1)
	}

	if actorEvent.ActorAction().ActionEnum()[0] != checkout_action.NewOrderAction {
		logger.Err("Received actorAction of event is not NewOrderAction, event: %v", event)
		returnChannel := make(chan future.IDataFuture, 1)
		returnChannel <- future.IDataFuture{Data: nil, Ex: future.FutureError{Code: future.InternalError, Reason: "Unknown Error"}}
		defer close(returnChannel)
		return future.NewFuture(returnChannel, 1, 1)
	}

	newOrder := actorEvent.Order()
	checkoutActionState.updateOrderStates(ctx, &newOrder, event.Timestamp())

	order, err := global.Singletons.OrderRepository.Save(newOrder)
	if err != nil {
		logger.Err("Save NewOrder Failed, error: %s, newOrder: %v", err, newOrder)
		returnChannel := make(chan future.IDataFuture, 1)
		returnChannel <- future.IDataFuture{Data: nil, Ex: future.FutureError{Code: future.InternalError, Reason: "Unknown Error"}}
		defer close(returnChannel)
		return future.NewFuture(returnChannel, 1, 1)
	}

	return stockState.ActionLauncher(ctx, *order, nil, nil)
}

func (checkoutActionState checkoutActionListener) updateOrderStates(ctx context.Context, newOrder *entities.Order, timestamp time.Time) {
	//for i := 0; i < len(newOrder.Items); i++ {
	//	newOrder.Items[i].Tracking.CurrentState.ActionName = checkoutActionState.ActionName()
	//	newOrder.Items[i].Tracking.CurrentState.Index = checkoutActionState.Index()
	//	newOrder.Items[i].Tracking.CurrentState.Type = checkoutActionState.Actions().ActionType().ActionName()
	//	newOrder.Items[i].Tracking.CurrentState.CreatedAt = time.Now().UTC()
	//	newOrder.Items[i].Tracking.CurrentState.Result = true
	//	newOrder.Items[i].Tracking.CurrentState.Reason = ""
	//
	//	newOrder.Items[i].Tracking.CurrentState.AcceptedAction.ActionName = checkout_action.NewOrderAction.String()
	//	newOrder.Items[i].Tracking.CurrentState.AcceptedAction.Type = actors.Checkout.String()
	//	newOrder.Items[i].Tracking.CurrentState.AcceptedAction.Base = actions.ActorAction.String()
	//	newOrder.Items[i].Tracking.CurrentState.AcceptedAction.Get = nil
	//	newOrder.Items[i].Tracking.CurrentState.AcceptedAction.Time = &timestamp
	//
	//	newOrder.Items[i].Tracking.CurrentState.Actions = []entities.Actions{newOrder.Items[i].Tracking.CurrentState.AcceptedAction}
	//
	//	stateHistory := entities.StateHistory {
	//		ActionName: newOrder.Items[i].Tracking.CurrentState.ActionName,
	//		Index: newOrder.Items[i].Tracking.CurrentState.Index,
	//		Type: newOrder.Items[i].Tracking.CurrentState.Type,
	//		Actions: newOrder.Items[i].Tracking.CurrentState.AcceptedAction,
	//		Result: newOrder.Items[i].Tracking.CurrentState.Result,
	//		Reason: newOrder.Items[i].Tracking.CurrentState.Reason,
	//		CreatedAt:newOrder.Items[i].Tracking.CurrentState.CreatedAt,
	//	}
	//
	//	newOrder.Items[i].Tracking.History[len(newOrder.Items[i].Tracking.History)].History = append(newOrder.Items[i].Tracking.History[len(newOrder.Items[i].Tracking.History)].History, stateHistory)
	//}
}
