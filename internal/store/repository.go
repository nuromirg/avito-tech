package store

import (
	"avito_task/internal/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

type UserRepository struct {
	 Storage *Storage
}

// I'll use msql stored procedures in future commits

func (ur *UserRepository) CreateUser() (*model.User, error) {
	u := &model.User{}
	if err := ur.Storage.Db.QueryRow("INSERT INTO users (balance) VALUE (?)", 0,
		).Err(); err != nil {
		return nil, err
	}
	if err := ur.Storage.Db.QueryRow("SELECT id, balance FROM users ORDER BY id DESC LIMIT 1",
		).Scan(&u.Id, &u.Balance); err != nil {
		return nil, err
	}
	err := ur.Storage.r.Set(ur.Storage.r.Context(), strconv.Itoa(u.Id),
		strconv.FormatInt(u.Balance, 10), 0).Err()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) FindById(id int) (*model.User, error) {
	if id <= 0 {
		logrus.Errorf("Error: Id can't be 0")
		return nil, nil
	}

	balance, err := ur.Storage.r.Get(ur.Storage.r.Context(), strconv.Itoa(id)).Result()
	if err != nil {
		return nil, err
	}
	balanceInt64, _ := strconv.ParseInt(balance, 10, 64)
	u := &model.User{
		Id:			id,
		Balance: 	balanceInt64,
	}
	return u, nil
}

func (ur *UserRepository) GetBalanceById(id int) (int64, error) {
	u, err := ur.FindById(id)
	if err != nil {
		logrus.Error(err)
		return -1, err
	}

	return u.Balance, nil
}

// Метод начисления/списания средств на баланс. Принимает id пользователя и сколько средств зачислить.

func (ur *UserRepository) ChangeFunds(id int, sum int64) (*model.User, error) {
	u, err := ur.FindById(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if sum < 0 && u.Balance < sum {
		return nil, fmt.Errorf("not enough funds to withdraw")
	}
	u.Balance += sum
	if err := ur.Storage.r.Set(ur.Storage.r.Context(), strconv.Itoa(u.Id), u.Balance, 0).Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	_, err = ur.Storage.Db.Query("UPDATE users SET balance = ? WHERE id = ?", &u.Balance, &u.Id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return u, nil
}

// Метод перевода средств от пользователя к пользователю.
// Принимает id пользователя с которого нужно списать средства,
// id пользователя которому должны зачислить средства, а также сумму.

func (ur *UserRepository) TransactionFunds(t1idFrom, t2idTo int, sum int64) error {
	u1From, _ := ur.FindById(t1idFrom)
	u2To, _ := ur.FindById(t2idTo)
	if u1From.Balance >= sum {
		u1From.Balance -= sum
	} else {
		return fmt.Errorf("error: Not enough funds: from user")
	}
	u2To.Balance += sum
	if err := ur.Storage.r.Set(ur.Storage.r.Context(), strconv.Itoa(u1From.Id), u1From.Balance, 0).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	_, err := ur.Storage.Db.Query("UPDATE users SET balance = ? WHERE id = ?", &u1From.Balance, &u1From.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err := ur.Storage.r.Set(ur.Storage.r.Context(), strconv.Itoa(u2To.Id), u2To.Balance, 0).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	_, err = ur.Storage.Db.Query("UPDATE users SET balance = ? WHERE id = ?", &u2To.Balance, &u2To.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

