package usecases

import (
	"errors"

	"zuck-my-clothe/zuck-my-clothe-backend/model"
	repo "zuck-my-clothe/zuck-my-clothe-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderUsecase struct {
	orderDetailRepo repo.OrderDetailRepository
	orderHeaderRepo repo.OrderHeaderRepository
	userRepo        repo.UserRepository
	machineRepo     repo.MachineRepository
	paymentUsecase  model.PaymentUsecase
}

type OrderUsecase interface {
	CreateNewOrder(newOrder *model.NewOrder) (*model.FullOrder, error)
	GetAllHeader() (*[]model.OrderHeader, error)
	GetByHeaderID(orderHeaderID string, isAdminView bool, option string) (interface{}, error)
	GetByBranchID(branchID string, managerUserID string) ([]interface{}, error)
	GetByUserID(userID string) ([]interface{}, error)
	UpdateStatus(order model.UpdateOrder) (interface{}, error)
	UpdateReview(review model.OrderReview) (*model.FullOrder, error)
	SoftDelete(orderHeaderID string, deletedBy string) (*model.FullOrder, error)
}

func CreateOrderUsecase(orderHeaderRepository repo.OrderHeaderRepository, orderDetailRepository repo.OrderDetailRepository, userRepository repo.UserRepository, machineRepo repo.MachineRepository, paymentUsecase model.PaymentUsecase) OrderUsecase {
	return &orderUsecase{
		orderHeaderRepo: orderHeaderRepository,
		orderDetailRepo: orderDetailRepository,
		userRepo:        userRepository,
		machineRepo:     machineRepo,
		paymentUsecase:  paymentUsecase,
	}
}

func toUserDetailDTO(userModel *model.Users) *model.UserDetailDTO {
	detail := model.UserDetailDTO{
		UserID:          userModel.UserID,
		GoogleID:        userModel.GoogleID,
		Email:           userModel.Email,
		Phone:           userModel.Phone,
		FirstName:       userModel.FirstName,
		LastName:        userModel.LastName,
		ProfileImageURL: userModel.ProfileImageURL,
		Role:            userModel.Role,
	}
	return &detail
}

func serVicePriceMapper(weight int) int {
	var price int = 0
	switch weight {
	case 7:
		price = model.SevenKGMachinePrice
	case 14:
		price = model.FourteenKGMachinePrice
	case 21:
		price = model.TwentyOneKGMachinePrice
	}
	return price
}

func combineFullOrder(h *model.OrderHeader, d *[]model.OrderDetail, user model.Users, isAdminView bool) *model.FullOrder {
	fullOrder := model.FullOrder{
		OrderHeaderID:   h.OrderHeaderID,
		UserID:          h.UserID,
		UserDetail:      *toUserDetailDTO(&user),
		BranchID:        h.BranchID,
		OrderNote:       h.OrderNote,
		PaymentID:       h.PaymentID,
		ZuckOnsite:      h.ZuckOnsite,
		DeliveryAddress: h.DeliveryAddress,
		DeliveryLat:     h.DeliveryLat,
		DeliveryLong:    h.DeliveryLong,
		StarRating:      h.StarRating,
		ReviewComment:   h.ReviewComment,
		CreatedAt:       &h.CreatedAt,
		CreatedBy:       &h.CreatedBy,
		UpdatedAt:       &h.UpdatedAt,
		UpdatedBy:       &h.UpdatedBy,
		DeletedAt:       &h.DeletedAt,
		DeletedBy:       h.DeletedBy,
		OrderDetails:    *d,
	}

	if !isAdminView {
		fullOrder.DeletedAt = nil
		fullOrder.DeletedBy = nil
	}

	return &fullOrder
}

