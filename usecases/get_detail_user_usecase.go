package usecases

import (
	"auth-service/domains"
	"auth-service/utils"
	"fmt"
	"time"
)

type getDetailUserUsecase struct {
	redisRepo domains.RedisRepository
	userRepo  domains.UserRepository
}

func NewGetDetailUserUsecase(redisRepo domains.RedisRepository, userRepo domains.UserRepository) domains.GetDetailUserUsecase {
	return &getDetailUserUsecase{
		redisRepo,
		userRepo,
	}
}

func (gudu *getDetailUserUsecase) GetDetailUserOnRedis(user_id string) (domains.DetailUserData, error) {
	var data domains.DetailUserData

	use_data, err := gudu.redisRepo.RedisUserHMGetAll(user_id)

	if err != nil {
		return data, err
	}

	if use_data["Email"] == "" {
		fmt.Println("user_id on Redis does not exists")
		return data, nil
	}

	data.UserID = user_id
	data.Email = use_data["Email"]
	data.FullName = use_data["FullName"]
	data.PhoneNumber = use_data["PhoneNumber"]
	data.BirthDay = use_data["BirthDay"]

	return data, nil
}

func (gudu *getDetailUserUsecase) GetDetailUserOnDatabase(user_id string) (domains.DetailUserData, error) {
	var data domains.DetailUserData
	var user_data_dto domains.UserDTO

	user_data, err := gudu.userRepo.GetByIDPreload(user_id, "LoginMethod", "Profiles")

	if err != nil {
		return data, err
	}

	// convert user_model sang DTO , DetailUserData
	data.UserID = user_id
	data.Email = user_data.Email
	data.CreatedAt = user_data.CreatedAt

	user_data_dto.UserID = fmt.Sprint(user_data.UserID)
	user_data_dto.Email = user_data.Email
	user_data_dto.UserStatus = string(user_data.UserStatus)
	user_data_dto.LoginMethod = user_data.LoginMethod.LoginMethod
	user_data_dto.LoginMethodID = user_data.LoginMethodID
	user_data_dto.CreatedAt = user_data.CreatedAt
	user_data_dto.UpdatedAt = user_data.UpdatedAt

	for _, profile := range user_data.Profiles {
		switch profile.ProfileKey {
		case "full_name":
			data.FullName = profile.ProfileValue
			user_data_dto.FullName = profile.ProfileValue
		case "birth_day":
			data.BirthDay = profile.ProfileValue
			user_data_dto.BirthDay = profile.ProfileValue
		case "phone_number":
			data.PhoneNumber = profile.ProfileValue
			user_data_dto.PhoneNumber = profile.ProfileValue
		}
	}

	// đông thời set lên lại Redis với data này
	userDataMapString := utils.StructureToMapString(user_data_dto)

	// lưu trữ thông tin user trên Redis với key là user_id
	timeToLiveUserData := time.Hour * 24 // Thời gian cache là 1 ngày
	gudu.redisRepo.RedisUserHMSet(user_id, userDataMapString, timeToLiveUserData)
	// lưu trữ thông tin user trên Redis với key là email để cache dữ liệu
	gudu.redisRepo.RedisUserHMSet(user_id, userDataMapString, timeToLiveUserData)

	return data, nil
}
