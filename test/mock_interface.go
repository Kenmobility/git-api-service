package test

import "github.com/kenmobility/git-api-service/internal/repository"

type Store interface {
	repository.CommitRepository
	repository.RepoMetadataRepository
}
