package service

import (
	"errors"

	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
)

// OrderService реализация бизнес логики обработки номера заказа
type BalanceService struct {
	repo repository.Balance
}

//NewOrderService конструктор
func NewBalanceService(repo repository.Balance) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

// Create создание нового баланса
// *(используется при авторизации нового пользователя)
func (b *BalanceService) Create(userID int) error {
	return b.repo.Create(userID)
}

// Get получить текущее количество баллов клиента и общее количество использованных баллов за все время
func (b *BalanceService) Get(userID int) (*core.Balance, error) {
	return b.repo.Get(userID)
}

// UpdateCurrent обновление текущего количества баллов клиента
// *получаю из базы имеющиеся баллы и добавляю к ним баллы из нового заказа
func (b *BalanceService) UpdateCurrent(userID int, order *core.Order) error {
	userBalance, err := b.repo.Get(userID)
	if err != nil {
		return err
	}
	userBalance.Current += order.Accrual
	return b.repo.UpdateCurrent(userBalance)
}

// Withdrawal списывую баллы с  баланса клиента и увеличивую общее количество использованных баллов за все время
func (b *BalanceService) Withdrawal(userID int, orderNumber float64) error {
	userBalance, err := b.repo.Get(userID)
	if err != nil {
		return err
	}
	userBalance.Current -= orderNumber
	userBalance.Withdrawn += orderNumber

	if userBalance.Current < 0 {
		return errors.New("there are not enough points on the balance")
	}

	if err := b.repo.UpdateCurrent(userBalance); err != nil {
		return err
	}
	if err := b.repo.UpdateWithdrawn(userBalance); err != nil {
		return err
	}
	return nil
}
