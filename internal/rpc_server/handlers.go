package RPCServer

import (
	"github.com/google/uuid"
	"log"
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

func (ps *ProjectServer) CreateDefaultProject(projectPayload ProjectRPCPayload, resp *ProjectRPCPayload) error {
	log.Println("CreateDefaultProject", projectPayload)
	if err := ps.projectService.DetermineStrategy(projectPayload.UserUuid.String(), "customer"); err != nil {
		return err
	}
	if err := ps.projectService.CreateDefaultProject(projectPayload.UserUuid); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectServer) GetProject(projectToken string, resp *ProjectRPCPayload) error {
	project, err := ps.projectService.GetProjectByToken(projectToken)
	if err != nil {
		return err
	}
	resp = NewProjectRPCPayload(project.Id, project.UserUuid, project.Name, project.Token)
	return nil
}

func (ps *ProjectServer) GetProjectList(userUuid uuid.UUID, resp *[]ProjectRPCPayload) error {
	projectList, err := ps.projectService.GetProjectListByUserUuid(userUuid, uuid.Nil)
	if err != nil {
		return err
	}
	for _, project := range projectList {
		payload := NewProjectRPCPayload(project.Id, project.UserUuid, project.Name, project.Token)
		*resp = append(*resp, *payload)
	}
	return nil
}
