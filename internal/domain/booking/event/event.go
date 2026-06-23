package event

import "encoding/json"

// PaymentMessage là payload của topic payment.request.
type PaymentMessage struct {
	BookingID uint64 `json:"booking_id"`
}

func (m PaymentMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func ParsePaymentMessage(data []byte) (PaymentMessage, error) {
	var m PaymentMessage
	err := json.Unmarshal(data, &m)
	return m, err
}
