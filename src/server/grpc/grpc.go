package grpc_server

import (
	"context"
	"path"
	"runtime/debug"
	"strconv"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"gitlab.faza.io/order-project/order-service/domain/actions"
	buyer_action "gitlab.faza.io/order-project/order-service/domain/actions/buyer"
	operator_action "gitlab.faza.io/order-project/order-service/domain/actions/operator"
	scheduler_action "gitlab.faza.io/order-project/order-service/domain/actions/scheduler"
	seller_action "gitlab.faza.io/order-project/order-service/domain/actions/seller"
	"gitlab.faza.io/order-project/order-service/domain/events"
	"gitlab.faza.io/order-project/order-service/domain/models"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	"gitlab.faza.io/order-project/order-service/domain/states"
	"gitlab.faza.io/order-project/order-service/infrastructure/frame"
	applog "gitlab.faza.io/order-project/order-service/infrastructure/logger"
	"gitlab.faza.io/order-project/order-service/infrastructure/utils"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"gitlab.faza.io/order-project/order-service/app"
	"gitlab.faza.io/order-project/order-service/domain"
	"gitlab.faza.io/order-project/order-service/infrastructure/future"
	pb "gitlab.faza.io/protos/order"
	pg "gitlab.faza.io/protos/payment-gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"net"

	"gitlab.faza.io/go-framework/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type RequestADT string
type RequestType string
type RequestName string
type RequestMethod string
type UserType string
type SortDirection string
type FilterType string
type FilterValue string
type ActionType string
type Action string

const (
	PostMethod RequestMethod = "POST"
	GetMethod  RequestMethod = "GET"
)

const (
	OrderStateFilterType FilterType = "OrderState"
)

const (
	NewOrderFilter                 FilterValue = "NewOrder"
	PaymentPendingFilter           FilterValue = "PaymentPending"
	PaymentSuccessFilter           FilterValue = "PaymentSuccess"
	PaymentFailedFilter            FilterValue = "PaymentFailed"
	OrderVerificationPendingFilter FilterValue = "OrderVerificationPending"
	OrderVerificationSuccessFilter FilterValue = "OrderVerificationSuccess"
	OrderVerificationFailedFilter  FilterValue = "OrderVerificationFailed"
	ApprovalPendingFilter          FilterValue = "ApprovalPending"
	CanceledBySellerFilter         FilterValue = "CanceledBySeller"
	CanceledByBuyerFilter          FilterValue = "CanceledByBuyer"
	ShipmentPendingFilter          FilterValue = "ShipmentPending"
	ShipmentDelayedFilter          FilterValue = "ShipmentDelayed"
	ShippedFilter                  FilterValue = "Shipped"
	DeliveryPendingFilter          FilterValue = "DeliveryPending"
	DeliveryDelayedFilter          FilterValue = "DeliveryDelayed"
	DeliveredFilter                FilterValue = "Delivered"
	DeliveryFailedFilter           FilterValue = "DeliveryFailed"
	ReturnRequestPendingFilter     FilterValue = "ReturnRequestPending"
	ReturnRequestRejectedFilter    FilterValue = "ReturnRequestRejected"
	ReturnCanceledFilter           FilterValue = "ReturnCanceled"
	ReturnShipmentPendingFilter    FilterValue = "ReturnShipmentPending"
	ReturnShippedFilter            FilterValue = "ReturnShipped"
	ReturnDeliveryPendingFilter    FilterValue = "ReturnDeliveryPending"
	ReturnDeliveryDelayedFilter    FilterValue = "ReturnDeliveryDelayed"
	ReturnDeliveredFilter          FilterValue = "ReturnDelivered"
	ReturnDeliveryFailedFilter     FilterValue = "ReturnDeliveryFailed"
	ReturnRejectedFilter           FilterValue = "ReturnRejected"
	PayToBuyerFilter               FilterValue = "PayToBuyer"
	PayToSellerFilter              FilterValue = "PayToSeller"

	AllOrdersFilter             FilterValue = "AllOrders"
	AllCanceledFilter           FilterValue = "AllCanceled"
	DashboardReportFilter       FilterValue = "DashboardReport"
	ShipmentReportFilter        FilterValue = "ShipmentReport"
	ReturnReportFilter          FilterValue = "ReturnReport"
	DeliveredReportFilter       FilterValue = "DeliveredReport"
	CanceledReportFilter        FilterValue = "CanceledReport"
	AllReportFilter             FilterValue = "AllReport"
	ApprovalPendingReportFilter FilterValue = "ApprovalPendingReport"
)

const (
	DeliverAction             Action = "Deliver"
	DeliveryFailAction        Action = "DeliveryFail"
	DeliveryDelayAction       Action = "DeliveryDelay"
	DeliveryPendingAction     Action = "DeliveryPending"
	SubmitReturnRequestAction Action = "SubmitReturnRequest"
	EnterShipmentDetailAction Action = "EnterShipmentDetail"
	ApproveAction             Action = "Approve"
	RejectAction              Action = "Reject"
	CancelAction              Action = "Cancel"
	AcceptAction              Action = "Accept"
	CloseAction               Action = "Close"
)

const (
	DataReqType   RequestType = "Data"
	ActionReqType RequestType = "Action"
)

const (
	ListType   RequestADT = "List"
	SingleType RequestADT = "Single"
)

const (
	OperatorUser  UserType = "Operator"
	SellerUser    UserType = "Seller"
	BuyerUser     UserType = "Buyer"
	SchedulerUser UserType = "Schedulers"
	FinanceUser   UserType = "Finance"
)

const (
	//SellerAllOrders             		RequestName = "SellerAllOrders"
	SellerOrderList                   RequestName = "SellerOrderList"
	SellerOrderDetail                 RequestName = "SellerOrderDetail"
	SellerReturnOrderList             RequestName = "SellerReturnOrderList"
	SellerReturnOrderDetail           RequestName = "SellerReturnOrderDetail"
	SellerOrderDashboardReports       RequestName = "SellerOrderDashboardReports"
	SellerOrderShipmentReports        RequestName = "SellerOrderShipmentReports"
	SellerOrderDeliveredReports       RequestName = "SellerOrderDeliveredReports"
	SellerOrderReturnReports          RequestName = "SellerOrderReturnReports"
	SellerOrderCancelReports          RequestName = "SellerOrderCancelReports"
	SellerAllOrderReports             RequestName = "SellerAllOrderReports"
	SellerApprovalPendingOrderReports RequestName = "SellerApprovalPendingOrderReports"

	//BuyerAllOrders			   RequestName = "BuyerAllOrders"
	//BuyerAllReturnOrders       RequestName = "BuyerAllReturnOrders"
	BuyerOrderDetailList       RequestName = "BuyerOrderDetailList"
	BuyerAllOrderReports       RequestName = "BuyerAllOrderReports"
	BuyerReturnOrderReports    RequestName = "BuyerReturnOrderReports"
	BuyerReturnOrderDetailList RequestName = "BuyerReturnOrderDetailList"

	//OperatorAllOrders	RequestName = "OperatorAllOrders"
	OperatorOrderList          RequestName = "OperatorOrderList"
	OperatorOrderDetail        RequestName = "OperatorOrderDetail"
	OperatorOrderInvoiceDetail RequestName = "OperatorOrderInvoiceDetail"
)

const (
	ASC  SortDirection = "ASC"
	DESC SortDirection = "DESC"
)

const (
	// ISO8601 standard time format
	ISO8601 = "2006-01-02T15:04:05-0700"
)

type stackTraceDisabler struct{}

func (s stackTraceDisabler) Enabled(zapcore.Level) bool {
	return false
}

type FilterState struct {
	expectedState []states.IEnumState
	actualState   states.IEnumState
}

type FilterQueryState struct {
	state     states.IEnumState
	queryPath string
}

type Server struct {
	pb.UnimplementedOrderServiceServer
	pg.UnimplementedBankResultHookServer
	flowManager          domain.IFlowManager
	address              string
	port                 uint16
	requestFilters       map[RequestName][]FilterValue
	buyerFilterStates    map[FilterValue][]FilterState
	buyerAllStatesMap    map[string][]states.IEnumState
	buyerReturnStatesMap map[string][]states.IEnumState
	sellerFilterStates   map[FilterValue][]FilterState
	sellerStatesMap      map[string][]states.IEnumState
	operatorFilterStates map[FilterValue][]FilterState
	queryPathStates      map[FilterValue]FilterQueryState
	actionStates         map[UserType][]actions.IAction
	reasonsMap           reasonsMap
}

func NewServer(address string, port uint16, flowManager domain.IFlowManager) Server {
	buyerStatesMap := make(map[FilterValue][]FilterState, 8)
	buyerStatesMap[NewOrderFilter] = []FilterState{{[]states.IEnumState{states.NewOrder}, states.NewOrder}}
	buyerStatesMap[PaymentPendingFilter] = []FilterState{{[]states.IEnumState{states.PaymentPending}, states.PaymentPending}}
	buyerStatesMap[PaymentSuccessFilter] = []FilterState{{[]states.IEnumState{states.PaymentSuccess}, states.ApprovalPending}}
	buyerStatesMap[PaymentFailedFilter] = []FilterState{{[]states.IEnumState{states.PaymentFailed}, states.PaymentFailed}}
	buyerStatesMap[OrderVerificationPendingFilter] = []FilterState{{[]states.IEnumState{states.OrderVerificationPending}, states.ApprovalPending}}
	buyerStatesMap[OrderVerificationSuccessFilter] = []FilterState{{[]states.IEnumState{states.OrderVerificationSuccess}, states.ApprovalPending}}
	buyerStatesMap[OrderVerificationFailedFilter] = []FilterState{{[]states.IEnumState{states.OrderVerificationFailed}, states.PayToBuyer}}
	buyerStatesMap[ApprovalPendingFilter] = []FilterState{{[]states.IEnumState{states.ApprovalPending}, states.ApprovalPending}}
	buyerStatesMap[CanceledBySellerFilter] = []FilterState{{[]states.IEnumState{states.CanceledBySeller}, states.PayToBuyer}}
	buyerStatesMap[CanceledByBuyerFilter] = []FilterState{{[]states.IEnumState{states.CanceledByBuyer}, states.PayToBuyer}}
	buyerStatesMap[ShipmentPendingFilter] = []FilterState{{[]states.IEnumState{states.ShipmentPending}, states.ShipmentPending}}
	buyerStatesMap[ShipmentDelayedFilter] = []FilterState{{[]states.IEnumState{states.ShipmentDelayed}, states.ShipmentDelayed}}
	buyerStatesMap[ShippedFilter] = []FilterState{{[]states.IEnumState{states.Shipped}, states.Shipped}}
	buyerStatesMap[DeliveredFilter] = []FilterState{{[]states.IEnumState{states.DeliveryPending}, states.DeliveryPending}, {[]states.IEnumState{states.DeliveryDelayed}, states.DeliveryDelayed}, {[]states.IEnumState{states.Delivered}, states.Delivered}}
	buyerStatesMap[DeliveryFailedFilter] = []FilterState{{[]states.IEnumState{states.DeliveryFailed}, states.PayToBuyer}}
	buyerStatesMap[ReturnRequestPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnRequestPending}, states.ReturnRequestPending}, {[]states.IEnumState{states.ReturnRequestRejected}, states.ReturnRequestRejected}}
	buyerStatesMap[ReturnShipmentPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnShipmentPending}, states.ReturnShipmentPending}}
	buyerStatesMap[ReturnShippedFilter] = []FilterState{{[]states.IEnumState{states.ReturnShipped}, states.ReturnShipped}}
	buyerStatesMap[ReturnDeliveredFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryPending}, states.ReturnDeliveryPending}, {[]states.IEnumState{states.ReturnDeliveryDelayed}, states.ReturnDeliveryDelayed}, {[]states.IEnumState{states.ReturnDelivered}, states.ReturnDelivered}}
	buyerStatesMap[ReturnDeliveryFailedFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryFailed}, states.PayToSeller}}

	buyerAllStatesMapping := make(map[string][]states.IEnumState, 16)
	buyerAllStatesMapping[states.NewOrder.StateName()] = []states.IEnumState{states.NewOrder}
	buyerAllStatesMapping[states.PaymentPending.StateName()] = []states.IEnumState{states.PaymentPending}
	buyerAllStatesMapping[states.PaymentSuccess.StateName()] = []states.IEnumState{states.PaymentSuccess}
	buyerAllStatesMapping[states.PaymentFailed.StateName()] = []states.IEnumState{states.PaymentFailed}
	buyerAllStatesMapping[states.OrderVerificationPending.StateName()] = []states.IEnumState{states.OrderVerificationPending}
	buyerAllStatesMapping[states.OrderVerificationSuccess.StateName()] = []states.IEnumState{states.OrderVerificationSuccess}
	buyerAllStatesMapping[states.OrderVerificationFailed.StateName()] = []states.IEnumState{states.PayToBuyer}
	buyerAllStatesMapping[states.ApprovalPending.StateName()] = []states.IEnumState{states.ApprovalPending}
	buyerAllStatesMapping[states.ShipmentPending.StateName()] = []states.IEnumState{states.ShipmentPending}
	buyerAllStatesMapping[states.ShipmentDelayed.StateName()] = []states.IEnumState{states.ShipmentDelayed}
	buyerAllStatesMapping[states.Shipped.StateName()] = []states.IEnumState{states.Shipped}
	buyerAllStatesMapping[states.DeliveryPending.StateName()] = []states.IEnumState{states.DeliveryPending}
	buyerAllStatesMapping[states.DeliveryDelayed.StateName()] = []states.IEnumState{states.DeliveryDelayed}
	buyerAllStatesMapping[states.Delivered.StateName()] = []states.IEnumState{states.Delivered}
	buyerAllStatesMapping[states.PayToBuyer.StateName()] = []states.IEnumState{states.PayToBuyer}

	buyerReturnStatesMapping := make(map[string][]states.IEnumState, 16)
	buyerReturnStatesMapping[states.ReturnRequestPending.StateName()] = []states.IEnumState{states.ReturnRequestPending}
	buyerReturnStatesMapping[states.ReturnRequestRejected.StateName()] = []states.IEnumState{states.ReturnRequestRejected}
	buyerReturnStatesMapping[states.ReturnShipmentPending.StateName()] = []states.IEnumState{states.ReturnShipmentPending}
	buyerReturnStatesMapping[states.ReturnShipped.StateName()] = []states.IEnumState{states.ReturnShipped}
	buyerReturnStatesMapping[states.ReturnDeliveryPending.StateName()] = []states.IEnumState{states.ReturnDeliveryPending}
	buyerReturnStatesMapping[states.ReturnDeliveryDelayed.StateName()] = []states.IEnumState{states.ReturnDeliveryDelayed}
	buyerReturnStatesMapping[states.ReturnDelivered.StateName()] = []states.IEnumState{states.ReturnDelivered}
	buyerReturnStatesMapping[states.ReturnRejected.StateName()] = []states.IEnumState{states.ReturnRejected}
	buyerReturnStatesMapping[states.PayToBuyer.StateName()] = []states.IEnumState{states.ReturnRejected, states.ReturnDelivered}

	operatorFilterStatesMap := make(map[FilterValue][]FilterState, 30)
	operatorFilterStatesMap[NewOrderFilter] = []FilterState{{[]states.IEnumState{states.NewOrder}, states.NewOrder}}
	operatorFilterStatesMap[PaymentPendingFilter] = []FilterState{{[]states.IEnumState{states.PaymentPending}, states.PaymentPending}}
	operatorFilterStatesMap[PaymentSuccessFilter] = []FilterState{{[]states.IEnumState{states.PaymentSuccess}, states.ApprovalPending}}
	operatorFilterStatesMap[PaymentFailedFilter] = []FilterState{{[]states.IEnumState{states.PaymentFailed}, states.PaymentFailed}}
	operatorFilterStatesMap[OrderVerificationPendingFilter] = []FilterState{{[]states.IEnumState{states.OrderVerificationPending}, states.ApprovalPending}}
	operatorFilterStatesMap[OrderVerificationSuccessFilter] = []FilterState{{[]states.IEnumState{states.OrderVerificationSuccess}, states.ApprovalPending}}
	operatorFilterStatesMap[OrderVerificationFailedFilter] = []FilterState{{[]states.IEnumState{states.OrderVerificationFailed}, states.PayToBuyer}}
	operatorFilterStatesMap[ApprovalPendingFilter] = []FilterState{{[]states.IEnumState{states.ApprovalPending}, states.ApprovalPending}}
	operatorFilterStatesMap[CanceledBySellerFilter] = []FilterState{{[]states.IEnumState{states.CanceledBySeller}, states.PayToBuyer}}
	operatorFilterStatesMap[CanceledByBuyerFilter] = []FilterState{{[]states.IEnumState{states.CanceledByBuyer}, states.PayToBuyer}}
	operatorFilterStatesMap[ShipmentPendingFilter] = []FilterState{{[]states.IEnumState{states.ShipmentPending}, states.ShipmentPending}}
	operatorFilterStatesMap[ShipmentDelayedFilter] = []FilterState{{[]states.IEnumState{states.ShipmentDelayed}, states.ShipmentDelayed}}
	operatorFilterStatesMap[ShippedFilter] = []FilterState{{[]states.IEnumState{states.Shipped}, states.Shipped}}
	operatorFilterStatesMap[DeliveryPendingFilter] = []FilterState{{[]states.IEnumState{states.DeliveryPending}, states.DeliveryPending}}
	operatorFilterStatesMap[DeliveryDelayedFilter] = []FilterState{{[]states.IEnumState{states.DeliveryDelayed}, states.DeliveryDelayed}}
	operatorFilterStatesMap[DeliveredFilter] = []FilterState{{[]states.IEnumState{states.Delivered}, states.Delivered}}
	operatorFilterStatesMap[DeliveryFailedFilter] = []FilterState{{[]states.IEnumState{states.DeliveryFailed}, states.PayToBuyer}}
	operatorFilterStatesMap[ReturnRequestPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnRequestPending}, states.ReturnRequestPending}}
	operatorFilterStatesMap[ReturnRequestRejectedFilter] = []FilterState{{[]states.IEnumState{states.ReturnRequestRejected}, states.ReturnRequestRejected}}
	operatorFilterStatesMap[ReturnCanceledFilter] = []FilterState{{[]states.IEnumState{states.ReturnCanceled}, states.PayToSeller}}
	operatorFilterStatesMap[ReturnShipmentPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnShipmentPending}, states.ReturnShipmentPending}}
	operatorFilterStatesMap[ReturnShippedFilter] = []FilterState{{[]states.IEnumState{states.ReturnShipped}, states.ReturnShipped}}
	operatorFilterStatesMap[ReturnDeliveryPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryPending}, states.ReturnDeliveryPending}}
	operatorFilterStatesMap[ReturnDeliveryDelayedFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryDelayed}, states.ReturnDeliveryDelayed}}
	operatorFilterStatesMap[ReturnDeliveredFilter] = []FilterState{{[]states.IEnumState{states.ReturnDelivered}, states.ReturnDelivered}}
	operatorFilterStatesMap[ReturnDeliveryFailedFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryFailed}, states.PayToSeller}}
	operatorFilterStatesMap[ReturnRejectedFilter] = []FilterState{{[]states.IEnumState{states.ReturnRejected}, states.ReturnRejected}}
	operatorFilterStatesMap[PayToBuyerFilter] = []FilterState{{[]states.IEnumState{states.PayToBuyer}, states.PayToBuyer}}
	operatorFilterStatesMap[PayToSellerFilter] = []FilterState{{[]states.IEnumState{states.PayToSeller}, states.PayToSeller}}

	sellerFilterStatesMap := make(map[FilterValue][]FilterState, 30)
	sellerFilterStatesMap[ApprovalPendingFilter] = []FilterState{{[]states.IEnumState{states.ApprovalPending}, states.ApprovalPending}}
	sellerFilterStatesMap[CanceledBySellerFilter] = []FilterState{{[]states.IEnumState{states.CanceledBySeller}, states.PayToBuyer}}
	sellerFilterStatesMap[CanceledByBuyerFilter] = []FilterState{{[]states.IEnumState{states.CanceledByBuyer}, states.PayToBuyer}}
	sellerFilterStatesMap[ShipmentPendingFilter] = []FilterState{{[]states.IEnumState{states.ShipmentPending}, states.ShipmentPending}}
	sellerFilterStatesMap[ShipmentDelayedFilter] = []FilterState{{[]states.IEnumState{states.ShipmentDelayed}, states.ShipmentDelayed}}
	sellerFilterStatesMap[ShippedFilter] = []FilterState{{[]states.IEnumState{states.Shipped}, states.Shipped}}
	sellerFilterStatesMap[DeliveryPendingFilter] = []FilterState{{[]states.IEnumState{states.DeliveryPending}, states.DeliveryPending}}
	sellerFilterStatesMap[DeliveryDelayedFilter] = []FilterState{{[]states.IEnumState{states.DeliveryDelayed}, states.DeliveryDelayed}}
	sellerFilterStatesMap[DeliveredFilter] = []FilterState{{[]states.IEnumState{states.Delivered}, states.Delivered}}
	sellerFilterStatesMap[DeliveryFailedFilter] = []FilterState{{[]states.IEnumState{states.DeliveryFailed}, states.PayToBuyer}}
	sellerFilterStatesMap[ReturnRequestPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnRequestPending}, states.ReturnRequestPending}}
	sellerFilterStatesMap[ReturnRequestRejectedFilter] = []FilterState{{[]states.IEnumState{states.ReturnRequestRejected}, states.ReturnRequestRejected}}
	sellerFilterStatesMap[ReturnCanceledFilter] = []FilterState{{[]states.IEnumState{states.ReturnCanceled}, states.PayToSeller}}
	sellerFilterStatesMap[ReturnShipmentPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnShipmentPending}, states.ReturnShipmentPending}}
	sellerFilterStatesMap[ReturnShippedFilter] = []FilterState{{[]states.IEnumState{states.ReturnShipped}, states.ReturnShipped}}
	sellerFilterStatesMap[ReturnDeliveryPendingFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryPending}, states.ReturnDeliveryPending}}
	sellerFilterStatesMap[ReturnDeliveryDelayedFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryDelayed}, states.ReturnDeliveryDelayed}}
	sellerFilterStatesMap[ReturnDeliveredFilter] = []FilterState{{[]states.IEnumState{states.ReturnDelivered}, states.ReturnDelivered}}
	sellerFilterStatesMap[ReturnDeliveryFailedFilter] = []FilterState{{[]states.IEnumState{states.ReturnDeliveryFailed}, states.PayToSeller}}
	sellerFilterStatesMap[ReturnRejectedFilter] = []FilterState{{[]states.IEnumState{states.ReturnRejected}, states.ReturnRejected}}
	sellerFilterStatesMap[PayToSellerFilter] = []FilterState{{[]states.IEnumState{states.ReturnCanceled, states.ReturnDeliveryFailed, states.ReturnShipmentPending, states.ReturnRequestRejected, states.Delivered, states.ReturnRejected}, states.PayToSeller}}

	sellerStatesMapping := make(map[string][]states.IEnumState, 30)
	sellerStatesMapping[states.ApprovalPending.StateName()] = []states.IEnumState{states.ApprovalPending}
	sellerStatesMapping[states.ShipmentPending.StateName()] = []states.IEnumState{states.ShipmentPending}
	sellerStatesMapping[states.ShipmentDelayed.StateName()] = []states.IEnumState{states.ShipmentDelayed}
	sellerStatesMapping[states.Shipped.StateName()] = []states.IEnumState{states.Shipped}
	sellerStatesMapping[states.DeliveryPending.StateName()] = []states.IEnumState{states.DeliveryPending}
	sellerStatesMapping[states.DeliveryDelayed.StateName()] = []states.IEnumState{states.DeliveryDelayed}
	sellerStatesMapping[states.Delivered.StateName()] = []states.IEnumState{states.Delivered}
	sellerStatesMapping[states.ReturnRequestPending.StateName()] = []states.IEnumState{states.ReturnRequestPending}
	sellerStatesMapping[states.ReturnRequestRejected.StateName()] = []states.IEnumState{states.ReturnRequestRejected}
	sellerStatesMapping[states.ReturnShipmentPending.StateName()] = []states.IEnumState{states.ReturnShipmentPending}
	sellerStatesMapping[states.ReturnShipped.StateName()] = []states.IEnumState{states.ReturnShipped}
	sellerStatesMapping[states.ReturnDeliveryPending.StateName()] = []states.IEnumState{states.ReturnDeliveryPending}
	sellerStatesMapping[states.ReturnDeliveryDelayed.StateName()] = []states.IEnumState{states.ReturnDeliveryDelayed}
	sellerStatesMapping[states.ReturnDelivered.StateName()] = []states.IEnumState{states.ReturnDelivered}
	sellerStatesMapping[states.ReturnRejected.StateName()] = []states.IEnumState{states.ReturnRejected}
	sellerStatesMapping[states.PayToSeller.StateName()] = []states.IEnumState{states.ReturnCanceled, states.ReturnDeliveryFailed, states.ReturnShipmentPending, states.ReturnRequestRejected, states.ReturnRejected, states.Delivered}
	sellerStatesMapping[states.PayToBuyer.StateName()] = []states.IEnumState{states.CanceledBySeller, states.CanceledByBuyer, states.DeliveryFailed, states.ReturnRejected, states.ReturnDelivered}

	queryPathStatesMap := make(map[FilterValue]FilterQueryState, 30)
	queryPathStatesMap[NewOrderFilter] = FilterQueryState{states.NewOrder, "packages.subpackages.status"}
	queryPathStatesMap[PaymentPendingFilter] = FilterQueryState{states.PaymentPending, "packages.subpackages.status"}
	queryPathStatesMap[PaymentSuccessFilter] = FilterQueryState{states.PaymentSuccess, "packages.subpackages.status"}
	queryPathStatesMap[PaymentFailedFilter] = FilterQueryState{states.PaymentFailed, "packages.subpackages.status"}
	queryPathStatesMap[OrderVerificationPendingFilter] = FilterQueryState{states.OrderVerificationPending, "packages.subpackages.status"}
	queryPathStatesMap[OrderVerificationSuccessFilter] = FilterQueryState{states.OrderVerificationSuccess, "packages.subpackages.status"}
	queryPathStatesMap[OrderVerificationFailedFilter] = FilterQueryState{states.OrderVerificationFailed, "packages.subpackages.tracking.history.name"}
	queryPathStatesMap[ApprovalPendingFilter] = FilterQueryState{states.ApprovalPending, "packages.subpackages.status"}
	queryPathStatesMap[CanceledBySellerFilter] = FilterQueryState{states.CanceledBySeller, "packages.subpackages.tracking.history.name"}
	queryPathStatesMap[CanceledByBuyerFilter] = FilterQueryState{states.CanceledByBuyer, "packages.subpackages.tracking.history.name"}
	queryPathStatesMap[ShipmentPendingFilter] = FilterQueryState{states.ShipmentPending, "packages.subpackages.status"}
	queryPathStatesMap[ShipmentDelayedFilter] = FilterQueryState{states.ShipmentDelayed, "packages.subpackages.status"}
	queryPathStatesMap[ShippedFilter] = FilterQueryState{states.Shipped, "packages.subpackages.status"}
	queryPathStatesMap[DeliveryPendingFilter] = FilterQueryState{states.DeliveryPending, "packages.subpackages.status"}
	queryPathStatesMap[DeliveryDelayedFilter] = FilterQueryState{states.DeliveryDelayed, "packages.subpackages.status"}
	queryPathStatesMap[DeliveredFilter] = FilterQueryState{states.Delivered, "packages.subpackages.status"}
	queryPathStatesMap[DeliveryFailedFilter] = FilterQueryState{states.DeliveryFailed, "packages.subpackages.tracking.history.name"}
	queryPathStatesMap[ReturnRequestPendingFilter] = FilterQueryState{states.ReturnRequestPending, "packages.subpackages.status"}
	queryPathStatesMap[ReturnRequestRejectedFilter] = FilterQueryState{states.ReturnRequestRejected, "packages.subpackages.status"}
	queryPathStatesMap[ReturnCanceledFilter] = FilterQueryState{states.ReturnCanceled, "packages.subpackages.tracking.history.name"}
	queryPathStatesMap[ReturnShipmentPendingFilter] = FilterQueryState{states.ReturnShipmentPending, "packages.subpackages.status"}
	queryPathStatesMap[ReturnShippedFilter] = FilterQueryState{states.ReturnShipped, "packages.subpackages.status"}
	queryPathStatesMap[ReturnDeliveryPendingFilter] = FilterQueryState{states.ReturnDeliveryPending, "packages.subpackages.status"}
	queryPathStatesMap[ReturnDeliveryDelayedFilter] = FilterQueryState{states.ReturnDeliveryDelayed, "packages.subpackages.status"}
	queryPathStatesMap[ReturnDeliveredFilter] = FilterQueryState{states.ReturnDelivered, "packages.subpackages.status"}
	queryPathStatesMap[ReturnDeliveryFailedFilter] = FilterQueryState{states.ReturnDeliveryFailed, "packages.subpackages.tracking.history.name"}
	queryPathStatesMap[ReturnRejectedFilter] = FilterQueryState{states.ReturnRejected, "packages.subpackages.status"}
	queryPathStatesMap[PayToBuyerFilter] = FilterQueryState{states.PayToBuyer, "packages.subpackages.status"}
	queryPathStatesMap[PayToSellerFilter] = FilterQueryState{states.PayToSeller, "packages.subpackages.status"}

	actionStateMap := make(map[UserType][]actions.IAction, 8)
	actionStateMap[SellerUser] = []actions.IAction{
		seller_action.New(seller_action.Approve),
		seller_action.New(seller_action.Reject),
		seller_action.New(seller_action.Cancel),
		seller_action.New(seller_action.Accept),
		seller_action.New(seller_action.Deliver),
		seller_action.New(seller_action.DeliveryFail),
		seller_action.New(seller_action.EnterShipmentDetail),
	}
	actionStateMap[BuyerUser] = []actions.IAction{
		buyer_action.New(buyer_action.DeliveryDelay),
		buyer_action.New(buyer_action.Cancel),
		buyer_action.New(buyer_action.SubmitReturnRequest),
		buyer_action.New(buyer_action.EnterShipmentDetail),
	}
	actionStateMap[OperatorUser] = []actions.IAction{
		operator_action.New(operator_action.DeliveryDelay),
		operator_action.New(operator_action.Deliver),
		operator_action.New(operator_action.DeliveryFail),
		operator_action.New(operator_action.Accept),
		operator_action.New(operator_action.Reject),
		operator_action.New(operator_action.Deliver),
	}
	actionStateMap[SchedulerUser] = []actions.IAction{
		scheduler_action.New(scheduler_action.Cancel),
		scheduler_action.New(scheduler_action.Close),
		scheduler_action.New(scheduler_action.PaymentFail),
		scheduler_action.New(scheduler_action.DeliveryDelay),
		scheduler_action.New(scheduler_action.Deliver),
		scheduler_action.New(scheduler_action.DeliveryPending),
		scheduler_action.New(scheduler_action.Reject),
		scheduler_action.New(scheduler_action.Accept),
		scheduler_action.New(scheduler_action.Notification),
	}

	reqFilters := make(map[RequestName][]FilterValue, 8)
	reqFilters[SellerOrderList] = []FilterValue{
		ApprovalPendingFilter,
		CanceledBySellerFilter,
		CanceledByBuyerFilter,
		ShipmentPendingFilter,
		ShipmentDelayedFilter,
		ShippedFilter,
		DeliveryPendingFilter,
		DeliveryDelayedFilter,
		DeliveredFilter,
		DeliveryFailedFilter,
		AllCanceledFilter,
		AllOrdersFilter,
	}

	reqFilters[SellerOrderDetail] = []FilterValue{
		ApprovalPendingFilter,
		CanceledBySellerFilter,
		CanceledByBuyerFilter,
		ShipmentPendingFilter,
		ShipmentDelayedFilter,
		ShippedFilter,
		DeliveryPendingFilter,
		DeliveredFilter,
		DeliveryFailedFilter,
		AllCanceledFilter,
		AllOrdersFilter,
	}

	reqFilters[SellerReturnOrderList] = []FilterValue{
		ReturnRequestPendingFilter,
		ReturnRequestRejectedFilter,
		ReturnCanceledFilter,
		ReturnShipmentPendingFilter,
		ReturnShippedFilter,
		ReturnDeliveryPendingFilter,
		ReturnDeliveryDelayedFilter,
		ReturnDeliveredFilter,
		ReturnDeliveryFailedFilter,
		ReturnRejectedFilter,
	}

	reqFilters[SellerReturnOrderDetail] = []FilterValue{
		ReturnRequestPendingFilter,
		ReturnRequestRejectedFilter,
		ReturnCanceledFilter,
		ReturnShipmentPendingFilter,
		ReturnShippedFilter,
		ReturnDeliveryPendingFilter,
		ReturnDeliveryDelayedFilter,
		ReturnDeliveredFilter,
		ReturnDeliveryFailedFilter,
		ReturnRejectedFilter,
	}

	reqFilters[SellerOrderDashboardReports] = []FilterValue{}
	reqFilters[SellerOrderShipmentReports] = []FilterValue{}
	reqFilters[SellerOrderDeliveredReports] = []FilterValue{}
	reqFilters[SellerOrderReturnReports] = []FilterValue{}
	reqFilters[SellerOrderCancelReports] = []FilterValue{}
	reqFilters[SellerApprovalPendingOrderReports] = []FilterValue{}
	reqFilters[SellerAllOrderReports] = []FilterValue{}

	reqFilters[BuyerOrderDetailList] = []FilterValue{
		NewOrderFilter,
		PaymentPendingFilter,
		PaymentSuccessFilter,
		PaymentFailedFilter,
		OrderVerificationPendingFilter,
		OrderVerificationSuccessFilter,
		OrderVerificationFailedFilter,
		ApprovalPendingFilter,
		CanceledBySellerFilter,
		CanceledByBuyerFilter,
		ShipmentPendingFilter,
		ShipmentDelayedFilter,
		ShippedFilter,
		DeliveryPendingFilter,
		DeliveryDelayedFilter,
		DeliveredFilter,
		DeliveryFailedFilter,
		PayToBuyerFilter,
		AllOrdersFilter,
	}

	//reqFilters[BuyerReturnOrderReports] = []FilterValue{}

	reqFilters[BuyerReturnOrderDetailList] = []FilterValue{
		ReturnRequestPendingFilter,
		ReturnShipmentPendingFilter,
		ReturnShippedFilter,
		ReturnDeliveredFilter,
		ReturnDeliveryFailedFilter,
		AllOrdersFilter,
	}

	//reqFilters[BuyerAllReturnOrders] = []FilterValue{
	//	ReturnRequestPendingFilter,
	//	ReturnRequestRejectedFilter,
	//	ReturnCanceledFilter,
	//	ReturnShipmentPendingFilter,
	//	ReturnShippedFilter,
	//	ReturnDeliveryPendingFilter,
	//	ReturnDeliveryDelayedFilter,
	//	ReturnDeliveredFilter,
	//	ReturnDeliveryFailedFilter,
	//	ReturnRejectedFilter,
	//}
	rp := createReasonsMap()
	return Server{
		flowManager:          flowManager,
		address:              address,
		port:                 port,
		requestFilters:       reqFilters,
		buyerFilterStates:    buyerStatesMap,
		buyerAllStatesMap:    buyerAllStatesMapping,
		buyerReturnStatesMap: buyerReturnStatesMapping,
		sellerFilterStates:   sellerFilterStatesMap,
		sellerStatesMap:      sellerStatesMapping,
		operatorFilterStates: operatorFilterStatesMap,
		queryPathStates:      queryPathStatesMap,
		actionStates:         actionStateMap,
		reasonsMap:           rp,
	}
}

