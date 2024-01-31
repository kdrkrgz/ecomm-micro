package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

/*
	go test -run api_test.go
*/

type paymentMock struct {
	mock.Mock
}

func (p *paymentMock) Charge(ctx context.Context, order *domain.Order) error {
	args := p.Called(ctx, order)
	return args.Error(0)
}

type dbMock struct {
	mock.Mock
}

func (d *dbMock) Save(ctx context.Context, order *domain.Order) error {
	args := d.Called(ctx, order)
	return args.Error(0)
}

func (d *dbMock) Get(ctx context.Context, orderId int64) (domain.Order, error) {
	args := d.Called(ctx, orderId)
	return args.Get(0).(domain.Order), args.Error(1)
}

func orderObj() domain.Order {
	return domain.Order{
		CustomerId: 1,
		OrderItems: []domain.OrderItem{
			{
				Id:          1,
				ProductCode: "P1",
				Quantity:    1,
				UnitPrice:   10,
			},
			{
				Id:          2,
				ProductCode: "P2",
				Quantity:    2,
				UnitPrice:   20,
			},
		},
		//CreatedAt: 0,
	}
}

func TestPlaceOrder(t *testing.T) {
	payment := &paymentMock{}
	db := &dbMock{}

	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), orderObj())
	fmt.Println(err)
	assert.Nil(t, err)
}

func Test_Should_Return_Error_When_Db_Persistence_Fail(t *testing.T) {
	payment := new(paymentMock)
	db := new(dbMock)

	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), orderObj())
	assert.NotNil(t, err)
	assert.EqualError(t, err, "db error")
}

func Test_Should_Return_Error_When_Payment_Fail(t *testing.T) {
	payment := new(paymentMock)
	db := new(dbMock)

	payment.On("Charge", mock.Anything, mock.Anything).Return(errors.New("payment error"))
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), orderObj())
	st, _ := status.FromError(err)
	assert.Equal(t, st.Code(), codes.InvalidArgument)
	assert.Equal(t, st.Message(), "payment error")
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = payment error")
}
