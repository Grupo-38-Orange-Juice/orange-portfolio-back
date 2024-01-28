package usecases

import "github.com/Grupo-38-Orange-Juice/orange-portfolio-back/domain/entities"

type ProjectUseCase struct {
	projectRepository ProjectRepository
	UserRepository    UserRepository
}

func NewProjectUseCase(projectRepository ProjectRepository, UserRepository UserRepository) ProjectUseCase {
	return ProjectUseCase{
		projectRepository: projectRepository,
		UserRepository:    UserRepository,
	}
}

func (p ProjectUseCase) CreateProject(project *entities.Project, userId string) error {
	user, err := p.UserRepository.FindUserById(userId)
	if user == nil {
		return entityNotFound("user", "id", map[string]string{userId: userId})
	}
	if err != nil {
		return err
	}

	founded, err := p.projectRepository.FindProjectByNameAndUserId(project.Name, userId)
	if founded != nil {
		possibleValues := map[string]string{project.Name: project.Name}
		return entityAlreadyExists("project", "name", possibleValues)
	}
	if err != nil {
		return err
	}

	err = p.projectRepository.CreateProject(project, userId)
	if err != nil {
		return err
	}
	return nil
}