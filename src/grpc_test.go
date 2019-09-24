package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"

	pb "gitlab.faza.io/protos/payment"
)

func createOrderObject() *pb.OrderPaymentRequest {
	req := createPaymentRequestSampleFull()

	order := &pb.OrderPaymentRequest{
		Amount: &pb.Amount{},
		Buyer: &pb.Buyer{
			Finance: &pb.BuyerFinance{},
			Address: &pb.BuyerAddress{},
		},
	}

	order.Amount.Total = float32(req.Amount.Total)
	order.Amount.Payable = float32(req.Amount.Payable)
	order.Amount.Discount = float32(req.Amount.Discount)

	order.Buyer.LastName = req.Buyer.LastName
	order.Buyer.FirstName = req.Buyer.FirstName
	order.Buyer.Email = req.Buyer.Email
	order.Buyer.Mobile = req.Buyer.Mobile
	order.Buyer.NationalId = req.Buyer.NationalId
	order.Buyer.Ip = req.Buyer.IP

	order.Buyer.Finance.Iban = req.Buyer.Finance.Iban

	order.Buyer.Address.Address = req.Buyer.Address.Address
	order.Buyer.Address.State = req.Buyer.Address.State
	order.Buyer.Address.Phone = req.Buyer.Address.Phone
	order.Buyer.Address.ZipCode = req.Buyer.Address.ZipCode
	order.Buyer.Address.City = req.Buyer.Address.City
	order.Buyer.Address.Country = req.Buyer.Address.Country
	order.Buyer.Address.Lat = req.Buyer.Address.Lat
	order.Buyer.Address.Lan = req.Buyer.Address.Lan

	i := pb.Item{
		Price:    &pb.ItemPrice{},
		Shipment: &pb.ItemShipment{},
		Seller: &pb.ItemSeller{
			Address: &pb.ItemSellerAddress{},
			Finance: &pb.ItemSellerFinance{},
		},
	}
	i.Sku = req.Items[0].Sku
	i.Brand = req.Items[0].Brand
	i.Categories = req.Items[0].Categories
	i.Title = req.Items[0].Title
	i.Warranty = req.Items[0].Warranty
	i.Quantity = req.Items[0].Quantity

	i.Price.Discount = float32(req.Items[0].Price.Discount)
	i.Price.Payable = float32(req.Items[0].Price.Payable)
	i.Price.Total = float32(req.Items[0].Price.Total)
	i.Price.SellerCommission = float32(req.Items[0].Price.SellerCommission)
	i.Price.Unit = float32(req.Items[0].Price.Unit)

	i.Shipment.ShipmentDetail = req.Items[0].Shipment.ShipmentDetail
	i.Shipment.ShippingTime = req.Items[0].Shipment.ShippingTime
	i.Shipment.ReturnTime = req.Items[0].Shipment.ReturnTime
	i.Shipment.ReactionTime = req.Items[0].Shipment.ReactionTime
	i.Shipment.ProviderName = req.Items[0].Shipment.ProviderName

	i.Seller.Title = req.Items[0].Seller.Title
	i.Seller.NationalId = req.Items[0].Seller.NationalId
	i.Seller.Mobile = req.Items[0].Seller.Mobile
	i.Seller.FirstName = req.Items[0].Seller.FirstName
	i.Seller.LastName = req.Items[0].Seller.LastName
	i.Seller.Email = req.Items[0].Seller.Email
	i.Seller.RegistrationName = req.Items[0].Seller.RegistrationName
	i.Seller.CompanyName = req.Items[0].Seller.CompanyName

	i.Seller.Address.Address = req.Items[0].Seller.Address.Address
	i.Seller.Address.Lan = req.Items[0].Seller.Address.Lan
	i.Seller.Address.Lat = req.Items[0].Seller.Address.Lat
	i.Seller.Address.Country = req.Items[0].Seller.Address.Country
	i.Seller.Address.City = req.Items[0].Seller.Address.City
	i.Seller.Address.ZipCode = req.Items[0].Seller.Address.ZipCode
	i.Seller.Address.Phone = req.Items[0].Seller.Address.Phone
	i.Seller.Address.State = req.Items[0].Seller.Address.State
	i.Seller.Address.Title = req.Items[0].Seller.Address.Title

	i.Seller.Finance.Iban = req.Items[0].Seller.Finance.Iban

	order.Items = append(order.Items, &i)
	return order
}

