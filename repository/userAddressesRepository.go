package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type userAddressReopository struct {
	db *platform.Postgres
}

func CreateNewUserAddressesRepository(db *platform.Postgres) model.UserAddressesRepository {
	return &userAddressReopository{db: db}
}

func (u *userAddressReopository) AddUserAddress(newUserAddress *model.UserAddresses) error {
	dbTx := u.db.Create(newUserAddress)
	return dbTx.Error
}

func (u *userAddressReopository) FindUserAddressByID(addressID string) (*model.UserAddresses, error) {
	data := new(model.UserAddresses)
	dbTx := u.db.First(data, "address_id = ?", addressID)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return data, nil
}

func (u *userAddressReopository) FindUserAddresByOwnerID(ownerID string) (*[]model.UserAddresses, error) {
	dataList := new([]model.UserAddresses)
	dbTx := u.db.Find(dataList, "user_id = ?", ownerID)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return dataList, nil
}

func (u *userAddressReopository) UpdateUserAddressData(userID string, addressID string, updatedAddressData *model.UpdateUserAddressDTO) error {
	dbTx := u.db.Table("UserAddresses").Where("user_id = ? AND address_id = ?", userID, addressID).Updates(updatedAddressData)

	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return dbTx.Error
}

func (u *userAddressReopository) DeleteUserAddress(userID string, addressID string) error {
	dbTx := u.db.Table("UserAddresses").Where("user_id = ? ", userID).Delete(&model.UserAddresses{AddressID: addressID})
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return dbTx.Error
}
