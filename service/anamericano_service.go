package service

import (
	"context"
	"strconv"

	"github.com/sunrin-ana/anamericano-golang"
)

type AnAmericanoService struct {
	client *anamericano.Client
}

func NewAnAmericanoService() *AnAmericanoService {
	return &AnAmericanoService{
		client: anamericano.NewClient(&anamericano.ContextTokenAuth{}, nil),
	}
}

func (s *AnAmericanoService) Check(
	token string,
	userID int64,
	relation string,
	ns string,
	targetId string) (bool, error) {

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

func (s *AnAmericanoService) Write(
	token string,
	userID int64,
	relation string,
	ns string,
	targetId string) (*anamericano.Permission, error) {
	// TODO: 맞나?
	ctx := anamericano.WithToken(context.Background(), token)

	resp, err := s.client.WritePermission(ctx, &anamericano.PermissionWriteRequest{
		ObjectNamespace: ns,
		ObjectID:        targetId,
		Relation:        relation,
		SubjectType:     "user",
		SubjectID:       strconv.FormatInt(userID, 10),
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
