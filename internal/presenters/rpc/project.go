package RpcServer

import (
	"github.com/google/uuid"
)

type ProjectRPCPayload struct {
	Uuid     string
	UserUuid uuid.UUID
	Name     string
	Token    string
}

func NewProjectRPCPayload(projectUuid string, userUuid uuid.UUID, name string, token string) *ProjectRPCPayload {
	return &ProjectRPCPayload{
		Uuid:     projectUuid,
		UserUuid: userUuid,
		Name:     name,
		Token:    token,
	}
}

func (s Server) CreateDefaultProject(projectPayload ProjectRPCPayload, resp *ProjectRPCPayload) error {
	if err := s.projectUsecase.DetermineStrategy(projectPayload.UserUuid.String(), "customer"); err != nil {
		return err
	}
	if err := s.projectUsecase.CreateDefaultProject(projectPayload.UserUuid); err != nil {
		return err
	}
	return nil
}

func (s Server) GetProject(projectToken string, resp *ProjectRPCPayload) error {
	project, err := s.projectUsecase.GetProjectByToken(projectToken)
	if err != nil {
		return err
	}
	resp.UserUuid = project.CustomerUUID
	resp.Name = project.Name
	resp.Token = project.Token
	return nil
}

func (s Server) GetProjectList(userUuid uuid.UUID, resp *[]ProjectRPCPayload) error {
	if err := s.projectUsecase.DetermineStrategy(userUuid.String(), "customer"); err != nil {
		return err
	}
	projectList, err := s.projectUsecase.GetProjectListByCustomerUuid(userUuid, uuid.Nil)
	if err != nil {
		return err
	}
	for _, project := range projectList {
		payload := NewProjectRPCPayload(project.Uuid.String(), project.CustomerUUID, project.Name, project.Token)
		*resp = append(*resp, *payload)
	}
	return nil
}
