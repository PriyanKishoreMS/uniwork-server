package data

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const (
	uniqueViolationErr = pq.ErrorCode("23505")
)

type PaymentModel struct {
	DB *sql.DB
}

func (p *PaymentModel) CreateOrder(orderData *OrderData) error {
	query := `INSERT INTO payment_details (razorpay_order_id, amount, task_owner_id, task_request_id) VALUES ($1, $2, $3, $4)`

	ctx, cancel := handlectx()
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, orderData.OrderID, orderData.Amount, orderData.TaskOwnerID, orderData.TaskRequestID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case uniqueViolationErr:
				return errors.New("order already exists")
			default:
				return err
			}
		}
	}
	return nil
}

func (p *PaymentModel) CheckExistingOrder(taskRequestId int64) *string {
	// TODO return razorpay_orderId
	return nil
}
