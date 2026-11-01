package main

import (
	"errors"
	"fmt"

	"example.com/online-shop/product"
)

type Payment interface {
	Pay(amount float64) error
}

type CashPayment struct{}

type CardPayment struct{}

type OnlinePayment struct{}

func paymentProcess(p Payment, amount float64) error {
	return p.Pay(amount)
}

func (p *CashPayment) Pay(amount float64) error {
	err := validatePayment(amount)

	if err != nil {
		return err
	}

	fmt.Printf("Paying with cash: %.2f\n", amount)
	return nil
}

func (p *CardPayment) Pay(amount float64) error {
	err := validatePayment(amount)

	if err != nil {
		return err
	}

	fmt.Printf("Paying with card: %.2f\n", amount)
	return nil
}

func (p *OnlinePayment) Pay(amount float64) error {
	err := validatePayment(amount)

	if err != nil {
		return err
	}

	fmt.Printf("Paying online: %.2f\n", amount)
	return nil
}

func validatePayment(amount float64) error {
	if amount <= 0 {
		return errors.New("Payment should be positive")
	}

	return nil
}

func Purchase(p *product.Product, amountToBuy int, paymentType Payment) error {
	if amountToBuy <= 0 {
		return errors.New("Quantity should be positive")
	} else if amountToBuy > p.Stock {
		return errors.New("Not enough items in stock")
	}

	price := float64(amountToBuy) * p.Price
	err := paymentProcess(paymentType, price)

	if err != nil {
		return err
	}

	err = p.Sell(amountToBuy)

	if err != nil {
		return err
	}

	return nil
}
