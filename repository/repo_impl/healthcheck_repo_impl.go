package repo_impl

import "log"

type HealthcheckRepoImpl struct{}

func (r *HealthcheckRepoImpl) CheckHealth() {
	log.Println("Performing health check")
}