func (u *orderUsecase) CreateNewOrder(newOrder *model.NewOrder) (*model.FullOrder, error) {

	if !newOrder.ZuckOnsite &&
		(newOrder.DeliveryAddress == nil ||
			newOrder.DeliveryLat == nil ||
			newOrder.DeliveryLong == nil) {
		return nil, errors.New("no delivery address for online order")
	}

	// create a payment
	//payId := "2709093b-1ee9-44a1-bd60-e7b092012c8d"

	var allWashingWeight int16 = 0
	var allDryingweight int16 = 0

	var washingBasketCount int = 0
	var dryinBasketCount int = 0
	var basketWeight int16 = 0
	var dryingWeight int16 = 0

	var isDeliveryExist bool = false
	var isPickupExist bool = false
	var isAgentsExist bool = false

	for _, detail := range newOrder.OrderDetails {
		var serviceType model.ServiceType = detail.ServiceType
		if serviceType == "Washing" {
			allWashingWeight += detail.Weight
			washingBasketCount += 1
			basketWeight = detail.Weight
		} else if serviceType == "Drying" {
			dryingWeight = detail.Weight
			allDryingweight += detail.Weight
			dryinBasketCount += 1
		} else if serviceType == "Pickup" {
			isPickupExist = true
		} else if serviceType == "Delivery" {
			isDeliveryExist = true
		} else if serviceType == "Agents" {
			isAgentsExist = true
		}
	}

	// In case that customer select both washing and drying service
	// Front-end should combind all weight into single drying weight
	if washingBasketCount > 0 && dryinBasketCount > 1 {
		return nil, errors.New("ERR: number of drying request exceeded limit in this operation")
	}
	if washingBasketCount == 0 && isAgentsExist {
		return nil, errors.New("ERR: cannot use agent when you only drying your clothes")
	}

	if allWashingWeight > 0 && allWashingWeight > 21 {
		return nil, errors.New("ERR: washing weight exceded 21 Kg")
	} else if allWashingWeight == 0 && allDryingweight > 21 {
		return nil, errors.New("ERR: washing weight exceded 21 Kg")
	} else if allWashingWeight == 0 && allDryingweight == 0 {
		return nil, errors.New("ERR: empty order")
	}

	if isPickupExist != isDeliveryExist {
		return nil, errors.New("ERR: cannot select Pickup or Delivery individualy")
	} else if newOrder.ZuckOnsite == isPickupExist {
		if newOrder.ZuckOnsite {
			return nil, errors.New("ERR: cannot select Pickup or Delivery when using onsite service")
		} else {
			return nil, errors.New("ERR: Pickup and Delivery are needed when using online service")
		}
	}

	//Available machine validation
	var availableWasher *[]model.MachineInBranch
	var availableDryer *[]model.MachineInBranch
	if allWashingWeight > 0 {
		var err error
		availableWasher, err = u.machineRepo.GetMachineToAssign(newOrder.BranchID, "Washer", int(basketWeight), washingBasketCount)
		if err != nil {
			return nil, err
		}
	}
	if allDryingweight > 0 {
		var err error
		availableDryer, err = u.machineRepo.GetMachineToAssign(newOrder.BranchID, "Dryer", int(dryingWeight), dryinBasketCount)
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println("Washer >>")
	// fmt.Println(availableWasher)
	// fmt.Println("Dryer >>")
	// fmt.Println(availableDryer)

	//Price Calculation
	var calculatedPrice float64 = 0.0

	washingUnitPrice := serVicePriceMapper(int(basketWeight))
	calculatedPrice += float64(washingBasketCount * washingUnitPrice)
	dryingUnitPrice := serVicePriceMapper(int(dryingWeight))
	calculatedPrice += float64(dryingUnitPrice * dryinBasketCount)

	if !newOrder.ZuckOnsite {
		calculatedPrice += float64(model.DeliveryPrice + model.PickupPrice)
	}
	if isAgentsExist {
		calculatedPrice += float64(model.AgentsPrice)
	}

	payment := model.Payments{Amount: calculatedPrice}
	//Create new payment
	paymentResponse, err := u.paymentUsecase.CreatePayment(payment)
	if err != nil {
		return nil, errors.New("ERR: cannont create payment")
	}

	orderHeader := model.OrderHeader{
		OrderHeaderID:   uuid.New().String(),
		UserID:          newOrder.UserID,
		BranchID:        newOrder.BranchID,
		OrderNote:       newOrder.OrderNote,
		PaymentID:       paymentResponse.PaymentID, // temp solution
		ZuckOnsite:      newOrder.ZuckOnsite,
		DeliveryAddress: newOrder.DeliveryAddress,
		DeliveryLat:     newOrder.DeliveryLat,
		DeliveryLong:    newOrder.DeliveryLong,
		StarRating:      nil,
		ReviewComment:   nil,
		CreatedBy:       newOrder.UserID,
		UpdatedBy:       newOrder.UserID,
	}

	header, err := u.orderHeaderRepo.CreateOrderHeader(&orderHeader)

	if err != nil {
		return nil, err
	}

	var orderDetails []model.OrderDetail
	var washerIndexer int = 0
	var dryerIndexer int = 0
	for _, detail := range newOrder.OrderDetails {
		d := model.OrderDetail{
			OrderBasketID: uuid.New().String(),
			OrderHeaderID: orderHeader.OrderHeaderID,
			Weight:        detail.Weight,
			OrderStatus:   model.Waiting,
			ServiceType:   detail.ServiceType,
			FinishedAt:    nil,
			CreatedBy:     &newOrder.UserID,
			UpdatedBy:     &newOrder.UserID,
		}
		if d.ServiceType == "Washing" {
			d.MachineSerial = (&(*availableWasher)[washerIndexer].MachineSerial)
			washerIndexer += 1
		} else if d.ServiceType == "Drying" {
			d.MachineSerial = (&(*availableDryer)[dryerIndexer].MachineSerial)
			dryerIndexer += 1
		} else {
			d.MachineSerial = (&(*availableWasher)[0].MachineSerial)
		}
		orderDetails = append(orderDetails, d)
	}

	details, err := u.orderDetailRepo.CreateOrderDetails(&orderDetails)

	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindUserByUserID(newOrder.UserID)
	if err != nil {
		return nil, err
	}
	res := combineFullOrder(header, details, *user, false)

	return res, nil
}

func (u *orderUsecase) GetAllHeader() (*[]model.OrderHeader, error) {
	var headers *[]model.OrderHeader

	headers, err := u.orderHeaderRepo.GetAll()
	if err != nil {
		return &[]model.OrderHeader{}, err
	}

	for index, header := range *headers {
		user, err := u.userRepo.FindUserByUserID(header.UserID)
		if err != nil {
			(*headers)[index].UserDetail = model.UserDetailDTO{}
			continue
		}

		(*headers)[index].UserDetail = *toUserDetailDTO(user)
	}
	return headers, err
}

func (u *orderUsecase) GetByHeaderID(orderHeaderID string, isAdminView bool, option string) (interface{}, error) {
	headers, err := u.orderHeaderRepo.GetByID(orderHeaderID, isAdminView)
	if err != nil {
		return nil, err
	}

	if headers.OrderHeaderID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	if option == "header" {
		return headers, err
	}

	detail, err := u.orderDetailRepo.GetByHeaderID(orderHeaderID, isAdminView)

	if err != nil {
		return nil, err
	}

	if option == "detail" {
		return detail, err
	}
	user, err := u.userRepo.FindUserByUserID(headers.UserID)
	if err != nil {
		return nil, err
	}
	fullOrder := combineFullOrder(headers, detail, *user, isAdminView)

	return fullOrder, err
}

func (u *orderUsecase) GetByBranchID(branchID string, managerUserID string) ([]interface{}, error) {

	manager, err := u.userRepo.FindUserByUserID(managerUserID)
	if err != nil {
		return []interface{}{}, err
	}
	if manager.Role != "SuperAdmin" && manager.UserID != managerUserID {
		return []interface{}{}, errors.New("ERR: forbidden manager try to access unautherized branch")
	}

	headers, err := u.orderHeaderRepo.GetByBranchID(branchID)
	if err != nil {
		return []interface{}{}, err
	}

	detail, err := u.orderDetailRepo.GetAll()
	if err != nil {
		return []interface{}{}, err
	}

	if len(*headers) == 0 {
		return []interface{}{}, nil
	}

	var users []model.Users
	for _, header := range *headers {
		user, err := u.userRepo.FindUserByUserID(header.UserID)
		users = append(users, *user)
		if err != nil {
			return nil, errors.New("ERR: error occured when trying to query user data")
		}
	}

	var fullOrder []interface{}
	for i, h := range *headers {
		var thisDetail []model.OrderDetail
		for _, d := range *detail {
			if d.OrderHeaderID == h.OrderHeaderID {
				thisDetail = append(thisDetail, d)
			}
		}
		fullOrder = append(fullOrder, combineFullOrder(&h, &thisDetail, users[i], true))
	}

	return fullOrder, err
}

func (u *orderUsecase) GetByUserID(userID string) ([]interface{}, error) {
	headers, err := u.orderHeaderRepo.GetByUserID(userID)

	if err != nil {
		return []interface{}{}, err
	}

	detail, err := u.orderDetailRepo.GetAll()
	if err != nil {
		return []interface{}{}, err
	}

	if len(*headers) == 0 {
		return []interface{}{}, nil
	}

	var users []model.Users
	for _, header := range *headers {
		user, err := u.userRepo.FindUserByUserID(header.UserID)
		users = append(users, *user)
		if err != nil {
			return nil, errors.New("ERR: error occured when trying to query user data")
		}
	}

	var fullOrder []interface{}
	for i, h := range *headers {
		var thisDetail []model.OrderDetail
		for _, d := range *detail {
			if d.OrderHeaderID == h.OrderHeaderID {
				thisDetail = append(thisDetail, d)
			}
		}
		fullOrder = append(fullOrder, combineFullOrder(&h, &thisDetail, users[i], true))
	}

	return fullOrder, err
}

func (u *orderUsecase) UpdateStatus(order model.UpdateOrder) (interface{}, error) {

	updatedOrder := model.OrderDetail{
		OrderBasketID: order.OrderBasketID,
		MachineSerial: order.MachineSerial,
		OrderStatus:   order.OrderStatus,
		FinishedAt:    order.FinishedAt,
		UpdatedBy:     &order.UpdatedBy,
	}

	orderDetail, err := u.orderDetailRepo.UpdateStatus(updatedOrder)

	if err != nil {
		return nil, err
	}

	fullOrder, err := u.GetByHeaderID(orderDetail.OrderHeaderID, true, "full")

	if err != nil {
		return nil, err
	}

	return &fullOrder, err
}

func (u *orderUsecase) UpdateReview(review model.OrderReview) (*model.FullOrder, error) {
	orderModel := model.OrderHeader{
		OrderHeaderID: review.OrderHeaderID,
		UserID:        review.UserID,
		StarRating:    review.StarRating,
		ReviewComment: review.ReviewComment,
		UpdatedBy:     review.UserID,
	}

	herdead, err := u.orderHeaderRepo.GetByID(orderModel.OrderHeaderID, false)
	if err != nil {
		return nil, err
	}

	if herdead.UserID != orderModel.UserID {
		return nil, errors.New("err: forbidden review update")
	}

	orderHeader, err := u.orderHeaderRepo.UpdateReview(orderModel)
	if err != nil {
		return nil, err
	}

	orderDetails, err := u.orderDetailRepo.GetByHeaderID(orderHeader.OrderHeaderID, false)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.FindUserByUserID(orderHeader.UserID)
	if err != nil {
		return nil, err
	}
	fullOrder := combineFullOrder(orderHeader, orderDetails, *user, false)

	return fullOrder, err
}

func (u *orderUsecase) SoftDelete(orderHeaderID string, deletedBy string) (*model.FullOrder, error) {
	orderHeader, err := u.orderHeaderRepo.SoftDelete(orderHeaderID, deletedBy)

	if err != nil {
		return nil, err
	}

	orderDetails, err := u.orderDetailRepo.DeleteByHeaderID(orderHeaderID, deletedBy)

	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.FindUserByUserID(orderHeader.UserID)
	if err != nil {
		return nil, err
	}
	fullOrder := combineFullOrder(orderHeader, orderDetails, *user, true)

	return fullOrder, err
}
