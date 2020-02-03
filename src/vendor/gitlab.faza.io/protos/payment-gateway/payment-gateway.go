// Code generated by protoc-gen-go. DO NOT EDIT.
// source: payment-gateway.proto

package payment_gateway

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PaymentRequest_Status int32

const (
	PaymentRequest_PENDING PaymentRequest_Status = 0
	PaymentRequest_SUCCESS PaymentRequest_Status = 1
	PaymentRequest_FAIL    PaymentRequest_Status = 2
)

var PaymentRequest_Status_name = map[int32]string{
	0: "PENDING",
	1: "SUCCESS",
	2: "FAIL",
}

var PaymentRequest_Status_value = map[string]int32{
	"PENDING": 0,
	"SUCCESS": 1,
	"FAIL":    2,
}

func (x PaymentRequest_Status) String() string {
	return proto.EnumName(PaymentRequest_Status_name, int32(x))
}

func (PaymentRequest_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{5, 0}
}

type MPGValidateRequest struct {
	Message              string   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	HostResponse         string   `protobuf:"bytes,2,opt,name=HostResponse,proto3" json:"HostResponse,omitempty"`
	HostResponseSign     string   `protobuf:"bytes,3,opt,name=HostResponseSign,proto3" json:"HostResponseSign,omitempty"`
	StatusCode           int64    `protobuf:"varint,4,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	UniqueTranID         int64    `protobuf:"varint,5,opt,name=UniqueTranID,proto3" json:"UniqueTranID,omitempty"`
	PaymentID            string   `protobuf:"bytes,7,opt,name=PaymentID,proto3" json:"PaymentID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MPGValidateRequest) Reset()         { *m = MPGValidateRequest{} }
func (m *MPGValidateRequest) String() string { return proto.CompactTextString(m) }
func (*MPGValidateRequest) ProtoMessage()    {}
func (*MPGValidateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{0}
}

func (m *MPGValidateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MPGValidateRequest.Unmarshal(m, b)
}
func (m *MPGValidateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MPGValidateRequest.Marshal(b, m, deterministic)
}
func (m *MPGValidateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MPGValidateRequest.Merge(m, src)
}
func (m *MPGValidateRequest) XXX_Size() int {
	return xxx_messageInfo_MPGValidateRequest.Size(m)
}
func (m *MPGValidateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MPGValidateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MPGValidateRequest proto.InternalMessageInfo

func (m *MPGValidateRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *MPGValidateRequest) GetHostResponse() string {
	if m != nil {
		return m.HostResponse
	}
	return ""
}

func (m *MPGValidateRequest) GetHostResponseSign() string {
	if m != nil {
		return m.HostResponseSign
	}
	return ""
}

func (m *MPGValidateRequest) GetStatusCode() int64 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *MPGValidateRequest) GetUniqueTranID() int64 {
	if m != nil {
		return m.UniqueTranID
	}
	return 0
}

func (m *MPGValidateRequest) GetPaymentID() string {
	if m != nil {
		return m.PaymentID
	}
	return ""
}

type MPGValidateResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MPGValidateResponse) Reset()         { *m = MPGValidateResponse{} }
func (m *MPGValidateResponse) String() string { return proto.CompactTextString(m) }
func (*MPGValidateResponse) ProtoMessage()    {}
func (*MPGValidateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{1}
}

