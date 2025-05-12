package repo_impl

import "log"

type UserRepoImpl struct{}

func (r *UserRepoImpl) GetUser() {
    log.Println("Fetching user data")
}