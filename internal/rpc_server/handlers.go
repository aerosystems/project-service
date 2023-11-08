package RPCServer

type ProjectRPCPayload struct {
	ID     int
	UserId int
	Name   string
	Token  string
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
