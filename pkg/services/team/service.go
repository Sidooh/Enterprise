package team

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/entities"
)

type Service interface {
	FetchTeams() (*[]entities.Team, error)
	GetTeam(id int) (*entities.Team, error)
	CreateTeam(team *entities.Team) (*presenter.Team, error)

	FetchTeamsForEnterprise(enterpriseId int) (*[]entities.Team, error)
	GetTeamForEnterprise(enterpriseId int, id int) (*entities.Team, error)
}

type service struct {
	teamRepository Repository
}

func (s *service) FetchTeams() (*[]entities.Team, error) {
	return s.teamRepository.ReadTeams()
}

func (s *service) GetTeam(id int) (*entities.Team, error) {
	return s.teamRepository.ReadTeam(id)
}

func (s *service) CreateTeam(team *entities.Team) (*presenter.Team, error) {
	model, err := s.teamRepository.CreateTeam(team)

	return &presenter.Team{
		Id:           model.Id,
		Name:         model.Name,
		EnterpriseId: int(model.EnterpriseId),
	}, err
}

func (s *service) FetchTeamsForEnterprise(enterpriseId int) (*[]entities.Team, error) {
	return s.teamRepository.ReadEnterpriseTeams(enterpriseId)
}

func (s *service) GetTeamForEnterprise(enterpriseId int, id int) (*entities.Team, error) {
	return s.teamRepository.ReadEnterpriseTeam(enterpriseId, id)
}

func NewService(team Repository) Service {
	return &service{teamRepository: team}
}
