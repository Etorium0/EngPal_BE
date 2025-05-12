package repo_impl

import "log"

type AssignmentRepoImpl struct{}

func (r *AssignmentRepoImpl) GenerateAssignment() {
	log.Println("Generating assignment")
}