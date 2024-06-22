package data

import "database/sql"

type PaymentModel struct {
	DB *sql.DB
}

func (p *PaymentModel) CreateOrder(orderData *OrderData) error {
	query := `INSERT INTO payment_details (razorpay_order_id, amount, task_id, user_id, worker_id) VALUES ($1, $2, $3, $4, $5)`

	ctx, cancel := handlectx()
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, orderData.OrderID, orderData.Amount, orderData.TaskID, orderData.TaskOwnerID, orderData.WorkerID)
	if err != nil {
		return err
	}

	return nil
}