func TestAddGrpcStateRule(t *testing.T) {
	GrpcStatesRules.SellerApprovalPending = addStateRule(PaymentSuccess)

	_, ok := GrpcStatesRules.SellerApprovalPending[PaymentSuccess]
	assert.True(t, ok)
}

// Grpc test
func TestNewOrder(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnNewOrder, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	PaymentService := pb.NewOrderServiceClient(grpcConnNewOrder)

	order := createOrderObject()

	resOrder, err := PaymentService.NewOrder(ctx, order)
	assert.Nil(t, err)
	assert.NotNil(t, resOrder)
}
func TestSellerApprovalPendingApproved(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ApprovalRequest{
		OrderNumber: ppr.OrderNumber,
		Approval:    true,
		Reason:      "",
	}

	resApproval, err := OrderService.SellerApprovalPending(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentPending, savedOrder.Status.Current)
}
func TestSellerApprovalPendingRejected(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	ppr.Status.Current = SellerApprovalPending

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ApprovalRequest{
		OrderNumber: ppr.OrderNumber,
		Approval:    false,
		Reason:      "out of stock",
	}

	resApproval, err := OrderService.SellerApprovalPending(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+2, len(savedOrder.Status.History))
	assert.Equal(t, PayToBuyer, savedOrder.Status.Current)
}
func TestShipmentDetail(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	// insert to mongo
	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)
	time.Sleep(time.Second)
	// call grpc
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	shipmentDetail := &pb.ShipmentDetailRequest{
		OrderNumber:            ppr.OrderNumber,
		ShipmentTrackingNumber: "Track1234",
		ShipmentProvider:       "SnappBox",
	}
	resDetail, err := OrderService.ShipmentDetail(ctx, shipmentDetail)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resDetail.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, Shipped, savedOrder.Status.Current)
}
func TestBuyerCancel(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	request := &pb.BuyerCancelRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "its took soo much time",
	}

	resApproval, err := OrderService.BuyerCancel(ctx, request)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentCanceled, savedOrder.Status.Current)
}
func TestDelivered(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentDeliveredRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ShipmentDelivered(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentDelivered, savedOrder.Status.Current)
}
func TestShipmentDeliveryDelay(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDeliveryPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentDeliveryDelayedRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ShipmentDeliveryDelayed(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentDeliveryDelayed, savedOrder.Status.Current)
}
func TestReturnShipmentDeliveryDelay(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "this how i return it",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "return shipment delivery dayes reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDeliveryDelayedRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ReturnShipmentDeliveryDelayed(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentDeliveryDelayed, savedOrder.Status.Current)
}
func TestShipmentCanceled(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDeliveryPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDeliveryDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentCanceledRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "operator decide to cancel",
		Operator:    "operator",
	}

	resApproval, err := OrderService.ShipmentCanceled(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+2, len(savedOrder.Status.History))
	assert.Equal(t, PayToBuyer, savedOrder.Status.Current)
}
func TestReturnShipmentCanceledFromReturnShipmentDeliveryProblem(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i received return package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryProblem,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentCanceledRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "operator decide to cancel",
	}

	resApproval, err := OrderService.ReturnShipmentCanceled(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+2, len(savedOrder.Status.History))
	assert.Equal(t, PayToSeller, savedOrder.Status.Current)
}
func TestReturnShipmentCanceledFromReturnShipmentDeliveryDelayed(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentCanceledRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "operator decide to cancel",
	}

	resApproval, err := OrderService.ReturnShipmentCanceled(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+2, len(savedOrder.Status.History))
	assert.Equal(t, PayToSeller, savedOrder.Status.Current)
}
func TestShipmentDeliveryProblem(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentDeliveryProblemRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "operator decide to cancel",
	}

	resApproval, err := OrderService.ShipmentDeliveryProblem(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentDeliveryProblem, savedOrder.Status.Current)
}
func TestReturnShipmentDeliveryProblem(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i sent it",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got it",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDeliveryProblemRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "operator decide to cancel",
	}

	resApproval, err := OrderService.ReturnShipmentDeliveryProblem(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentDeliveryProblem, savedOrder.Status.Current)
}
func TestShipmentSuccessFromShipmentDelivered(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentSuccessRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ShipmentSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentSuccess, savedOrder.Status.Current)
}
func TestShipmentSuccessFromShipmentDeliveryProblem(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDeliveryProblem,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentSuccessRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ShipmentSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentSuccess, savedOrder.Status.Current)
}
func TestShipmentSuccessFromReturnShipmentDetailDelayed(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i need to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "x days past",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ShipmentSuccessRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ShipmentSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ShipmentSuccess, savedOrder.Status.Current)
}
func TestReturnShipmentPendingFromShipmentDeliveryProblem(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDeliveryProblem,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentPendingRequest{
		OrderNumber: ppr.OrderNumber,
		Operator:    "buyer",
		Reason:      "package broken",
	}

	resApproval, err := OrderService.ReturnShipmentPending(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentPending, savedOrder.Status.Current)
}
func TestReturnShipmentPendingFromShipmentDelivered(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentPendingRequest{
		OrderNumber: ppr.OrderNumber,
		Operator:    "buyer",
		Reason:      "package broken",
	}

	resApproval, err := OrderService.ReturnShipmentPending(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentPending, savedOrder.Status.Current)
}
func TestReturnShipmentDetailFromReturnShipmentPending(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDetailRequest{
		OrderNumber:            ppr.OrderNumber,
		ShipmentTrackingNumber: "10232153",
		Description:            "i send it via some 3pl",
		ShipmentProvider:       "snapp box",
	}

	resApproval, err := OrderService.ReturnShipmentDetail(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipped, savedOrder.Status.Current)
}
func TestReturnShipmentDetailFromReturnShipmentDetailDelayed(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDetailRequest{
		OrderNumber:            ppr.OrderNumber,
		ShipmentTrackingNumber: "10232153",
		Description:            "i send it via some 3pl",
		ShipmentProvider:       "snapp box",
	}

	resApproval, err := OrderService.ReturnShipmentDetail(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipped, savedOrder.Status.Current)
}
func TestReturnShipmentDeliveredFromReturnShipmentDeliveryDelayed(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "its how i send it",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDeliveredRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "",
		Operator:    "operator",
	}

	resApproval, err := OrderService.ReturnShipmentDelivered(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentDelivered, savedOrder.Status.Current)
}
func TestReturnShipmentDeliveredFromReturnShipmentDeliveryPending(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "its how i send it",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDeliveredRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "no action or x days",
		Operator:    "system",
	}

	resApproval, err := OrderService.ReturnShipmentDelivered(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentDelivered, savedOrder.Status.Current)
}
func TestReturnShipmentDeliveredFromReturnShipped(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentDeliveredRequest{
		OrderNumber: ppr.OrderNumber,
		Reason:      "no action or x days",
		Operator:    "system",
	}

	resApproval, err := OrderService.ReturnShipmentDelivered(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, ReturnShipmentDelivered, savedOrder.Status.Current)
}
func TestReturnShipmentSuccessFromReturnShipmentDeliveryProblem(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryProblem,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentSuccessRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ReturnShipmentSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+2, len(savedOrder.Status.History))
	assert.Equal(t, PayToBuyer, savedOrder.Status.Current)
}
func TestReturnShipmentSuccessFromReturnShipmentDelivered(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.ReturnShipmentSuccessRequest{
		OrderNumber: ppr.OrderNumber,
	}

	resApproval, err := OrderService.ReturnShipmentSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+2, len(savedOrder.Status.History))
	assert.Equal(t, PayToBuyer, savedOrder.Status.Current)
}
func TestPayToBuyerSuccess(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "package is okey",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToBuyer,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToBuyerFailed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "cant pay",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.PayToBuyerSuccessRequest{
		OrderNumber: ppr.OrderNumber,
		Description: "paid by XXXXX",
	}

	resApproval, err := OrderService.PayToBuyerSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, PayToBuyerSuccess, savedOrder.Status.Current)
}
func TestPayToSellerSuccess(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryProblem,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentCanceled,
		CreatedAt: time.Now().UTC(),
		Agent:     "operator",
		Reason:    "not returnable",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToSeller,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToSellerFailed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "cant pay",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.PayToSellerSuccessRequest{
		OrderNumber: ppr.OrderNumber,
		Description: "paid by XXXXX",
	}

	resApproval, err := OrderService.PayToSellerSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, PayToSellerSuccess, savedOrder.Status.Current)
}
func TestPayToMarketSuccess(t *testing.T) {
	// Create ppr
	ppr := createPaymentRequestSampleFull()
	// Delete test order
	_, err := App.mongo.DeleteOne(MongoDB, Orders, bson.D{{"ordernumber", ppr.OrderNumber}})
	assert.Nil(t, err)
	statusHistory := StatusHistory{
		Status:    PaymentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PaymentSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    SellerApprovalPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto approval",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "hale",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    Shipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "seller add detail",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "shipping days reached",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentPending,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i want to return",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDetailDelayed,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "no action for x days",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipped,
		CreatedAt: time.Now().UTC(),
		Agent:     "buyer",
		Reason:    "i send the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDelivered,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i got the package",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentDeliveryProblem,
		CreatedAt: time.Now().UTC(),
		Agent:     "seller",
		Reason:    "i need support",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    ReturnShipmentCanceled,
		CreatedAt: time.Now().UTC(),
		Agent:     "operator",
		Reason:    "not returnable",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToSeller,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "auto",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToSellerSuccess,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "paid",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToMarket,
		CreatedAt: time.Now().UTC(),
		Agent:     "system",
		Reason:    "paid",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)
	statusHistory = StatusHistory{
		Status:    PayToMarketFailed,
		CreatedAt: time.Now().UTC(),
		Agent:     "operator",
		Reason:    "paid by XXXXX",
	}
	ppr.Status.History = append(ppr.Status.History, statusHistory)

	ppr.Status.Current = ppr.Status.History[(len(ppr.Status.History) - 1)].Status

	_, err = App.mongo.InsertOne(MongoDB, Orders, ppr)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	grpcConnOrderApproved, err := grpc.DialContext(ctx, ":"+fmt.Sprint(App.config.App.Port), grpc.WithInsecure())
	assert.Nil(t, err)
	OrderService := pb.NewOrderServiceClient(grpcConnOrderApproved)

	approveRequest := &pb.PayToMarketSuccessRequest{
		OrderNumber: ppr.OrderNumber,
		Description: "paid by XXXXX",
	}

	resApproval, err := OrderService.PayToMarketSuccess(ctx, approveRequest)
	assert.Nil(t, err)
	assert.Equal(t, ppr.OrderNumber, resApproval.OrderNumber)

	savedOrder, err := GetOrder(ppr.OrderNumber)
	assert.Nil(t, err)
	assert.Equal(t, len(ppr.Status.History)+1, len(savedOrder.Status.History))
	assert.Equal(t, PayToMarketSuccess, savedOrder.Status.Current)
}
