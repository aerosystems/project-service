package RPCServer

type ProjectRPCPayload struct {
	Id     int
	UserId int
	Name   string
	Token  string
}

func NewProjectRPCPayload(id int, userId int, name string, token string) *ProjectRPCPayload {
	return &ProjectRPCPayload{
		Id:     id,
		UserId: userId,
		Name:   name,
		Token:  token,
	}
}

func (ps *ProjectServer) CreateDefaultProject(projectPayload ProjectRPCPayload, resp *ProjectRPCPayload) error {
	if err := ps.projectService.DetermineStrategy(projectPayload.UserId, "customer"); err != nil {
		return err
	}
	if err := ps.projectService.CreateDefaultProject(projectPayload.UserId); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectServer) GetProject(projectToken string, resp *ProjectRPCPayload) error {
	project, err := ps.projectService.GetProjectByToken(projectToken)
	if err != nil {
		return err
	}
	resp = NewProjectRPCPayload(project.Id, project.UserId, project.Name, project.Token)
	return nil
}

func (ps *ProjectServer) GetProjectList(userId int, resp *[]ProjectRPCPayload) error {
	projectList, err := ps.projectService.GetProjectListByUserId(userId, 0)
	if err != nil {
		return err
	}
	for _, project := range projectList {
		payload := NewProjectRPCPayload(project.Id, project.UserId, project.Name, project.Token)
		*resp = append(*resp, *payload)
	}
	return nil
}
