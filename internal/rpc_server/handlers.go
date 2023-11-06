package RPCServer

type CreateProjectRPCPayload struct {
	UserId int
	Name   string
}

type ProjectRPCPayload struct {
	ID     int
	UserId int
	Name   string
	Token  string
}

func (ps *ProjectServer) CreateProject(payload CreateProjectRPCPayload, resp *string) error {
	//projectList, err := ps.projectRepo.GetByUserId(payload.UserId)
	//if err != nil {
	//	return err
	//}
	//
	//for _, project := range projectList {
	//	if project.Name == payload.Name {
	//		err := fmt.Errorf("project with Name %s already exists", payload.Name)
	//		return err
	//	}
	//}
	//
	//var newProject = models.Project{
	//	UserId: payload.UserId,
	//	Name:   payload.Name,
	//}
	//
	//if err = ps.projectRepo.Create(&newProject); err != nil {
	//	return err
	//}
	//
	//*resp = fmt.Sprintf("project %s successfully created", payload.Name)
	return nil
}

func (ps *ProjectServer) GetProject(projectToken string, resp *ProjectRPCPayload) error {
	//project, err := ps.projectRepo.GetByToken(projectToken)
	//if err != nil {
	//	return err
	//}
	//
	//*resp = ProjectRPCPayload{
	//	ID:     project.ID,
	//	UserId: project.UserId,
	//	Name:   project.Name,
	//	Token:  project.Token,
	//}
	return nil
}

func (ps *ProjectServer) GetProjectList(userID int, resp *[]ProjectRPCPayload) error {
	//projectList, err := ps.projectRepo.GetByUserId(userID)
	//if err != nil {
	//	return err
	//}
	//
	//for _, project := range projectList {
	//	*resp = append(*resp, ProjectRPCPayload{
	//		ID:     project.ID,
	//		UserId: project.UserId,
	//		Name:   project.Name,
	//		Token:  project.Token,
	//	})
	//}
	return nil
}
