package voucher_service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.faza.io/go-framework/logger"
	"gitlab.faza.io/order-project/order-service/configs"
	voucherProto "gitlab.faza.io/protos/cart"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var config *configs.Cfg
var voucherSrv iVoucherServiceImpl

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

	voucherSrv = iVoucherServiceImpl{
		voucherClient:  nil,
		grpcConnection: nil,
		serverAddress:  config.VoucherService.Address,
		serverPort:     config.VoucherService.Port,
	}
}

func TestVoucherSettlement(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testName1 := "test-" + strconv.Itoa(rand.Int())

	cT := &voucherProto.CouponTemplate{
		Title:           testName1,
		Prefix:          testName1,
		UseLimit:        1,
		Count:           1,
		Length:          5,
		StartDate:       time.Date(2019, 07, 24, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		EndDate:         time.Date(2019, 07, 25, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		Categories:      nil,
		Products:        nil,
		Users:           []string{"1000002"},
		Sellers:         nil,
		IsFirstPurchase: true,
		CouponDiscount: &voucherProto.CouponDiscount{
			Type:             "fixed",
			Amount:           100000,
			MaxDiscountValue: 0,
			MinBasketValue:   1500000,
		},
	}

	err := voucherSrv.Connect()
	assert.Nil(t, err)

	defer voucherSrv.Disconnect()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := voucherSrv.voucherClient.CreateCouponTemplate(ctx, cT)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 200, int(result.Code))

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	voucherRequest := &voucherProto.GetVoucherByTemplateNameRequest{
		Page:        1,
		Perpage:     1,
		VoucherName: testName1,
	}
	allVouchers, err := voucherSrv.voucherClient.GetVoucherByTemplateName(ctx, voucherRequest)
	assert.Nil(t, err)
	assert.NotNil(t, allVouchers.Vouchers[0])
	assert.NotEmpty(t, allVouchers.Vouchers[0].Code)

	ctx, _ = context.WithCancel(context.Background())
	iFuture := voucherSrv.VoucherSettlement(ctx, allVouchers.Vouchers[0].Code, 123456789776, 1000002)
	futureData := iFuture.Get()
	assert.Nil(t, futureData.Data())
	assert.Nil(t, futureData.Error())
}