func (server *Server) RequestHandler(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {

	userAcl, err := app.Globals.UserService.AuthenticateContextToken(ctx)
	if err != nil {
		app.Globals.Logger.FromContext(ctx).Error("UserService.AuthenticateContextToken failed", "fn", "RequestHandler", "error", err)
		return nil, status.Error(codes.Code(future.Forbidden), "User Not Authorized")
	}

	if uint64(userAcl.User().UserID) != req.Meta.UID {
		app.Globals.Logger.FromContext(ctx).Error("request userId mismatch with token userId", "fn", "RequestHandler",
			"userId", req.Meta.UID, "token", userAcl.User().UserID)
		return nil, status.Error(codes.Code(future.Forbidden), "User Not Authorized")
	}

	if req.Meta.UTP == string(OperatorUser) {
		if !userAcl.UserPerm().Has("order.state.all.view") && RequestType(req.Type) == DataReqType {
			return nil, status.Error(codes.Code(future.Forbidden), "User Not Permitted")
		}

		if !userAcl.UserPerm().Has("order.state.all.action") && RequestType(req.Type) == ActionReqType {
			return nil, status.Error(codes.Code(future.Forbidden), "User Not Permitted")
		}
	}

	if ctx.Value(string(utils.CtxUserID)) == nil {
		ctx = context.WithValue(ctx, string(utils.CtxUserID), uint64(req.Meta.UID))
		ctx = context.WithValue(ctx, string(utils.CtxUserACL), userAcl)
	}

	reqType := RequestType(req.Type)
	if reqType == DataReqType {
		return server.requestDataHandler(ctx, req)
	} else {
		return server.requestActionHandler(ctx, req)
	}
}

func (server *Server) SchedulerMessageHandler(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {

	app.Globals.Logger.FromContext(ctx).Debug("Received scheduler action request",
		"fn", "SchedulerMessageHandler",
		"request", req)

	if ctx.Value(string(utils.CtxUserID)) == nil {
		ctx = context.WithValue(ctx, string(utils.CtxUserID), uint64(0))
	}

	userType := SchedulerUser
	var userAction actions.IAction

	var schedulerActionRequest pb.SchedulerActionRequest
	if err := ptypes.UnmarshalAny(req.Data, &schedulerActionRequest); err != nil {
		app.Globals.Logger.Error("Could not unmarshal schedulerActionRequest from request anything field", "fn", "SchedulerMessageHandler",
			"request", req, "error", err)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	for _, orderReq := range schedulerActionRequest.Orders {
		userActions, ok := server.actionStates[userType]
		if !ok {
			app.Globals.Logger.Error("requested scheduler action not supported", "fn", "SchedulerMessageHandler", "request", req)
			return nil, status.Error(codes.Code(future.BadRequest), "Scheduler Action Invalid")
		}

		for _, action := range userActions {
			if action.ActionEnum().ActionName() == orderReq.ActionState {
				userAction = action
				break
			}
		}

		if userAction == nil {
			app.Globals.Logger.Error("scheduler action invalid", "fn", "SchedulerMessageHandler", "request", req)
			return nil, status.Error(codes.Code(future.BadRequest), "Action Invalid")
		}

		if userAction.ActionEnum() == scheduler_action.PaymentFail {
			event := events.New(events.Action, orderReq.OID, 0, 0,
				orderReq.StateIndex, userAction,
				time.Unix(req.Time.GetSeconds(), int64(req.Time.GetNanos())), nil)

			app.Globals.Logger.Debug("scheduler action event paymentFail",
				"fn", "SchedulerMessageHandler",
				"oid", event.OrderId(),
				"uid", event.UserId(),
				"event", event)

			iFuture := future.Factory().SetCapacity(1).Build()
			iFrame := frame.Factory().SetFuture(iFuture).SetEvent(event).Build()
			server.flowManager.MessageHandler(ctx, iFrame)
			futureData := iFuture.Get()
			if futureData.Error() != nil {
				app.Globals.Logger.FromContext(ctx).Error("flowManager.MessageHandler failed",
					"fn", "SchedulerMessageHandler",
					"event", event,
					"error", futureData.Error().Reason())
			}

		} else {
			for _, pkgReq := range orderReq.Packages {
				subpackages := make([]events.ActionSubpackage, 0, len(pkgReq.Subpackages))
				for _, subPkgReq := range pkgReq.Subpackages {

					subpackage := events.ActionSubpackage{
						SId:   subPkgReq.SID,
						Items: nil,
					}
					subpackage.Items = make([]events.ActionItem, 0, len(subPkgReq.Items))
					for _, item := range subPkgReq.Items {
						actionItem := events.ActionItem{
							InventoryId: item.InventoryId,
							Quantity:    item.Quantity,
						}
						subpackage.Items = append(subpackage.Items, actionItem)
					}
					subpackages = append(subpackages, subpackage)
				}

				actionData := events.ActionData{
					SubPackages:    subpackages,
					Carrier:        "",
					TrackingNumber: "",
				}

				event := events.New(events.Action, orderReq.OID, pkgReq.PID, 0,
					orderReq.StateIndex, userAction,
					time.Unix(req.Time.GetSeconds(), int64(req.Time.GetNanos())), actionData)

				iFuture := future.Factory().SetCapacity(1).Build()
				iFrame := frame.Factory().SetFuture(iFuture).SetEvent(event).Build()

				app.Globals.Logger.Debug("scheduler action event",
					"fn", "SchedulerMessageHandler",
					"oid", event.OrderId(),
					"uid", event.UserId(),
					"event", event)

				server.flowManager.MessageHandler(ctx, iFrame)
				futureData := iFuture.Get()
				if futureData.Error() != nil {
					app.Globals.Logger.Error("flowManager.MessageHandler failed", "fn", "SchedulerMessageHandler", "event", event, "error", futureData.Error().Reason())
				}
			}
		}
	}

	response := &pb.MessageResponse{
		Entity: "ActionResponse",
		Meta:   nil,
		Data:   nil,
	}
	return response, nil
}

func (server *Server) FinanceOrderItems(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	app.Globals.Logger.FromContext(ctx).Debug("Received financeOrderItem request",
		"fn", "FinanceOrderItems",
		"request", req)

	if ctx.Value(string(utils.CtxUserID)) == nil {
		ctx = context.WithValue(ctx, string(utils.CtxUserID), uint64(0))
	}

	var state string
	var filterValue FilterValue
	var sortName string
	var sortDirection int
	var startTimestamp time.Time
	var endTimestamp time.Time

	if req.Meta.Filters != nil {
		filterValue = FilterValue(req.Meta.Filters[0].Value)
		if filterValue != PayToSellerFilter && filterValue != PayToBuyerFilter {
			app.Globals.Logger.Error("filterValue invalid",
				"fn", "FinanceOrderItems",
				"filterValue", filterValue,
				"request", req)
			return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
		}

		if filterValue == PayToSellerFilter {
			state = states.PayToSeller.StateName()
		} else if filterValue == PayToBuyerFilter {
			state = states.PayToBuyer.StateName()
		}
	} else {
		app.Globals.Logger.Error("filters is empty",
			"fn", "FinanceOrderItems",
			"request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	if req.Meta.StartTimestamp != "" {
		temp, err := time.Parse(ISO8601, req.Meta.StartTimestamp)
		if err != nil {
			app.Globals.Logger.Error("StartTimestamp invalid",
				"fn", "FinanceOrderItems",
				"StartTimestamp", req.Meta.StartTimestamp,
				"request", req)
			return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
		}

		startTimestamp = temp
	} else {
		app.Globals.Logger.Error("StartTimestamp is empty",
			"fn", "FinanceOrderItems",
			"request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	if req.Meta.EndTimestamp != "" {
		temp, err := time.Parse(ISO8601, req.Meta.EndTimestamp)
		if err != nil {
			app.Globals.Logger.Error("EndTimestamp invalid",
				"fn", "FinanceOrderItems",
				"endTimestamp", req.Meta.EndTimestamp,
				"request", req)
			return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
		}

		endTimestamp = temp
	} else {
		app.Globals.Logger.Error("EndTimestamp is empty",
			"fn", "FinanceOrderItems",
			"request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	if req.Meta.Sorts != nil {
		sortName = req.Meta.Sorts[0].Name
		if SortDirection(req.Meta.Sorts[0].Direction) == ASC {
			sortDirection = 1
		} else {
			sortDirection = -1
		}
	}

	if req.Meta.Page <= 0 {
		app.Globals.Logger.Error("Page invalid",
			"fn", "FinanceOrderItems",
			"page", req.Meta.Page,
			"request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	if req.Meta.PerPage <= 0 {
		app.Globals.Logger.Error("PerPage invalid",
			"fn", "FinanceOrderItems",
			"perPage", req.Meta.Page,
			"request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	finances, total, err := app.Globals.FinanceReportRepository.FindAllWithPageAndSort(ctx, state, startTimestamp,
		endTimestamp, int64(req.Meta.Page), int64(req.Meta.PerPage), sortName, sortDirection)

	if err != nil {
		app.Globals.Logger.Error("PerPage invalid",
			"fn", "FinanceOrderItems",
			"perPage", req.Meta.PerPage,
			"request", req,
			"error", err.Reason())
		return nil, status.Error(codes.Code(err.Code()), err.Message())
	}

	financeReportList := make([]*pb.FinanceOrderItemDetailList_OrderItemDetail, 0, len(finances))
	for _, finance := range finances {
		if finance.RawSellerShippingNet == nil || finance.RoundupSellerShippingNet == nil {
			app.Globals.Logger.Info("OrderItem doesn't support finance",
				"fn", "FinanceOrderItems",
				"oid", finance.OrderId,
				"pid", finance.PId,
				"sid", finance.SId)
			continue
		}

		financeOrderItem := &pb.FinanceOrderItemDetailList_OrderItemDetail{
			OId:      finance.OrderId,
			SellerId: finance.PId,
			ShipmentAmount: &pb.Money{
				Amount:   finance.ShipmentAmount.Amount,
				Currency: finance.ShipmentAmount.Currency,
			},
			RawShippingNet: &pb.Money{
				Amount:   finance.RawSellerShippingNet.Amount,
				Currency: finance.RawSellerShippingNet.Currency,
			},
			RoundupShippingNet: &pb.Money{
				Amount:   finance.RoundupSellerShippingNet.Amount,
				Currency: finance.RoundupSellerShippingNet.Currency,
			},
			Items:          nil,
			CreatedAt:      finance.CreatedAt.Format(ISO8601),
			UpdatedAt:      finance.UpdatedAt.Format(ISO8601),
			OrderCreatedAt: finance.CreatedAt.Format(ISO8601),
		}

		financeItemList := make([]*pb.FinanceOrderItemDetailList_OrderItemDetail_Item, 0, len(finance.Items))
		for _, item := range finance.Items {
			financeItem := &pb.FinanceOrderItemDetailList_OrderItemDetail_Item{
				SId:         finance.SId,
				Sku:         item.SKU,
				InventoryId: item.InventoryId,
				Title:       item.Title,
				Brand:       item.Brand,
				Category:    item.Category,
				Guaranty:    item.Guaranty,
				Image:       item.Image,
				Returnable:  item.Returnable,
				Quantity:    item.Quantity,
				Attributes:  nil,
				Invoice: &pb.FinanceOrderItemDetailList_OrderItemDetail_Item_ItemInvoice{
					Commission: &pb.FinanceOrderItemDetailList_OrderItemDetail_Item_ItemInvoice_ItemCommission{
						ItemCommission: item.Invoice.Commission.ItemCommission,
						RawUnitPrice: &pb.Money{
							Amount:   item.Invoice.Commission.RawUnitPrice.Amount,
							Currency: item.Invoice.Commission.RawUnitPrice.Currency,
						},
						RoundupUnitPrice: &pb.Money{
							Amount:   item.Invoice.Commission.RoundupUnitPrice.Amount,
							Currency: item.Invoice.Commission.RoundupUnitPrice.Currency,
						},
						RawTotalPrice: &pb.Money{
							Amount:   item.Invoice.Commission.RawTotalPrice.Amount,
							Currency: item.Invoice.Commission.RawTotalPrice.Currency,
						},
						RoundupTotalPrice: &pb.Money{
							Amount:   item.Invoice.Commission.RoundupUnitPrice.Amount,
							Currency: item.Invoice.Commission.RoundupUnitPrice.Currency,
						},
					},
					Share: &pb.FinanceOrderItemDetailList_OrderItemDetail_Item_ItemInvoice_ItemShare{
						RawItemNet: &pb.Money{
							Amount:   item.Invoice.Share.RawItemNet.Amount,
							Currency: item.Invoice.Share.RawItemNet.Currency,
						},
						RoundupItemNet: &pb.Money{
							Amount:   item.Invoice.Share.RoundupItemNet.Amount,
							Currency: item.Invoice.Share.RoundupItemNet.Currency,
						},
						RawTotalNet: &pb.Money{
							Amount:   item.Invoice.Share.RawTotalNet.Amount,
							Currency: item.Invoice.Share.RawTotalNet.Currency,
						},
						RoundupTotalNet: &pb.Money{
							Amount:   item.Invoice.Share.RoundupTotalNet.Amount,
							Currency: item.Invoice.Share.RoundupTotalNet.Currency,
						},
						RawUnitSellerShare: &pb.Money{
							Amount:   item.Invoice.Share.RawUnitSellerShare.Amount,
							Currency: item.Invoice.Share.RawUnitSellerShare.Currency,
						},
						RoundupUnitSellerShare: &pb.Money{
							Amount:   item.Invoice.Share.RoundupUnitSellerShare.Amount,
							Currency: item.Invoice.Share.RoundupUnitSellerShare.Currency,
						},
						RawTotalSellerShare: &pb.Money{
							Amount:   item.Invoice.Share.RawTotalSellerShare.Amount,
							Currency: item.Invoice.Share.RawTotalSellerShare.Currency,
						},
						RoundupTotalSellerShare: &pb.Money{
							Amount:   item.Invoice.Share.RoundupTotalSellerShare.Amount,
							Currency: item.Invoice.Share.RoundupTotalSellerShare.Currency,
						},
					},
					SSO: &pb.FinanceOrderItemDetailList_OrderItemDetail_Item_ItemInvoice_ItemSSO{
						Rate:      item.Invoice.SSO.Rate,
						IsObliged: item.Invoice.SSO.IsObliged,
						RawUnitPrice: &pb.Money{
							Amount:   item.Invoice.SSO.RawUnitPrice.Amount,
							Currency: item.Invoice.SSO.RawUnitPrice.Currency,
						},
						RoundupUnitPrice: &pb.Money{
							Amount:   item.Invoice.SSO.RoundupUnitPrice.Amount,
							Currency: item.Invoice.SSO.RoundupUnitPrice.Currency,
						},
						RawTotalPrice: &pb.Money{
							Amount:   item.Invoice.SSO.RawTotalPrice.Amount,
							Currency: item.Invoice.SSO.RawTotalPrice.Currency,
						},
						RoundupTotalPrice: &pb.Money{
							Amount:   item.Invoice.SSO.RoundupTotalPrice.Amount,
							Currency: item.Invoice.SSO.RoundupTotalPrice.Currency,
						},
					},
					VAT: &pb.FinanceOrderItemDetailList_OrderItemDetail_Item_ItemInvoice_ItemVAT{
						Rate:      item.Invoice.VAT.SellerVat.Rate,
						IsObliged: item.Invoice.VAT.SellerVat.IsObliged,
						RawUnitPrice: &pb.Money{
							Amount:   item.Invoice.VAT.SellerVat.RawUnitPrice.Amount,
							Currency: item.Invoice.VAT.SellerVat.RawUnitPrice.Currency,
						},
						RoundupUnitPrice: &pb.Money{
							Amount:   item.Invoice.VAT.SellerVat.RoundupUnitPrice.Amount,
							Currency: item.Invoice.VAT.SellerVat.RoundupUnitPrice.Currency,
						},
						RawTotalPrice: &pb.Money{
							Amount:   item.Invoice.VAT.SellerVat.RawTotalPrice.Amount,
							Currency: item.Invoice.VAT.SellerVat.RawTotalPrice.Currency,
						},
						RoundupTotalPrice: &pb.Money{
							Amount:   item.Invoice.VAT.SellerVat.RoundupTotalPrice.Amount,
							Currency: item.Invoice.VAT.SellerVat.RoundupTotalPrice.Currency,
						},
					},
				},
			}

			if item.Attributes != nil {
				financeItem.Attributes = make(map[string]*pb.Attribute, len(item.Attributes))
				for attrKey, attribute := range item.Attributes {
					keyTranslates := make(map[string]string, len(attribute.KeyTranslate))
					for keyTran, value := range attribute.KeyTranslate {
						keyTranslates[keyTran] = value
					}

					valTranslates := make(map[string]string, len(attribute.ValueTranslate))
					for valTran, value := range attribute.ValueTranslate {
						valTranslates[valTran] = value
					}

					financeItem.Attributes[attrKey] = &pb.Attribute{
						KeyTrans:   keyTranslates,
						ValueTrans: valTranslates,
					}
				}
			}

			financeItemList = append(financeItemList, financeItem)
		}

		financeReportList = append(financeReportList, financeOrderItem)
	}

	financeOrderItemDetailList := &pb.FinanceOrderItemDetailList{
		OrderItems: financeReportList,
	}

	serializedData, e := proto.Marshal(financeOrderItemDetailList)
	if e != nil {
		app.Globals.Logger.FromContext(ctx).Error("could not marshal financeOrderItemDetailList",
			"fn", "FinanceOrderItems",
			"request", req,
			"error", e)
		return nil, status.Error(codes.Code(future.InternalError), "Unknown Error")
	}

	response := &pb.MessageResponse{
		Entity: "financeOrderItemDetailList",
		Meta: &pb.ResponseMetadata{
			Total:   uint32(total),
			Page:    req.Meta.Page,
			PerPage: req.Meta.PerPage,
		},
		Data: &any.Any{
			TypeUrl: "baman.io/" + proto.MessageName(financeOrderItemDetailList),
			Value:   serializedData,
		},
	}

	return response, nil
}

func (server *Server) requestDataHandler(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	reqName := RequestName(req.Name)
	userType := UserType(req.Meta.UTP)
	//reqADT := RequestADT(req.ADT)

	//var filterType FilterType
	var filterValue FilterValue
	var sortName string
	var sortDirection SortDirection
	var buyerMobile string
	if req.Meta.Filters != nil {
		//filterType = FilterType(req.Meta.Filters[0].UTP)
		filterValue = FilterValue(req.Meta.Filters[0].Value)
	}

	if req.Meta.Sorts != nil {
		sortName = req.Meta.Sorts[0].Name
		sortDirection = SortDirection(req.Meta.Sorts[0].Direction)
	}

	//if reqName == SellerOrderList && filterType != OrderStateFilterType {
	//	logger.Err("requestDataHandler() => request name %s mismatch with %s filter, request: %v", reqName, filterType, req)
	//	return nil, status.Error(codes.Code(future.BadRequest), "Mismatch Request name with filter")
	//}

	//if (reqName == SellerReturnOrderList || reqName == BuyerReturnOrderDetailList) && filterType != OrderReturnStateFilter {
	//	logger.Err("requestDataHandler() => request name %s mismatch with %s filterType, request: %v", reqName, filterType, req)
	//	return nil, status.Error(codes.Code(future.BadRequest), "Mismatch Request name with filterType")
	//}

	if userType == SellerUser &&
		reqName != SellerOrderList &&
		reqName != SellerOrderDetail &&
		reqName != SellerReturnOrderList &&
		reqName != SellerReturnOrderDetail &&
		reqName != SellerOrderDeliveredReports &&
		reqName != SellerOrderReturnReports &&
		reqName != SellerOrderShipmentReports &&
		reqName != SellerOrderDashboardReports &&
		reqName != SellerOrderCancelReports &&
		reqName != SellerApprovalPendingOrderReports &&
		reqName != SellerAllOrderReports {
		app.Globals.Logger.FromContext(ctx).Error("RequestName with userType mismatch", "fn", "requestDataHandler", "rn", reqName, "utp", userType, "request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "RN UTP Invalid")
	} else if userType == BuyerUser &&
		reqName != BuyerOrderDetailList &&
		reqName != BuyerAllOrderReports &&
		reqName != BuyerReturnOrderReports &&
		reqName != BuyerReturnOrderDetailList {
		app.Globals.Logger.FromContext(ctx).Error("RequestName with userType mismatch", "fn", "requestDataHandler", "rn", reqName, "utp", userType, "request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "RN UTP Invalid")
	} else if userType == OperatorUser &&
		reqName != OperatorOrderList &&
		reqName != OperatorOrderDetail &&
		reqName != OperatorOrderInvoiceDetail {
		app.Globals.Logger.FromContext(ctx).Error("RequestName with userType mismatch", "fn", "requestDataHandler", "rn", reqName, "utp", userType, "request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "RN UTP Invalid")
	}

	//if req.Meta.OID > 0 && reqADT == ListType {
	//	logger.Err("requestDataHandler() => %s orderId mismatch with %s requestADT, request: %v", userType, reqADT, req)
	//	return nil, status.Error(codes.Code(future.BadRequest), "Mismatch Request name with RequestADT")
	//}

	//if req.Meta.OID > 0 && reqName != SellerOrderList && reqName != OperatorOrderList {
	//	logger.Err("requestDataHandler() => %s orderId mismatch with %s requestName, request: %v", userType, reqName, req)
	//	return nil, status.Error(codes.Code(future.BadRequest), "Mismatch OrderId with Request name")
	//}

	if userType == BuyerUser && reqName != BuyerAllOrderReports && reqName != BuyerReturnOrderReports {
		if reqName == BuyerOrderDetailList {
			if filterValue != "" {
				var findFlag = false
				for _, filter := range server.requestFilters[reqName] {
					if filter == filterValue {
						findFlag = true
						break
					}
				}

				if !findFlag && req.Meta.OID <= 0 {
					app.Globals.Logger.FromContext(ctx).Error("RequestName with filter mismatch", "fn", "requestDataHandler", "rn", reqName, "filter", filterValue, "request", req)
					return nil, status.Error(codes.Code(future.BadRequest), "RN Filter Invalid")
				}
			}
		} else {
			var findFlag = false
			for _, filter := range server.requestFilters[reqName] {
				if filter == filterValue {
					findFlag = true
					break
				}
			}

			if !findFlag {
				app.Globals.Logger.FromContext(ctx).Error("RequestName with filter mismatch", "fn", "requestDataHandler", "rn", reqName, "filter", filterValue, "request", req)
				return nil, status.Error(codes.Code(future.BadRequest), "RN Filter Invalid")
			}
		}
	} else if userType == SellerUser &&
		reqName != SellerOrderDashboardReports &&
		reqName != SellerOrderShipmentReports &&
		reqName != SellerOrderDeliveredReports &&
		reqName != SellerOrderReturnReports &&
		reqName != SellerOrderCancelReports &&
		reqName != SellerApprovalPendingOrderReports &&
		reqName != SellerAllOrderReports {
		var findFlag = false
		for _, filter := range server.requestFilters[reqName] {
			if filter == filterValue {
				findFlag = true
				break
			}
		}

		if !findFlag {
			app.Globals.Logger.FromContext(ctx).Error("RequestName with filter mismatch", "fn", "requestDataHandler", "rn", reqName, "filter", filterValue, "request", req)
			return nil, status.Error(codes.Code(future.BadRequest), "RN Filter Invalid")
		}
	}

	if reqName == OperatorOrderDetail && filterValue != "" {
		app.Globals.Logger.FromContext(ctx).Error("RequestName doesn't need any filter", "fn", "requestDataHandler", "rn", reqName, "filter", filterValue, "request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "RN Filter Invalid")
	} else if reqName == OperatorOrderList {
		if req.Meta.Ext != nil {
			buyerMobile = req.Meta.Ext["buyerMobile"]
		}
	}

	//if req.Meta.OID > 0 && reqName == SellerOrderList {
	//	return server.sellerGetOrderByIdHandler(ctx, , req.Meta.PID, filterValue)
	//}

	switch reqName {
	case SellerOrderList:
		return server.sellerOrderListHandler(ctx, req.Meta.OID, req.Meta.PID, filterValue, req.Meta.Page, req.Meta.PerPage, sortName, sortDirection)
	case SellerOrderDetail:
		return server.sellerOrderDetailHandler(ctx, req.Meta.PID, req.Meta.OID, filterValue)
	case SellerReturnOrderList:
		return server.sellerReturnOrderListHandler(ctx, req.Meta.PID, filterValue, req.Meta.Page, req.Meta.PerPage, sortName, sortDirection)
	case SellerReturnOrderDetail:
		return server.sellerReturnOrderDetailHandler(ctx, req.Meta.PID, req.Meta.OID, filterValue)

	case SellerOrderDashboardReports:
		return server.sellerOrderDashboardReportsHandler(ctx, req.Meta.UID)
	case SellerOrderShipmentReports:
		return server.sellerOrderShipmentReportsHandler(ctx, req.Meta.UID)
	case SellerOrderReturnReports:
		return server.sellerOrderReturnReportsHandler(ctx, req.Meta.UID)
	case SellerOrderDeliveredReports:
		return server.sellerOrderDeliveredReportsHandler(ctx, req.Meta.UID)
	case SellerOrderCancelReports:
		return server.sellerOrderCancelReportsHandler(ctx, req.Meta.UID)
	case SellerAllOrderReports:
		return server.sellerAllOrderReportsHandler(ctx, req.Meta.UID)
	case SellerApprovalPendingOrderReports:
		return server.sellerApprovalPendingOrderReportsHandler(ctx, req.Meta.UID)

	case BuyerOrderDetailList:
		return server.buyerOrderDetailListHandler(ctx, req.Meta.OID, req.Meta.UID, filterValue, req.Meta.Page, req.Meta.PerPage, sortName, sortDirection)
	case BuyerAllOrderReports:
		return server.buyerAllOrderReportsHandler(ctx, req.Meta.UID)
	case BuyerReturnOrderReports:
		return server.buyerReturnOrderReportsHandler(ctx, req.Meta.UID)
	case BuyerReturnOrderDetailList:
		return server.buyerReturnOrderDetailListHandler(ctx, req.Meta.UID, filterValue, req.Meta.Page, req.Meta.PerPage, sortName, sortDirection)

	case OperatorOrderList:
		return server.operatorOrderListHandler(ctx, req.Meta.OID, buyerMobile, filterValue, req.Meta.Page, req.Meta.PerPage, sortName, sortDirection)
	case OperatorOrderDetail:
		return server.operatorOrderDetailHandler(ctx, req.Meta.OID)
	case OperatorOrderInvoiceDetail:
		return server.operatorOrderInvoiceDetailHandler(ctx, req.Meta.OID)
	}

	return nil, status.Error(codes.Code(future.BadRequest), "Invalid Request")
}

func (server *Server) requestActionHandler(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	userType := UserType(req.Meta.UTP)
	var userAction actions.IAction

	app.Globals.Logger.FromContext(ctx).Debug("received request action", "fn", "requestActionHandler", "request", req)

	userActions, ok := server.actionStates[userType]
	if !ok {
		app.Globals.Logger.FromContext(ctx).Error("action userType not supported", "fn", "requestActionHandler", "utp", userType, "request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "User Action Invalid")
	}

	for _, action := range userActions {
		if action.ActionEnum().ActionName() == req.Meta.Action.ActionState {
			userAction = action
			break
		}
	}

	if userAction == nil {
		app.Globals.Logger.FromContext(ctx).Error("action invalid", "fn", "requestActionHandler", "action", req.Meta.Action.ActionState, "request", req)
		return nil, status.Error(codes.Code(future.BadRequest), "Action Invalid")
	}

	var reqActionData pb.ActionData
	if err := ptypes.UnmarshalAny(req.Data, &reqActionData); err != nil {
		app.Globals.Logger.FromContext(ctx).Error("Could not unmarshal reqActionData from request field", "fn", "requestActionHandler", "request", req, "error", err)
		return nil, status.Error(codes.Code(future.BadRequest), "Request Invalid")
	}

	subpackages := make([]events.ActionSubpackage, 0, len(reqActionData.Subpackages))
	for _, reqSubpackage := range reqActionData.Subpackages {
		subpackage := events.ActionSubpackage{
			SId: reqSubpackage.SID,
		}
		subpackage.Items = make([]events.ActionItem, 0, len(reqSubpackage.Items))
		for _, item := range reqSubpackage.Items {

			if item.Quantity <= 0 {
				app.Globals.Logger.FromContext(ctx).Error("action quantity invalid", "fn", "requestActionHandler", "action", req.Meta.Action.ActionState, "quantity", item.Quantity, "request", req)
				return nil, status.Error(codes.Code(future.BadRequest), "Action Quantity Invalid")
			}

			actionItem := events.ActionItem{
				InventoryId: item.InventoryId,
				Quantity:    item.Quantity,
			}
			if item.Reasons != nil {
				actionItem.Reasons = make([]models.Reason, 0, len(item.Reasons))
				for _, reason := range item.Reasons {
					rscnf, ok := server.reasonsMap[reason.Key]
					if !ok {
						return nil, status.Error(codes.Code(future.BadRequest), "reason not allowed")
					}
					rs := models.Reason{
						Key:         rscnf.Key,
						Translation: rscnf.Translation,
						Cancel:      rscnf.Cancel,
						Return:      rscnf.Return,
						Responsible: rscnf.Responsible,
					}
					if rscnf.HasDescription {
						rs.Description = reason.Description
					}
					actionItem.Reasons = append(actionItem.Reasons, rs)
				}
			}
			subpackage.Items = append(subpackage.Items, actionItem)
		}
		subpackages = append(subpackages, subpackage)
	}

	actionData := events.ActionData{
		SubPackages:    subpackages,
		Carrier:        reqActionData.Carrier,
		TrackingNumber: reqActionData.TrackingNumber,
	}

	event := events.New(events.Action, req.Meta.OID, req.Meta.PID, req.Meta.UID,
		req.Meta.Action.StateIndex, userAction,
		time.Unix(req.Time.GetSeconds(), int64(req.Time.GetNanos())), actionData)

	iFuture := future.Factory().SetCapacity(1).Build()
	iFrame := frame.Factory().SetFuture(iFuture).SetEvent(event).Build()
	server.flowManager.MessageHandler(ctx, iFrame)
	futureData := iFuture.Get()
	if futureData.Error() != nil {
		return nil, status.Error(codes.Code(futureData.Error().Code()), futureData.Error().Message())
	}

	eventResponse := futureData.Data().(events.ActionResponse)

	actionResponse := &pb.ActionResponse{
		OID:  eventResponse.OrderId,
		SIDs: eventResponse.SIds,
	}

	serializedResponse, err := proto.Marshal(actionResponse)
	if err != nil {
		app.Globals.Logger.FromContext(ctx).Error("could not marshal actionResponse", "fn", "requestActionHandler", "request", req, "response", actionResponse)
		return nil, status.Error(codes.Code(future.InternalError), "Unknown Error")
	}

	response := &pb.MessageResponse{
		Entity: "ActionResponse",
		Meta:   nil,
		Data: &any.Any{
			TypeUrl: "baman.io/" + proto.MessageName(actionResponse),
			Value:   serializedResponse,
		},
	}

	return response, nil
}

func (server *Server) PaymentGatewayHook(ctx context.Context, req *pg.PaygateHookRequest) (*pg.PaygateHookResponse, error) {

	app.Globals.Logger.FromContext(ctx).Debug("received payment response", "fn", "PaymentGatewayHook",
		"orderId", req.OrderID,
		"PaymentId", req.PaymentId,
		"InvoiceId", req.InvoiceId,
		"result", req.Result)
	futureData := server.flowManager.PaymentGatewayResult(ctx, req).Get()

	if futureData.Error() != nil {
		return nil, status.Error(codes.Code(futureData.Error().Code()), futureData.Error().Message())
	}

	return &pg.PaygateHookResponse{Ok: true}, nil
}

func (server Server) NewOrder(ctx context.Context, req *pb.RequestNewOrder) (*pb.ResponseNewOrder, error) {

	//ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	userAcl, err := app.Globals.UserService.AuthenticateContextToken(ctx)
	if err != nil {
		app.Globals.Logger.FromContext(ctx).Error("UserService.AuthenticateContextToken failed", "fn", "NewOrder",
			"error", err)
		return nil, status.Error(codes.Code(future.Forbidden), "User Not Authorized")
	}

	if uint64(userAcl.User().UserID) != req.Buyer.BuyerId {
		app.Globals.Logger.FromContext(ctx).Error("request userId with token userId mismatch", "fn", "NewOrder", "uid", req.Buyer.BuyerId, "token", userAcl.User().UserID)
		return nil, status.Error(codes.Code(future.Forbidden), "User Not Authorized")
	}

	if ctx.Value(string(utils.CtxUserID)) == nil {
		if userAcl != nil {
			ctx = context.WithValue(ctx, string(utils.CtxUserID), uint64(userAcl.User().UserID))
		} else {
			ctx = context.WithValue(ctx, string(utils.CtxUserID), uint64(0))
		}
	}

	if ctx.Value(string(utils.CtxUserACL)) == nil {
		if userAcl != nil {
			ctx = context.WithValue(ctx, string(utils.CtxUserACL), userAcl)
		}
	}

	iFuture := future.Factory().SetCapacity(1).Build()
	iFrame := frame.Factory().SetDefaultHeader(frame.HeaderNewOrder, req).SetFuture(iFuture).Build()
	server.flowManager.MessageHandler(ctx, iFrame)
	futureData := iFuture.Get()

	//if futureData.Error() != nil {
	//	futureErr := futureData.Error()
	//	return nil, status.Error(codes.Code(futureErr.Code()), futureErr.Message())
	//}

	var responseNewOrder pb.ResponseNewOrder

	if ipgResponse, ok := futureData.Data().(entities.PaymentIPGResponse); ok {
		responseNewOrder = pb.ResponseNewOrder{
			Action: pb.ResponseNewOrder_Redirect,
			Response: &pb.ResponseNewOrder_Ipg{
				Ipg: &pb.IPGResponse{
					CallbackUrl: ipgResponse.CallBackUrl,
				},
			},
		}

	} else if mpgResponse, ok := futureData.Data().(entities.PaymentMPGResponse); ok {
		responseNewOrder = pb.ResponseNewOrder{
			Action: pb.ResponseNewOrder_MPG,
			Response: &pb.ResponseNewOrder_Mpg{
				Mpg: &pb.MPGResponse{
					HostRequest:     mpgResponse.HostRequest,
					HostRequestSign: mpgResponse.HostRequestSign,
				},
			},
		}
	} else {
		app.Globals.Logger.FromContext(ctx).Error("NewOrder received data of futureData invalid", "fn", "NewOrder", "data", futureData.Data())
		return nil, status.Error(codes.Code(future.InternalError), "Unknown Error")
	}

	return &responseNewOrder, nil
}

func (server Server) ReasonsList(ctx context.Context, in *pb.ReasonsListRequest) (list *pb.ReasonsListResponse, err error) {
	ls := server.reasonsMap.toGRPC()
	list = &pb.ReasonsListResponse{
		Reasons: ls,
	}
	return
}

func (server Server) ReportOrderItems(req *pb.RequestReportOrderItems, srv pb.OrderService_ReportOrderItemsServer) error {

	//userAcl, err := app.Globals.UserService.AuthenticateContextToken(srv.Context())
	//if err != nil {
	//	app.Globals.Logger.Error("UserService.AuthenticateContextToken failed",
	//		"fn", "ReportOrderItems",
	//		"error", err)
	//	return status.Error(codes.Code(future.Forbidden), "User Not Authorized")
	//}

	//if userAcl.User().UserID <= 0 {
	//	app.Globals.Logger.Error("Token userId not authorized",
	//		"fn", "ReportOrderItems",
	//		"userId", userAcl.User().UserID)
	//	return status.Error(codes.Code(future.Forbidden), "User token not authorized")
	//}
	//
	//if !userAcl.UserPerm().Has("order.state.all.view") || !userAcl.UserPerm().Has("order.state.all.action") {
	//	return status.Error(codes.Code(future.Forbidden), "User Not Permitted")
	//}

	iFuture := server.flowManager.ReportOrderItems(srv.Context(), req, srv).Get()

	if iFuture.Error() != nil {
		return status.Error(codes.Code(iFuture.Error().Code()), iFuture.Error().Message())
	}

	return nil
}

func (server Server) VerifyUserSuccessOrder(ctx context.Context, req *pb.VerifyUserOrderRequest) (*pb.VerifyUserOrderResponse, error) {
	futureData := server.flowManager.VerifyUserSuccessOrder(ctx, req.UserId).Get()

	if futureData.Error() != nil {
		return nil, status.Error(codes.Code(futureData.Error().Code()), futureData.Error().Message())
	}

	app.Globals.Logger.FromContext(ctx).Debug("VerifyUserSuccessOrder received",
		"fn", "VerifyUserSuccessOrder",
		"uid", req.UserId,
		"IsSuccessOrder", futureData.Data().(bool))

	return &pb.VerifyUserOrderResponse{
		UserId:         req.UserId,
		IsSuccessOrder: futureData.Data().(bool),
	}, nil
}

func (server Server) Start() {
	port := strconv.Itoa(int(server.port))
	lis, err := net.Listen("tcp", server.address+":"+port)
	if err != nil {
		app.Globals.Logger.Error("Failed to listen to TCP on port", "fn", "Start", "port", port, "error", err)
	}
	app.Globals.Logger.Info("GRPC server started", "fn", "Start", "address", server.address, "port", port)

	customFunc := func(p interface{}) (err error) {
		app.Globals.Logger.Error("rpc panic recovered", "fn", "Start",
			"panic", p, "stacktrace", string(debug.Stack()))
		return grpc.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	//zapLogger, _ := zap.NewProduction()
	//stackDisableOpt := zap.AddStacktrace(stackTraceDisabler{})
	//noStackLogger := app.Globals.ZapLogger.WithOptions(stackDisableOpt)

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	uIntOpt := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(opts...),
		myUnaryLogger(app.Globals.Logger),
		//grpc_zap.UnaryServerInterceptor(zapLogger),
	))

	sIntOpt := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(opts...),
		//grpc_zap.StreamServerInterceptor(app.Globals.ZapLogger),
	))

	// enable grpc prometheus interceptors to log timing info for grpc APIs
	grpc_prometheus.EnableHandlingTimeHistogram()

	//Start GRPC server and register the server
	grpcServer := grpc.NewServer(uIntOpt, sIntOpt)
	pb.RegisterOrderServiceServer(grpcServer, &server)
	pg.RegisterBankResultHookServer(grpcServer, &server)
	if err := grpcServer.Serve(lis); err != nil {
		app.Globals.Logger.Error("GRPC server start field", "fn", "Start", "error", err.Error())
		panic("GRPC server start field")
	}
}

func (server Server) StartTest() {
	port := strconv.Itoa(int(server.port))
	lis, err := net.Listen("tcp", server.address+":"+port)
	if err != nil {
		applog.GLog.Logger.Error("Failed to listen to TCP",
			"port", port,
			"error", err.Error())
	}
	applog.GLog.Logger.Debug("app started", "address", server.address, "port", port)

	// Start GRPC server and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &server)
	pg.RegisterBankResultHookServer(grpcServer, &server)
	if err := grpcServer.Serve(lis); err != nil {
		panic("GRPC server start field")
	}
}

func myUnaryLogger(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		resp, err = handler(ctx, req)
		dur := time.Since(startTime)
		lg := log.FromContext(ctx)
		lg = lg.With(
			zap.Duration("took_sec", dur),
			zap.String("grpc.Method", path.Base(info.FullMethod)),
			zap.String("grpc.Service", path.Dir(info.FullMethod)[1:]),
			zap.String("grpc.Code", grpc.Code(err).String()),
		)
		lg.Debug("finished unary call")
		return
	}
}
