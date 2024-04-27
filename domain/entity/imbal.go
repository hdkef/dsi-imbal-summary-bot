package entity

import "time"

type Imbal struct {
	id     uint32
	date   time.Time
	amount float64
}

func (i *Imbal) SetID(id uint32) {
	i.id = id
}

func (i *Imbal) GetID() uint32 {
	return i.id
}

func (i *Imbal) SetDate(date time.Time) {
	i.date = date
}

func (i *Imbal) GetDate() time.Time {
	return i.date
}

func (i *Imbal) SetAmount(amount float64) {
	i.amount = amount
}

func (i *Imbal) GetAmount() float64 {
	return i.amount
}
