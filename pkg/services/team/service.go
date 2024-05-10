package team

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/account"
)

type Service interface {
	FetchTeams() (*[]entities.Team, error)
	GetTeam(id int) (*entities.Team, error)
	CreateTeam(team *entities.Team) (*presenter.Team, error)
	AddTeamAccount(team *entities.Team, accountId int) (*entities.Account, error)

	FetchTeamsForEnterprise(enterpriseId int) (*[]entities.Team, error)
	GetTeamForEnterprise(enterpriseId int, id int) (*entities.Team, error)
}

type service struct {
	teamRepository    Repository
	accountRepository account.Repository
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

func (s *service) AddTeamAccount(team *entities.Team, accountId int) (*entities.Account, error) {
	account, err := s.accountRepository.ReadAccount(accountId)
	if err != nil {
		return nil, pkg.ErrInvalidAccount
	}

	err = s.teamRepository.AddTeamAccount(team, account.Id)

	return account, err
}

func (s *service) FetchTeamsForEnterprise(enterpriseId int) (*[]entities.Team, error) {
	return s.teamRepository.ReadEnterpriseTeams(enterpriseId)
}

func (s *service) GetTeamForEnterprise(enterpriseId int, id int) (*entities.Team, error) {
	return s.teamRepository.ReadEnterpriseTeam(enterpriseId, id)
}

func NewService(team Repository) Service {
	accountRepository := account.NewRepo()
	return &service{teamRepository: team, accountRepository: accountRepository}
}
