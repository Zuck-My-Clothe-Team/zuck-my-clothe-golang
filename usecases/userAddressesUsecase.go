package usecases

import (
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/google/uuid"
)

type userAddressesUsecase struct {
	userAddressesRepository model.UserAddressesRepository
}

func CreateNewUserAddressesUsecase(userAddressesRepository model.UserAddressesRepository) model.UserAddressesUsecase {
	return &userAddressesUsecase{userAddressesRepository: userAddressesRepository}
}

func toUserAddressDetail(userAddress *model.UserAddresses) interface{} {
	result := model.UserAddressDetail{
		AddressID:   userAddress.AddressID,
		UserID:      userAddress.UserID,
		Address:     userAddress.Address,
		Province:    userAddress.Province,
		District:    userAddress.District,
		SubDistrict: userAddress.SubDistrict,
		Zipcode:     userAddress.Zipcode,
		Lat:         userAddress.Lat,
		Long:        userAddress.Long,
	}
	return result
}

func (u *userAddressesUsecase) AddUserAddress(owenrID string, newUserAddress *model.AddUserAddressDTO) (*interface{}, error) {
	data := model.UserAddresses{
		AddressID:   uuid.New().String(),
		UserID:      owenrID,
		Address:     newUserAddress.Address,
		Province:    newUserAddress.Province,
		District:    newUserAddress.District,
		SubDistrict: newUserAddress.SubDistrict,
		Zipcode:     newUserAddress.Zipcode,
		Lat:         newUserAddress.Lat,
		Long:        newUserAddress.Long,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := u.userAddressesRepository.AddUserAddress(&data); err != nil {
		return nil, err
	}
	createdRecord, err := u.userAddressesRepository.FindUserAddressByID(data.AddressID, owenrID)
	if err != nil {
		return nil, err
	}
	createdRecordDetail := toUserAddressDetail(createdRecord)
	return &createdRecordDetail, nil
}

func (u *userAddressesUsecase) FindUserAddressByID(addressID string, userID string) (*interface{}, error) {
	userAddress, err := u.userAddressesRepository.FindUserAddressByID(addressID, userID)
	if err != nil {
		return nil, err
	}
	var result interface{} = toUserAddressDetail(userAddress)
	return &result, nil
}

func (u *userAddressesUsecase) FindUserAddresByOwnerID(ownerID string) (*[]interface{}, error) {
	addressList, err := u.userAddressesRepository.FindUserAddresByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	var result []interface{}

	if len((*addressList)) == 0 {
		return &result, err
	}

	for _, address := range *addressList {
		result = append(result, toUserAddressDetail(&address))
	}
	return &result, err
}

func (u *userAddressesUsecase) UpdateUserAddressData(userID string, updatedAddressData *model.AddUserAddressDTO) (*interface{}, error) {
	data := model.UpdateUserAddressDTO{
		Address:     updatedAddressData.Address,
		Province:    updatedAddressData.Province,
		District:    updatedAddressData.District,
		SubDistrict: updatedAddressData.SubDistrict,
		Zipcode:     updatedAddressData.Zipcode,
		Lat:         updatedAddressData.Lat,
		Long:        updatedAddressData.Long,
		UpdatedAt:   time.Now().UTC(),
	}
	if err := u.userAddressesRepository.UpdateUserAddressData(userID, updatedAddressData.AddressID, &data); err != nil {
		return nil, err
	}
	updatedData, err := u.userAddressesRepository.FindUserAddressByID(updatedAddressData.AddressID, userID)
	if err != nil {
		return nil, err
	}
	result := toUserAddressDetail(updatedData)
	return &result, nil
}

func (u *userAddressesUsecase) DeleteUserAddress(userID string, addressID string) error {
	err := u.userAddressesRepository.DeleteUserAddress(userID, addressID)
	return err
}
