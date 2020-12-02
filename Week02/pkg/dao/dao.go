

package dao

import (
"time"

xerrors "github.com/pkg/errors"
)

type UserAttributes struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

func (UserAttributes) TableName() string {
	return "user"
}

type UserDAO interface {
	SelectByEmail(email string) (*UserAttributes, error)
	Save(user *UserAttributes) error
}

type UserDAOImpl struct {
}

func (userDAO *UserDAOImpl) SelectByEmail(email string) (*UserAttributes, error) {
	user := &UserAttributes{}
	err := db.Where("email = ?", email).First(user).Error

	return user, xerrors.Wrapf(err, "SelectByEmail error")
}

func (userDAO *UserDAOImpl) Save(user *UserAttributes) error {
	return db.Create(user).Error
}