func (m *MPGValidateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MPGValidateResponse.Unmarshal(m, b)
}
func (m *MPGValidateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MPGValidateResponse.Marshal(b, m, deterministic)
}
func (m *MPGValidateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MPGValidateResponse.Merge(m, src)
}
func (m *MPGValidateResponse) XXX_Size() int {
	return xxx_messageInfo_MPGValidateResponse.Size(m)
}
func (m *MPGValidateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MPGValidateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MPGValidateResponse proto.InternalMessageInfo

func (m *MPGValidateResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type MPGStartRequest struct {
	Amount               int64    `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	Currency             string   `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	OrderID              string   `protobuf:"bytes,3,opt,name=orderID,proto3" json:"orderID,omitempty"`
	Mobile               string   `protobuf:"bytes,4,opt,name=mobile,proto3" json:"mobile,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MPGStartRequest) Reset()         { *m = MPGStartRequest{} }
func (m *MPGStartRequest) String() string { return proto.CompactTextString(m) }
func (*MPGStartRequest) ProtoMessage()    {}
func (*MPGStartRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{2}
}

func (m *MPGStartRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MPGStartRequest.Unmarshal(m, b)
}
func (m *MPGStartRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MPGStartRequest.Marshal(b, m, deterministic)
}
func (m *MPGStartRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MPGStartRequest.Merge(m, src)
}
func (m *MPGStartRequest) XXX_Size() int {
	return xxx_messageInfo_MPGStartRequest.Size(m)
}
func (m *MPGStartRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MPGStartRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MPGStartRequest proto.InternalMessageInfo

func (m *MPGStartRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *MPGStartRequest) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *MPGStartRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *MPGStartRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

type MPGStartResponse struct {
	HostRequest          string   `protobuf:"bytes,1,opt,name=hostRequest,proto3" json:"hostRequest,omitempty"`
	HostRequestSign      string   `protobuf:"bytes,2,opt,name=hostRequestSign,proto3" json:"hostRequestSign,omitempty"`
	PaymentId            string   `protobuf:"bytes,3,opt,name=paymentId,proto3" json:"paymentId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MPGStartResponse) Reset()         { *m = MPGStartResponse{} }
func (m *MPGStartResponse) String() string { return proto.CompactTextString(m) }
func (*MPGStartResponse) ProtoMessage()    {}
func (*MPGStartResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{3}
}

func (m *MPGStartResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MPGStartResponse.Unmarshal(m, b)
}
func (m *MPGStartResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MPGStartResponse.Marshal(b, m, deterministic)
}
func (m *MPGStartResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MPGStartResponse.Merge(m, src)
}
func (m *MPGStartResponse) XXX_Size() int {
	return xxx_messageInfo_MPGStartResponse.Size(m)
}
func (m *MPGStartResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MPGStartResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MPGStartResponse proto.InternalMessageInfo

func (m *MPGStartResponse) GetHostRequest() string {
	if m != nil {
		return m.HostRequest
	}
	return ""
}

func (m *MPGStartResponse) GetHostRequestSign() string {
	if m != nil {
		return m.HostRequestSign
	}
	return ""
}

func (m *MPGStartResponse) GetPaymentId() string {
	if m != nil {
		return m.PaymentId
	}
	return ""
}

type GetPaymentResultByOrderIdRequest struct {
	OrderID              string   `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPaymentResultByOrderIdRequest) Reset()         { *m = GetPaymentResultByOrderIdRequest{} }
func (m *GetPaymentResultByOrderIdRequest) String() string { return proto.CompactTextString(m) }
func (*GetPaymentResultByOrderIdRequest) ProtoMessage()    {}
func (*GetPaymentResultByOrderIdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{4}
}

func (m *GetPaymentResultByOrderIdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPaymentResultByOrderIdRequest.Unmarshal(m, b)
}
func (m *GetPaymentResultByOrderIdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPaymentResultByOrderIdRequest.Marshal(b, m, deterministic)
}
func (m *GetPaymentResultByOrderIdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPaymentResultByOrderIdRequest.Merge(m, src)
}
func (m *GetPaymentResultByOrderIdRequest) XXX_Size() int {
	return xxx_messageInfo_GetPaymentResultByOrderIdRequest.Size(m)
}
func (m *GetPaymentResultByOrderIdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPaymentResultByOrderIdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPaymentResultByOrderIdRequest proto.InternalMessageInfo

func (m *GetPaymentResultByOrderIdRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

type PaymentRequest struct {
	OrderID              string                `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	PaymentId            string                `protobuf:"bytes,2,opt,name=paymentId,proto3" json:"paymentId,omitempty"`
	InvoiceId            int64                 `protobuf:"varint,3,opt,name=invoiceId,proto3" json:"invoiceId,omitempty"`
	Amount               int64                 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
	CardMask             string                `protobuf:"bytes,5,opt,name=cardMask,proto3" json:"cardMask,omitempty"`
	Status               PaymentRequest_Status `protobuf:"varint,6,opt,name=status,proto3,enum=payment_gateway.PaymentRequest_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PaymentRequest) Reset()         { *m = PaymentRequest{} }
func (m *PaymentRequest) String() string { return proto.CompactTextString(m) }
func (*PaymentRequest) ProtoMessage()    {}
func (*PaymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{5}
}

func (m *PaymentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PaymentRequest.Unmarshal(m, b)
}
func (m *PaymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PaymentRequest.Marshal(b, m, deterministic)
}
func (m *PaymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PaymentRequest.Merge(m, src)
}
func (m *PaymentRequest) XXX_Size() int {
	return xxx_messageInfo_PaymentRequest.Size(m)
}
func (m *PaymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PaymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PaymentRequest proto.InternalMessageInfo

func (m *PaymentRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *PaymentRequest) GetPaymentId() string {
	if m != nil {
		return m.PaymentId
	}
	return ""
}

func (m *PaymentRequest) GetInvoiceId() int64 {
	if m != nil {
		return m.InvoiceId
	}
	return 0
}

func (m *PaymentRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *PaymentRequest) GetCardMask() string {
	if m != nil {
		return m.CardMask
	}
	return ""
}

func (m *PaymentRequest) GetStatus() PaymentRequest_Status {
	if m != nil {
		return m.Status
	}
	return PaymentRequest_PENDING
}

type GenerateRedirRequest struct {
	Gateway              string   `protobuf:"bytes,1,opt,name=gateway,proto3" json:"gateway,omitempty"`
	Amount               int64    `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Currency             string   `protobuf:"bytes,3,opt,name=currency,proto3" json:"currency,omitempty"`
	OrderID              string   `protobuf:"bytes,4,opt,name=orderID,proto3" json:"orderID,omitempty"`
	Mobile               string   `protobuf:"bytes,5,opt,name=mobile,proto3" json:"mobile,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateRedirRequest) Reset()         { *m = GenerateRedirRequest{} }
func (m *GenerateRedirRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateRedirRequest) ProtoMessage()    {}
func (*GenerateRedirRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{6}
}

func (m *GenerateRedirRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateRedirRequest.Unmarshal(m, b)
}
func (m *GenerateRedirRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateRedirRequest.Marshal(b, m, deterministic)
}
func (m *GenerateRedirRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateRedirRequest.Merge(m, src)
}
func (m *GenerateRedirRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateRedirRequest.Size(m)
}
func (m *GenerateRedirRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateRedirRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateRedirRequest proto.InternalMessageInfo

func (m *GenerateRedirRequest) GetGateway() string {
	if m != nil {
		return m.Gateway
	}
	return ""
}

func (m *GenerateRedirRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *GenerateRedirRequest) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *GenerateRedirRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *GenerateRedirRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

type GenerateRedirResponse struct {
	CallbackUrl          string   `protobuf:"bytes,1,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
	InvoiceId            int64    `protobuf:"varint,2,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id,omitempty"`
	PaymentId            string   `protobuf:"bytes,3,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateRedirResponse) Reset()         { *m = GenerateRedirResponse{} }
func (m *GenerateRedirResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateRedirResponse) ProtoMessage()    {}
func (*GenerateRedirResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{7}
}

func (m *GenerateRedirResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateRedirResponse.Unmarshal(m, b)
}
func (m *GenerateRedirResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateRedirResponse.Marshal(b, m, deterministic)
}
func (m *GenerateRedirResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateRedirResponse.Merge(m, src)
}
func (m *GenerateRedirResponse) XXX_Size() int {
	return xxx_messageInfo_GenerateRedirResponse.Size(m)
}
func (m *GenerateRedirResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateRedirResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateRedirResponse proto.InternalMessageInfo

func (m *GenerateRedirResponse) GetCallbackUrl() string {
	if m != nil {
		return m.CallbackUrl
	}
	return ""
}

func (m *GenerateRedirResponse) GetInvoiceId() int64 {
	if m != nil {
		return m.InvoiceId
	}
	return 0
}

func (m *GenerateRedirResponse) GetPaymentId() string {
	if m != nil {
		return m.PaymentId
	}
	return ""
}

type ListGatewaysRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListGatewaysRequest) Reset()         { *m = ListGatewaysRequest{} }
func (m *ListGatewaysRequest) String() string { return proto.CompactTextString(m) }
func (*ListGatewaysRequest) ProtoMessage()    {}
func (*ListGatewaysRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{8}
}

func (m *ListGatewaysRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListGatewaysRequest.Unmarshal(m, b)
}
func (m *ListGatewaysRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListGatewaysRequest.Marshal(b, m, deterministic)
}
func (m *ListGatewaysRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListGatewaysRequest.Merge(m, src)
}
func (m *ListGatewaysRequest) XXX_Size() int {
	return xxx_messageInfo_ListGatewaysRequest.Size(m)
}
func (m *ListGatewaysRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListGatewaysRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListGatewaysRequest proto.InternalMessageInfo

type ListGatewaysResponse struct {
	Gateways             []*ListGatewaysResponse_GateWay `protobuf:"bytes,1,rep,name=gateways,proto3" json:"gateways,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *ListGatewaysResponse) Reset()         { *m = ListGatewaysResponse{} }
func (m *ListGatewaysResponse) String() string { return proto.CompactTextString(m) }
func (*ListGatewaysResponse) ProtoMessage()    {}
func (*ListGatewaysResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{9}
}

func (m *ListGatewaysResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListGatewaysResponse.Unmarshal(m, b)
}
func (m *ListGatewaysResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListGatewaysResponse.Marshal(b, m, deterministic)
}
func (m *ListGatewaysResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListGatewaysResponse.Merge(m, src)
}
func (m *ListGatewaysResponse) XXX_Size() int {
	return xxx_messageInfo_ListGatewaysResponse.Size(m)
}
func (m *ListGatewaysResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListGatewaysResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListGatewaysResponse proto.InternalMessageInfo

func (m *ListGatewaysResponse) GetGateways() []*ListGatewaysResponse_GateWay {
	if m != nil {
		return m.Gateways
	}
	return nil
}

type ListGatewaysResponse_GateWay struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	LogoImageAddress     string   `protobuf:"bytes,3,opt,name=logo_image_address,json=logoImageAddress,proto3" json:"logo_image_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListGatewaysResponse_GateWay) Reset()         { *m = ListGatewaysResponse_GateWay{} }
func (m *ListGatewaysResponse_GateWay) String() string { return proto.CompactTextString(m) }
func (*ListGatewaysResponse_GateWay) ProtoMessage()    {}
func (*ListGatewaysResponse_GateWay) Descriptor() ([]byte, []int) {
	return fileDescriptor_73578cb826d5ffb6, []int{9, 0}
}

func (m *ListGatewaysResponse_GateWay) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListGatewaysResponse_GateWay.Unmarshal(m, b)
}
func (m *ListGatewaysResponse_GateWay) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListGatewaysResponse_GateWay.Marshal(b, m, deterministic)
}
func (m *ListGatewaysResponse_GateWay) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListGatewaysResponse_GateWay.Merge(m, src)
}
func (m *ListGatewaysResponse_GateWay) XXX_Size() int {
	return xxx_messageInfo_ListGatewaysResponse_GateWay.Size(m)
}
func (m *ListGatewaysResponse_GateWay) XXX_DiscardUnknown() {
	xxx_messageInfo_ListGatewaysResponse_GateWay.DiscardUnknown(m)
}

var xxx_messageInfo_ListGatewaysResponse_GateWay proto.InternalMessageInfo

func (m *ListGatewaysResponse_GateWay) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ListGatewaysResponse_GateWay) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *ListGatewaysResponse_GateWay) GetLogoImageAddress() string {
	if m != nil {
		return m.LogoImageAddress
	}
	return ""
}

func init() {
	proto.RegisterEnum("payment_gateway.PaymentRequest_Status", PaymentRequest_Status_name, PaymentRequest_Status_value)
	proto.RegisterType((*MPGValidateRequest)(nil), "payment_gateway.MPGValidateRequest")
	proto.RegisterType((*MPGValidateResponse)(nil), "payment_gateway.MPGValidateResponse")
	proto.RegisterType((*MPGStartRequest)(nil), "payment_gateway.MPGStartRequest")
	proto.RegisterType((*MPGStartResponse)(nil), "payment_gateway.MPGStartResponse")
	proto.RegisterType((*GetPaymentResultByOrderIdRequest)(nil), "payment_gateway.GetPaymentResultByOrderIdRequest")
	proto.RegisterType((*PaymentRequest)(nil), "payment_gateway.PaymentRequest")
	proto.RegisterType((*GenerateRedirRequest)(nil), "payment_gateway.GenerateRedirRequest")
	proto.RegisterType((*GenerateRedirResponse)(nil), "payment_gateway.GenerateRedirResponse")
	proto.RegisterType((*ListGatewaysRequest)(nil), "payment_gateway.ListGatewaysRequest")
	proto.RegisterType((*ListGatewaysResponse)(nil), "payment_gateway.ListGatewaysResponse")
	proto.RegisterType((*ListGatewaysResponse_GateWay)(nil), "payment_gateway.ListGatewaysResponse.GateWay")
}

func init() { proto.RegisterFile("payment-gateway.proto", fileDescriptor_73578cb826d5ffb6) }

var fileDescriptor_73578cb826d5ffb6 = []byte{
	// 743 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x55, 0x6d, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x93, 0x34, 0x49, 0x27, 0x55, 0x1b, 0x4d, 0x5b, 0x64, 0x22, 0x0a, 0xa9, 0x69, 0xab,
	0x08, 0xb5, 0x41, 0x94, 0xbf, 0x08, 0xa9, 0x34, 0x25, 0x58, 0x6a, 0x4a, 0x70, 0x08, 0x48, 0x48,
	0x28, 0xda, 0xda, 0x4b, 0xb0, 0xea, 0xd8, 0xad, 0x77, 0x4d, 0x89, 0xc4, 0x01, 0x38, 0x01, 0x27,
	0xe2, 0x20, 0x5c, 0x83, 0x7f, 0xc8, 0xeb, 0x75, 0x62, 0xe7, 0xa3, 0xe1, 0x9f, 0xe7, 0xed, 0xac,
	0xf7, 0xcd, 0x9b, 0xb7, 0xb3, 0xb0, 0x7d, 0x4d, 0x46, 0x43, 0xea, 0xf2, 0xa3, 0x01, 0xe1, 0xf4,
	0x96, 0x8c, 0x1a, 0xd7, 0xbe, 0xc7, 0x3d, 0xdc, 0x90, 0x70, 0x5f, 0xc2, 0xda, 0x1f, 0x05, 0xb0,
	0xdd, 0x69, 0x7d, 0x20, 0x8e, 0x6d, 0x11, 0x4e, 0x0d, 0x7a, 0x13, 0x50, 0xc6, 0x51, 0x85, 0x62,
	0x9b, 0x32, 0x46, 0x06, 0x54, 0x55, 0x6a, 0x4a, 0x7d, 0xd5, 0x88, 0x43, 0xd4, 0x60, 0xed, 0x8d,
	0xc7, 0xb8, 0x41, 0xd9, 0xb5, 0xe7, 0x32, 0xaa, 0x66, 0xc5, 0x72, 0x0a, 0xc3, 0x27, 0x50, 0x49,
	0xc6, 0x5d, 0x7b, 0xe0, 0xaa, 0x39, 0x91, 0x37, 0x83, 0xe3, 0x43, 0x80, 0x2e, 0x27, 0x3c, 0x60,
	0xa7, 0x9e, 0x45, 0xd5, 0x7c, 0x4d, 0xa9, 0xe7, 0x8c, 0x04, 0x12, 0x9e, 0xd7, 0x73, 0xed, 0x9b,
	0x80, 0xbe, 0xf7, 0x89, 0xab, 0x37, 0xd5, 0x15, 0x91, 0x91, 0xc2, 0xf0, 0x01, 0xac, 0x76, 0xa2,
	0xba, 0xf4, 0xa6, 0x5a, 0x14, 0x07, 0x4d, 0x00, 0xed, 0x29, 0x6c, 0xa6, 0x2a, 0x94, 0x24, 0x55,
	0x28, 0x76, 0x03, 0xd3, 0xa4, 0x8c, 0x89, 0x12, 0x4b, 0x46, 0x1c, 0x6a, 0xb7, 0xb0, 0xd1, 0xee,
	0xb4, 0xba, 0x9c, 0xf8, 0x3c, 0xd6, 0xe3, 0x1e, 0x14, 0xc8, 0xd0, 0x0b, 0x5c, 0x2e, 0x72, 0x73,
	0x86, 0x8c, 0xb0, 0x0a, 0x25, 0x33, 0xf0, 0x7d, 0xea, 0x9a, 0x23, 0xa9, 0xc4, 0x38, 0x0e, 0x0f,
	0xf0, 0x7c, 0x8b, 0xfa, 0x7a, 0x53, 0x16, 0x1f, 0x87, 0xe1, 0xdf, 0x86, 0xde, 0xa5, 0xed, 0x44,
	0xf5, 0xae, 0x1a, 0x32, 0xd2, 0x7e, 0x40, 0x65, 0x72, 0xb0, 0xa4, 0x59, 0x83, 0xf2, 0x57, 0xa1,
	0x99, 0x20, 0x22, 0xbb, 0x91, 0x84, 0xb0, 0x0e, 0x1b, 0x89, 0x50, 0x88, 0x1d, 0x51, 0x99, 0x86,
	0x43, 0x9d, 0x64, 0xff, 0x75, 0x4b, 0x72, 0x9a, 0x00, 0xda, 0x0b, 0xa8, 0xb5, 0x28, 0x97, 0xba,
	0x19, 0x94, 0x05, 0x0e, 0x7f, 0x35, 0x7a, 0x2b, 0x28, 0x5b, 0x09, 0x5f, 0xc4, 0x35, 0x29, 0xa9,
	0x9a, 0xb4, 0x9f, 0x59, 0x58, 0x1f, 0xef, 0x5d, 0x92, 0x9c, 0x26, 0x92, 0x9d, 0x22, 0x12, 0xae,
	0xda, 0xee, 0x37, 0xcf, 0x36, 0xa9, 0xa4, 0x99, 0x33, 0x26, 0x40, 0xa2, 0x15, 0xf9, 0x99, 0x56,
	0x10, 0xdf, 0x6a, 0x13, 0x76, 0x25, 0x4c, 0x12, 0xb6, 0x42, 0xc6, 0xf8, 0x12, 0x0a, 0x4c, 0x58,
	0x4a, 0x2d, 0xd4, 0x94, 0xfa, 0xfa, 0xf1, 0x41, 0x63, 0xea, 0x1e, 0x34, 0xd2, 0xd4, 0x1b, 0x91,
	0x01, 0x0d, 0xb9, 0x4b, 0x3b, 0x84, 0x42, 0x84, 0x60, 0x19, 0x8a, 0x9d, 0xb3, 0x8b, 0xa6, 0x7e,
	0xd1, 0xaa, 0x64, 0xc2, 0xa0, 0xdb, 0x3b, 0x3d, 0x3d, 0xeb, 0x76, 0x2b, 0x0a, 0x96, 0x20, 0xff,
	0xfa, 0x44, 0x3f, 0xaf, 0x64, 0xb5, 0x5f, 0x0a, 0x6c, 0xb5, 0xa8, 0x4b, 0x7d, 0x61, 0x37, 0xcb,
	0xf6, 0x13, 0x82, 0xc8, 0xf3, 0x62, 0x41, 0x64, 0x98, 0x28, 0x2a, 0xbb, 0xd0, 0x5f, 0xb9, 0xc5,
	0xfe, 0xca, 0x2f, 0xf2, 0xd7, 0x4a, 0xca, 0x5f, 0xdf, 0x61, 0x7b, 0x8a, 0x97, 0x34, 0xd9, 0x2e,
	0xac, 0x99, 0xc4, 0x71, 0x2e, 0x89, 0x79, 0xd5, 0x0f, 0x7c, 0x27, 0x76, 0x59, 0x8c, 0xf5, 0x7c,
	0x07, 0x77, 0x00, 0x64, 0x0f, 0xfa, 0xb6, 0x25, 0x59, 0x26, 0xba, 0xb2, 0x03, 0x10, 0x4b, 0x6a,
	0xcf, 0xf1, 0xd6, 0x36, 0x6c, 0x9e, 0xdb, 0x8c, 0xb7, 0xa2, 0x72, 0x99, 0x14, 0x44, 0xfb, 0xad,
	0xc0, 0x56, 0x1a, 0x97, 0x84, 0x74, 0x28, 0x49, 0x69, 0xc2, 0xdb, 0x99, 0xab, 0x97, 0x8f, 0x8f,
	0x66, 0x5a, 0x36, 0x6f, 0x63, 0x23, 0x04, 0x3e, 0x92, 0x91, 0x31, 0xde, 0x5e, 0x25, 0x50, 0x94,
	0x20, 0x22, 0xe4, 0x5d, 0x32, 0x8c, 0x47, 0x9a, 0xf8, 0xc6, 0x2d, 0x58, 0xe1, 0x36, 0x77, 0xe2,
	0x41, 0x16, 0x05, 0x78, 0x08, 0xe8, 0x78, 0x03, 0xaf, 0x6f, 0x0f, 0xc9, 0x80, 0xf6, 0x89, 0x65,
	0xf9, 0xe1, 0x9c, 0x90, 0x33, 0x2c, 0x5c, 0xd1, 0xc3, 0x85, 0x93, 0x08, 0x3f, 0xfe, 0x9b, 0x1b,
	0x7b, 0x5f, 0x12, 0xc2, 0x2f, 0xb0, 0x99, 0x92, 0x9a, 0x9a, 0xbc, 0x67, 0x9c, 0xe3, 0xfe, 0x4c,
	0x15, 0xf3, 0x8c, 0x52, 0x3d, 0x58, 0x96, 0x16, 0x55, 0xab, 0x65, 0xf0, 0x33, 0xac, 0x25, 0x75,
	0xc0, 0xbd, 0x25, 0x32, 0x45, 0xff, 0xdf, 0xff, 0x2f, 0x31, 0xb5, 0x0c, 0xde, 0xc0, 0xfd, 0x45,
	0x33, 0xa1, 0x89, 0xcf, 0xe6, 0xb0, 0xbc, 0x7b, 0x7e, 0x54, 0x1f, 0x2d, 0xb9, 0x78, 0x5a, 0x06,
	0xdf, 0x41, 0x29, 0x1e, 0x82, 0x58, 0x9b, 0x49, 0x9f, 0x1a, 0xcc, 0xd5, 0xdd, 0x3b, 0x32, 0xc6,
	0x55, 0x7c, 0x82, 0x72, 0xe2, 0x05, 0xc0, 0xc7, 0xf3, 0xf6, 0x4c, 0xbd, 0x80, 0xd5, 0xbd, 0xbb,
	0x93, 0xe2, 0x7f, 0x5f, 0x16, 0xc4, 0xc3, 0xfa, 0xfc, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x54,
	0x93, 0x25, 0xbf, 0x71, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PaymentGatewayClient is the client API for PaymentGateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PaymentGatewayClient interface {
	GenerateRedirectURL(ctx context.Context, in *GenerateRedirRequest, opts ...grpc.CallOption) (*GenerateRedirResponse, error)
	ListGateways(ctx context.Context, in *ListGatewaysRequest, opts ...grpc.CallOption) (*ListGatewaysResponse, error)
	GetPaymentResultByOrderID(ctx context.Context, in *GetPaymentResultByOrderIdRequest, opts ...grpc.CallOption) (*PaymentRequest, error)
	MPGStart(ctx context.Context, in *MPGStartRequest, opts ...grpc.CallOption) (*MPGStartResponse, error)
	MPGValidate(ctx context.Context, in *MPGValidateRequest, opts ...grpc.CallOption) (*MPGValidateResponse, error)
}

type paymentGatewayClient struct {
	cc *grpc.ClientConn
}

func NewPaymentGatewayClient(cc *grpc.ClientConn) PaymentGatewayClient {
	return &paymentGatewayClient{cc}
}

func (c *paymentGatewayClient) GenerateRedirectURL(ctx context.Context, in *GenerateRedirRequest, opts ...grpc.CallOption) (*GenerateRedirResponse, error) {
	out := new(GenerateRedirResponse)
	err := c.cc.Invoke(ctx, "/payment_gateway.PaymentGateway/GenerateRedirectURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentGatewayClient) ListGateways(ctx context.Context, in *ListGatewaysRequest, opts ...grpc.CallOption) (*ListGatewaysResponse, error) {
	out := new(ListGatewaysResponse)
	err := c.cc.Invoke(ctx, "/payment_gateway.PaymentGateway/ListGateways", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentGatewayClient) GetPaymentResultByOrderID(ctx context.Context, in *GetPaymentResultByOrderIdRequest, opts ...grpc.CallOption) (*PaymentRequest, error) {
	out := new(PaymentRequest)
	err := c.cc.Invoke(ctx, "/payment_gateway.PaymentGateway/GetPaymentResultByOrderID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentGatewayClient) MPGStart(ctx context.Context, in *MPGStartRequest, opts ...grpc.CallOption) (*MPGStartResponse, error) {
	out := new(MPGStartResponse)
	err := c.cc.Invoke(ctx, "/payment_gateway.PaymentGateway/MPGStart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentGatewayClient) MPGValidate(ctx context.Context, in *MPGValidateRequest, opts ...grpc.CallOption) (*MPGValidateResponse, error) {
	out := new(MPGValidateResponse)
	err := c.cc.Invoke(ctx, "/payment_gateway.PaymentGateway/MPGValidate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentGatewayServer is the server API for PaymentGateway service.
type PaymentGatewayServer interface {
	GenerateRedirectURL(context.Context, *GenerateRedirRequest) (*GenerateRedirResponse, error)
	ListGateways(context.Context, *ListGatewaysRequest) (*ListGatewaysResponse, error)
	GetPaymentResultByOrderID(context.Context, *GetPaymentResultByOrderIdRequest) (*PaymentRequest, error)
	MPGStart(context.Context, *MPGStartRequest) (*MPGStartResponse, error)
	MPGValidate(context.Context, *MPGValidateRequest) (*MPGValidateResponse, error)
}

// UnimplementedPaymentGatewayServer can be embedded to have forward compatible implementations.
type UnimplementedPaymentGatewayServer struct {
}

func (*UnimplementedPaymentGatewayServer) GenerateRedirectURL(ctx context.Context, req *GenerateRedirRequest) (*GenerateRedirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateRedirectURL not implemented")
}
func (*UnimplementedPaymentGatewayServer) ListGateways(ctx context.Context, req *ListGatewaysRequest) (*ListGatewaysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGateways not implemented")
}
func (*UnimplementedPaymentGatewayServer) GetPaymentResultByOrderID(ctx context.Context, req *GetPaymentResultByOrderIdRequest) (*PaymentRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaymentResultByOrderID not implemented")
}
func (*UnimplementedPaymentGatewayServer) MPGStart(ctx context.Context, req *MPGStartRequest) (*MPGStartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MPGStart not implemented")
}
func (*UnimplementedPaymentGatewayServer) MPGValidate(ctx context.Context, req *MPGValidateRequest) (*MPGValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MPGValidate not implemented")
}

func RegisterPaymentGatewayServer(s *grpc.Server, srv PaymentGatewayServer) {
	s.RegisterService(&_PaymentGateway_serviceDesc, srv)
}

func _PaymentGateway_GenerateRedirectURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateRedirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).GenerateRedirectURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment_gateway.PaymentGateway/GenerateRedirectURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).GenerateRedirectURL(ctx, req.(*GenerateRedirRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentGateway_ListGateways_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGatewaysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).ListGateways(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment_gateway.PaymentGateway/ListGateways",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).ListGateways(ctx, req.(*ListGatewaysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentGateway_GetPaymentResultByOrderID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPaymentResultByOrderIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).GetPaymentResultByOrderID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment_gateway.PaymentGateway/GetPaymentResultByOrderID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).GetPaymentResultByOrderID(ctx, req.(*GetPaymentResultByOrderIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentGateway_MPGStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MPGStartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).MPGStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment_gateway.PaymentGateway/MPGStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).MPGStart(ctx, req.(*MPGStartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentGateway_MPGValidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MPGValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).MPGValidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment_gateway.PaymentGateway/MPGValidate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).MPGValidate(ctx, req.(*MPGValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PaymentGateway_serviceDesc = grpc.ServiceDesc{
	ServiceName: "payment_gateway.PaymentGateway",
	HandlerType: (*PaymentGatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateRedirectURL",
			Handler:    _PaymentGateway_GenerateRedirectURL_Handler,
		},
		{
			MethodName: "ListGateways",
			Handler:    _PaymentGateway_ListGateways_Handler,
		},
		{
			MethodName: "GetPaymentResultByOrderID",
			Handler:    _PaymentGateway_GetPaymentResultByOrderID_Handler,
		},
		{
			MethodName: "MPGStart",
			Handler:    _PaymentGateway_MPGStart_Handler,
		},
		{
			MethodName: "MPGValidate",
			Handler:    _PaymentGateway_MPGValidate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment-gateway.proto",
}
