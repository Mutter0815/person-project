package dto

type UpdatePersonRequest struct {
	Name        *string `json:"name" example:"Ivan"`
	Surname     *string `json:"surname" example:"Ivanov"`
	Patronymic  *string `json:"patronymic" example:"Ivanovich"`
	Age         *int    `json:"age" example:"30"`
	Gender      *string `json:"gender" example:"male" validate:"omitempty,oneof=male female"`
	Nationality *string `json:"nationality" example:"Russian"`
}

type CreatePersonRequest struct {
	Name       string  `json:"name" example:"Ivan"`
	Surname    string  `json:"surname" example:"Ivanov"`
	Patronymic *string `json:"patronymic" example:"Ivanovich"`
}
