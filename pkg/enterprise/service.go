package enterprise

import (
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/logger"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
	code, body, errors := s.apiClient.Get("http://localhost:8001/enterprises/" + strconv.Itoa(id)).Struct(e)
	if len(errors) > 0 {
		logger.ClientLog.Error(errors)
		return nil, errors[0]
	}

	fmt.Println(code, body)

	return e, errors[0]
}

func NewService() Service {
	return &service{apiClient: &fiber.Client{}}
}
