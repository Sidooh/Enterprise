package enterprise

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"sidooh-enterprise-gateway/pkg/entities"
	"strconv"
)

type Service interface {
	GetEnterprise(id int) (*entities.Enterprise, error)
}

type service struct {
	apiClient *fiber.Client
}

func (s *service) GetEnterprise(id int) (*entities.Enterprise, error) {

	var e *entities.Enterprise
	code, body, errors := s.apiClient.Get("http://localhost:8001/" + strconv.Itoa(id)).Struct(e)
	if len(errors) > 0 {
		log.Print(errors)
		return nil, errors[0]
	}

	fmt.Println(code, body)

	return e, errors[0]
}

func NewService() Service {
	return &service{apiClient: &fiber.Client{}}
}
