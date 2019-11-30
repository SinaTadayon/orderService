package operator_action_state

import (
	"context"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	"gitlab.faza.io/order-project/order-service/domain/events"
	"gitlab.faza.io/order-project/order-service/domain/states_old"
	listener_state "gitlab.faza.io/order-project/order-service/domain/states_old/listener"
	"gitlab.faza.io/order-project/order-service/infrastructure/promise"
)

const (
	actorType        = actions.Operator
	stateName string = "Operator_Action_State"
)

type operatorActionListener struct {
	*listener_state.BaseListenerImpl
}

func New(index int, childes, parents []states_old.IState, actions actions.IAction) listener_state.IListenerState {
	return &operatorActionListener{listener_state.NewBaseListener(stateName, index, childes, parents,
		actions, actorType)}
}

func NewOf(name string, index int, childes, parents []states_old.IState, actions actions.IAction) listener_state.IListenerState {
	return &operatorActionListener{listener_state.NewBaseListener(name, index, childes, parents,
		actions, actorType)}
}

func NewFrom(base *listener_state.BaseListenerImpl) listener_state.IListenerState {
	return &operatorActionListener{base}
}

func NewValueOf(base *listener_state.BaseListenerImpl, params ...interface{}) listener_state.IListenerState {
	panic("implementation required")
}

func (operatorAction operatorActionListener) ActionListener(ctx context.Context, event events.IEvent, param interface{}) promise.IPromise {
	panic("implementation required")
}
