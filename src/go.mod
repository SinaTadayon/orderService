module gitlab.faza.io/order-project/order-service

go 1.13

require (
	github.com/Netflix/go-env v0.0.0-20180529183433-1e80ef5003ef
	github.com/devfeel/mapper v0.7.2
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	gitlab.faza.io/go-framework/acl v0.0.3
	gitlab.faza.io/go-framework/logger v0.0.3
	gitlab.faza.io/go-framework/mongoadapter v0.0.8
	gitlab.faza.io/protos/cart v0.0.12
	gitlab.faza.io/protos/notification v0.0.3
	gitlab.faza.io/protos/order v0.0.51
	gitlab.faza.io/protos/payment-gateway v0.0.8
	gitlab.faza.io/protos/stock-proto.git v0.0.8
	gitlab.faza.io/services/user-app-client v0.0.20
	go.mongodb.org/mongo-driver v1.2.0
	google.golang.org/grpc v1.26.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
