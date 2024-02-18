package RpcServer

import (
	"github.com/google/uuid"
)

type ProjectRPCPayload struct {
	Id       int
	UserUuid uuid.UUID
	Name     string
	Token    string
}

func NewProjectRPCPayload(id int, UserUuid uuid.UUID, name string, token string) *ProjectRPCPayload {
	return &ProjectRPCPayload{
		Id:       id,
		UserUuid: UserUuid,
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
	resp = NewProjectRPCPayload(project.Id, project.UserUuid, project.Name, project.Token)
	return nil
}

func (s Server) GetProjectList(userUuid uuid.UUID, resp *[]ProjectRPCPayload) error {
	projectList, err := s.projectUsecase.GetProjectListByUserUuid(userUuid, uuid.Nil)
	if err != nil {
		return err
	}
	for _, project := range projectList {
		payload := NewProjectRPCPayload(project.Id, project.UserUuid, project.Name, project.Token)
		*resp = append(*resp, *payload)
	}
	return nil
}
