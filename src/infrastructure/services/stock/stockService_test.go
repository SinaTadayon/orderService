package stock_service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/configs"
	"gitlab.faza.io/order-project/order-service/domain/models/entities"
	stockProto "gitlab.faza.io/protos/stock-proto.git"
	"os"
	"testing"
	"time"
)

var config *configs.Cfg
var stock iStockServiceImpl

func createOrder() entities.Order {
	//currentTime := time.Now().UTC()

	paymentRequest := entities.PaymentRequest {
		Amount:	     	75400000,
		Currency:		"RR",
		Gateway: 		"AAP",
		CreatedAt:   	time.Now().UTC(),
	}

	paymentResponse	:= entities.PaymentResponse {
		Result:			true,
		Reason:			"",
		Description:	"",
		CallBackUrl:	"http://baman.io/payment-service",
		InvoiceId:		12345678946,
		PaymentId:		"r3r434ef45d",
		CreatedAt:   	time.Now().UTC(),
	}

	paymentResult := entities.PaymentResult {
		Result:			true,
		Reason:			"",
		PaymentId:      "r3r434ef45d",
		InvoiceId:		12345678946,
		Amount:    		75400000,
		ReqBody:   		"",
		ResBody:  		"",
		CardNumMask: 	"545498******4553",
		CreatedAt:   	time.Now().UTC(),
	}

	buyerInfo := entities.BuyerInfo {
		FirstName:			"Sina",
		LastName:   		"Tadayon",
		Mobile:     		"09123343534",
		Email:      		"sina.tadayon@baman.io",
		NationalId: 		"00598342521",
		Gender:				"male",
		IP:         		"127.0.0.1",
		FinanceInfo:   		entities.FinanceInfo {
			Iban:			"IR9450345802934803",
			CardNumber:		"4444555533332222",
			AccountNumber:	"293.6000.9439283.1",
			BankName:		"passargad",
		},
		ShippingAddress: 	entities.AddressInfo {
			Address:		"Tehran, Narmak, Golestan.st",
			Phone:   		"0217734873",
			Country: 		"Iran",
			City: 			"Tehran",
			Province: 		"Tehran",
			Neighbourhood:	"Chizar",
			Location:		entities.Location{
				Type:        "Point",
				Coordinates: []float64{-72.7738706, 41.6332836},
			},
			ZipCode: 		"1645630586",
		},
	}

	newOrder := entities.Order{
		OrderId: "",
		PaymentService: []entities.PaymentService{{
			PaymentRequest:  &paymentRequest,
			PaymentResponse: &paymentResponse,
			PaymentResult:   &paymentResult,
		}},
		SystemPayment: entities.SystemPayment{
			PayToBuyer: []entities.PayToBuyerInfo{{
				PaymentRequest:  &paymentRequest,
				PaymentResponse: &paymentResponse,
				PaymentResult:   &paymentResult,
			}},
			PayToSeller: []entities.PayToSellerInfo{{
				PaymentRequest:  &paymentRequest,
				PaymentResponse: &paymentResponse,
				PaymentResult:   &paymentResult,
			}},
			PayToMarket: []entities.PayToMarket{{
				PaymentRequest:  &paymentRequest,
				PaymentResponse: &paymentResponse,
				PaymentResult:   &paymentResult,
			}},
		},
		BuyerInfo: buyerInfo,
		Amount: entities.Amount {
			Total:         75400000,
			Subtotal:      73000000,
			Discount:      15600000,
			Currency:      "RR",
			ShipmentTotal: 5700000,
			PaymentMethod: "IPG",
			PaymentOption: "APP",
			Voucher:		&entities.Voucher{
				Amount:  230000,
				Code:    "Market",
				Details: nil,
			},
		},
		Items: []entities.Item{
			{
				ItemId:      "",
				InventoryId: "1111111111",
				Title:       "Mobile",
				Brand:       "Nokia",
				Guaranty:    "Sazegar",
				Category:    "Electronic",
				Image:       "",
				Returnable:  false,
				Quantity: 5,
				DeletedAt:   nil,
				Attributes:	 map[string]string {
					"Quantity":  "0",
					"Width":     "5cm",
					"Height":    "7cm",
					"Length":    "2m",
					"Weight":    "5kg",
					"Color":     "Blue",
					"Materials": "Stone",
				},
				SellerInfo: entities.SellerInfo{
					SellerId: 		  "129384234",
					Profile:            &entities.SellerProfile {
						Title:            "Sazgar",
						FirstName:        "Shahidi",
						LastName:         "nezhad",
						Mobile:           "019124343",
						Email:            "shahidi@samsong.com",
						NationalId:       "9793287434",
						CompanyName:      "Samservice",
						RegistrationName: "Sazgar",
						EconomicCode:     "342346434343",
						Finance: entities.FinanceInfo{
							Iban:          "IR92347299384782734",
							CardNumber:    "8888777766665555",
							AccountNumber: "983.234.2948723894.2",
							BankName:      "saderat",
						},
						ShippingAddress: entities.AddressInfo{
							Address:       "Tehran, Jordan",
							Phone:         "01249874345",
							Country:       "Iran",
							City:          "Tehran",
							Province:      "Tehran",
							Neighbourhood: "Narmak",
							Location: entities.Location{
								Type:        "Point",
								Coordinates: []float64{-104.7738706, 54.6332836},
							},
							ZipCode: "947534586",
						},
					},
				},
				Price: entities.Price{
					Unit:             1270000,
					Original:         7340000,
					Special:          1000000,
					SellerCommission: 5334444,
					Currency:         "RR",
				},
				ShipmentSpec: entities.ShipmentSpec {
					CarrierName:    "Post",
					CarrierProduct: "Post Express",
					CarrierType:    "Standard",
					ShippingCost:   1249348,
					VoucherAmount:  3242344,
					Currency:       "RR",
					ReactionTime:   2,
					ShippingTime:   8,
					ReturnTime:     24,
					Details:        "no return",
				},
				ShipmentDetails: entities.ShipmentDetails{
					SellerShipmentDetail: 	entities.ShipmentDetail{
						CarrierName: 		"Post",
						TrackingNumber:   	"545349534958349",
						Image:            	"",
						Description:      	"",
						CreatedAt:        	time.Now().UTC(),
					},
					BuyerReturnShipmentDetail: entities.ShipmentDetail{
						CarrierName: 			"Post",
						TrackingNumber:   		"545349534958349",
						Image:                  "",
						Description:            "",
						CreatedAt:              time.Now().UTC(),
					},
				},
				Progress: entities.Progress{
					CurrentStepName:  "0.NewOrder",
					CurrentStepIndex: 0,
					//CurrentState: entities.State {
					//	Name:  "0.New_Order_Process_State",
					//	Index: 0,
					//	Type: "LauncherAction",
					//	Actions: []entities.Action {{
					//		Name: "Success",
					//		Type: "NewOrder",
					//		Base: "Active",
					//		Data: nil,
					//		Time: &currentTime,
					//	}},
					//	AcceptedAction:entities.Action {
					//		Name: "Success",
					//		Type: "NewOrder",
					//		Base: "Active",
					//		Data: nil,
					//		Time: &currentTime,
					//	},
					//	Result: false,
					//	Reason:       "",
					//	CreatedAt:    time.Now().UTC(),
					//},
					CreatedAt: time.Now().UTC(),
					StepsHistory: []entities.StepHistory{{
						Name:      "0.NewOrder",
						Index:     0,
						CreatedAt: time.Now().UTC(),
						//StatesHistory: []entities.StateHistory{{
						//	Name:  "0.New_Order_Process_State",
						//	Index: 0,
						//	Type: "ListenerAction",
						//	Action: entities.Action{
						//		Name:           "Success",
						//		Type:           "NewOrder",
						//		Base:           "Active",
						//		Data:           nil,
						//		Time: 			&currentTime,
						//	},
						//	Result: 	  false,
						//	Reason:       "",
						//	CreatedAt:    time.Now().UTC(),
						//}},
					}},
				},
			},
			{
				ItemId:      "",
				InventoryId: "2222222222",
				Title:       "Laptop",
				Brand:       "Lenovo",
				Guaranty:    "Iranargham",
				Category:    "Electronic",
				Image:       "",
				Returnable:  true,
				Quantity: 5,
				DeletedAt:   nil,
				Attributes:	 map[string]string {
					"Quantity":  "0",
					"Width":     "5cm",
					"Height":    "7cm",
					"Length":    "2m",
					"Weight":    "5kg",
					"Color":     "Blue",
					"Materials": "Stone",
				},
				SellerInfo: entities.SellerInfo{
					SellerId:			"2384723083",
					Profile:			&entities.SellerProfile{
						Title:            "Avazhang",
						FirstName:        "Mostafavi",
						LastName:         "Rezaii",
						Mobile:           "0394739844",
						Email:            "mostafavi@samsong.com",
						NationalId:       "39458979455",
						CompanyName:      "Avazhang",
						RegistrationName: "Avazeh",
						EconomicCode:     "3045988273784",
						Finance: entities.FinanceInfo{
							Iban:          "IR209345882374",
							CardNumber:    "92384787263443443",
							AccountNumber: "983.234.2293452434.2",
							BankName:      "saderat",
						},
						ShippingAddress: entities.AddressInfo{
							Address:       "Tehran, Jordan",
							Phone:         "01249874345",
							Country:       "Iran",
							City:          "Tehran",
							Province:      "Tehran",
							Neighbourhood: "Navab",
							Location: entities.Location{
								Type:        "Point",
								Coordinates: []float64{-104.7738706, 54.6332836},
							},
							ZipCode: "947534586",
						},
					},
				},
				Price: entities.Price{
					Unit:             1270000,
					Original:         7340000,
					Special:          1000000,
					SellerCommission: 5334444,
					Currency:         "RR",
				},
				ShipmentSpec: entities.ShipmentSpec{
					CarrierName:    "Post",
					CarrierProduct: "Post Express",
					CarrierType:    "Standard",
					ShippingCost:   1249348,
					VoucherAmount:  3242344,
					Currency:       "RR",
					ReactionTime:   2,
					ShippingTime:   8,
					ReturnTime:     24,
					Details:        "no return",
				},
				ShipmentDetails: entities.ShipmentDetails{
					SellerShipmentDetail: 	entities.ShipmentDetail{
						CarrierName: 		"Post",
						TrackingNumber:   	"545349534958349",
						Image:            	"",
						Description:      	"",
						CreatedAt:        	time.Now().UTC(),
					},
					BuyerReturnShipmentDetail: entities.ShipmentDetail{
						CarrierName: 			"Post",
						TrackingNumber:   		"545349534958349",
						Image:                  "",
						Description:            "",
						CreatedAt:              time.Now().UTC(),
					},
				},
				Progress: entities.Progress{
					CurrentStepName:  "0.NewOrder",
					CurrentStepIndex: 0,
					//CurrentState: entities.State{
					//	Name:  "0.New_Order_Process_State",
					//	Index: 0,
					//	Actions: []entities.Action{{
					//		Name:           "Success",
					//		Type:           "NewOrder",
					//		Base:           "Active",
					//		Data:           nil,
					//		Time: 			&currentTime,
					//	}},
					//	AcceptedAction: entities.Action{
					//		Name:           "Success",
					//		Type:           "NewOrder",
					//		Base:           "Active",
					//		Data:           nil,
					//		Time: 			&currentTime,
					//	},
					//	Result: false,
					//	Reason:       "",
					//	CreatedAt:    time.Now().UTC(),
					//},
					CreatedAt: time.Now().UTC(),
					StepsHistory: []entities.StepHistory{{
						Name:      "0.NewOrder",
						Index:     0,
						CreatedAt: time.Now().UTC(),
						//StatesHistory: []entities.StateHistory{{
						//	Name:  "0.New_Order_Process_State",
						//	Index: 0,
						//	Type: "ListenerAction",
						//	Action: entities.Action{
						//		Name:           "Success",
						//		Type:           "NewOrder",
						//		Base:           "Active",
						//		Data:           nil,
						//		Time: 			&currentTime,
						//	},
						//
						//	Result: false,
						//	Reason:       "",
						//	CreatedAt:    time.Now().UTC(),
						//}},
					}},
				},
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
	}

	return newOrder
}

func init() {
	var err error
	var path string
	if os.Getenv("APP_ENV") == "dev" {
		path = "../../../testdata/.env"
	} else {
		path = ""
	}

	config, err = configs.LoadConfig(path)
	if err != nil {
		logger.Err(err.Error())
		panic("configs.LoadConfig failed")
	}

	stock = iStockServiceImpl{nil, nil,
		config.StockService.Address, config.StockService.Port}

	if err := stock.connectToStockService(); err != nil {
		logger.Err(err.Error())
		panic("stockService.connectToStockService() failed")
	}
}

func TestStockService_ReservedSuccess(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	order := createOrder()

	request := stockProto.StockRequest{
		Quantity:    5,
		InventoryId: order.Items[0].InventoryId,
	}
	_, err := stock.stockService.StockAllocate(ctx, &request)
	assert.Nil(t, err)

	promise := stock.BatchStockActions(ctx, order, "StockReserved")
	futuredata := promise.Data()

	assert.Nil(t, futuredata.Ex)
	assert.Nil(t, futuredata.Data)

}
