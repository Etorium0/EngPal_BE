package repo_impl

import "log"

type GitHubRepoImpl struct{}

func (r *GitHubRepoImpl) FetchData() {
    log.Println("Fetching data from GitHub")
}