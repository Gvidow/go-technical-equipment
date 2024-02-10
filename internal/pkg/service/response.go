package service

import (
	"time"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

type ResponseOk struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"success"`
	Body    any    `json:",omitempty"`
}

type ResponseError struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"fail"`
}

type responseBodyID struct {
	Id int `json:"id" example:"5"`
}

type responseWithEquipment struct {
	ID          int    `json:"id" example:"16"`
	Title       string `json:"title" example:"EquipmentTitle"`
	Picture     string `json:"picture" example:"http://technical-equipment.ru/equipment/picture/like.png"`
	Description string `json:"description" example:"info about this equipment"`
	Status      string `json:"status" example:"active"`
}

type responseWithFeedEquipment struct {
	Equipments []responseWithEquipment `json:"equipments"`
	DraftId    int                     `json:"last_request_id"`
}

type requestDTO struct {
	ID               int       `json:"id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	FormatedAt       time.Time `json:"formated_at"`
	CompletedAt      time.Time `json:"completed_at"`
	CreatorProfile   *ds.User  `json:"creator_profile"`
	ModeratorProfile *ds.User  `json:"moderator_profile"`
}

type responseRequests []requestDTO
