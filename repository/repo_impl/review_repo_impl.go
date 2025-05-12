package repo_impl

import "log"

type ReviewRepoImpl struct{}

func (r *ReviewRepoImpl) GenerateReview() {
	log.Println("Generating review")
}