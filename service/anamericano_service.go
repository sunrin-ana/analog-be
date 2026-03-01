package service

import (
	"context"
	"os"
	"strconv"

	"github.com/sunrin-ana/anamericano-golang"
)

type AnAmericanoService interface {
	Check(
		userID int64,
		relation string,
		ns string,
		targetId string) (bool, error)
	Write(
		userID int64,
		relation string,
		ns string,
		targetId int64) (*anamericano.Permission, error)
}

type AnAmericanoServiceImpl struct {
	client *anamericano.Client
}

func NewAnAmericanoService() AnAmericanoService {
	return &AnAmericanoServiceImpl{
		client: anamericano.NewClient(&anamericano.ContextTokenAuth{}, nil),
	}
}

func (s *AnAmericanoServiceImpl) Check(
	userID int64,
	relation string,
	ns string,
	targetId string) (bool, error) {

	token := os.Getenv("AN_ACCOUNT_API_TOKEN")

	ctx := anamericano.WithToken(context.Background(), token)

	resp, err := s.client.CheckPermission(ctx, &anamericano.PermissionCheckRequest{
		SubjectType:     "user",
		SubjectID:       strconv.FormatInt(userID, 10),
		Relation:        relation,
		ObjectNamespace: ns,
		ObjectID:        targetId,
	})

	if err != nil {
		return false, err
	}

	return resp.Allowed, nil
}

func (s *AnAmericanoServiceImpl) Write(
	userID int64,
	relation string,
	ns string,
	targetId int64) (*anamericano.Permission, error) {
	// TODO: 맞나?
	token := os.Getenv("AN_ACCOUNT_API_TOKEN")

	ctx := anamericano.WithToken(context.Background(), token)

	resp, err := s.client.WritePermission(ctx, &anamericano.PermissionWriteRequest{
		ObjectNamespace: ns,
		ObjectID:        strconv.FormatInt(targetId, 10),
		Relation:        relation,
		SubjectType:     "user",
		SubjectID:       strconv.FormatInt(userID, 10),
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
