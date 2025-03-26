package dkRegistry

import (
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnSlice"
)

const (
	ErrNotFoundRepository = "not_found_repository"
	ErrNotFoundTag        = "not_found_tag"
)

func HasRepository(
	args Registry,
	repository string,
) (has bool, err error) {
	var repositories []string
	if repositories, err = QueryRepositories(args); err != nil {
		return
	}

	if has = fnSlice.Has(repositories, func(v string) bool {
		return v == repository
	}); !has {
		err = fnError.NewFields(ErrNotFoundRepository, map[string]any{
			"repository": repository,
		})
		return
	}

	return
}

func HasTag(
	args Registry,
	repository string,
	tag string,
) (err error) {
	var has bool
	if has, err = HasRepository(args, repository); err != nil {
		return
	}

	if !has {
		err = fnError.NewFields(ErrNotFoundRepository, map[string]any{
			"repository": repository,
		})
		return
	}

	var tags []string
	if tags, err = QueryTags(args, repository); err != nil {
		return
	}

	if has = fnSlice.Has(tags, func(v string) bool {
		return v == tag
	}); !has {
		err = fnError.NewFields(ErrNotFoundTag, map[string]any{
			"repository": repository,
			"tag":        tag,
		})
		return
	}

	return
}
