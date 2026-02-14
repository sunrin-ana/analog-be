package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"bytes"
	"context"
	"time"

	"github.com/huantt/plaintext-extractor"

	"github.com/yuin/goldmark"
	"golang.org/x/sync/semaphore"
)

type LogService struct {
	logRepository      *repository.LogRepository
	commentRepository  *repository.CommentRepository
	anamericanoService *AnAmericanoService
	feedService        *FeedService
	plainExtractor     *plaintext.Extractor
	prerenderJobs      *semaphore.Weighted
}

func NewLogService(logRepository *repository.LogRepository, commentRepository *repository.CommentRepository, anamericanoService *AnAmericanoService) *LogService {
	return &LogService{
		logRepository:      logRepository,
		commentRepository:  commentRepository,
		anamericanoService: anamericanoService,
		plainExtractor:     plaintext.NewMarkdownExtractor(),
		prerenderJobs:      semaphore.NewWeighted(4),
	}
}

func (s *LogService) Get(ctx context.Context, id *entity.ID) (*entity.Log, error) {
	log, err := s.logRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (s *LogService) GetList(ctx context.Context, limit int, offset int) (*dto.PaginatedResult[*entity.Log], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	logs, total, err := s.logRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedResult[*entity.Log]{
		Items:  logs,
		Total:  *total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *LogService) Search(ctx context.Context, query string, limit int, offset int) (*dto.PaginatedResult[*entity.Log], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	logs, total, err := s.logRepository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedResult[*entity.Log]{
		Items:  logs,
		Total:  *total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *LogService) Create(ctx context.Context, req *dto.LogCreateRequest, authorID *entity.ID) (*entity.Log, error) {
	now := time.Now().UTC()

	log := &entity.Log{
		Title:       req.Title,
		Description: s.BuildDescription(req.Content),
		Generations: req.Generations,
		Content:     req.Content,
		CreatedAt:   now,
	}

	authorIDs := make([]entity.ID, len(req.CoAuthorIDs)+1)
	authorIDs = append(authorIDs, *authorID)
	for _, c := range req.CoAuthorIDs {
		authorIDs = append(authorIDs, c)
	}

	log, err := s.logRepository.Create(ctx, log, &req.TopicIDs, &req.CoAuthorIDs)
	if err != nil {
		return nil, err
	}

	_, err = s.anamericanoService.Write(*authorID, "owner", "analog_log", log.ID)
	if err != nil {
		return nil, err
	}

	for _, id := range authorIDs[1:] {
		_, err = s.anamericanoService.Write(id, "editor", "analog_log", log.ID)
		if err != nil {
			return nil, err
		}
	}

	s.feedService.UpdateFeed()

	return log, nil
}

func (s *LogService) Update(ctx context.Context, id *entity.ID, req *dto.LogUpdateRequest, authorID *entity.ID) (*entity.Log, error) {
	log, err := s.logRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		log.Title = *req.Title
	}
	if req.Generations != nil {
		log.Generations = *req.Generations
	}

	if req.Content != nil {
		log.Content = *req.Content
		log.Description = s.BuildDescription(*req.Content)

		go func() {
			gctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := s.prerenderJobs.Acquire(gctx, 1); err != nil {
				return
			}

			defer s.prerenderJobs.Release(1)

			err := s.PreRender(gctx, id)
			if err != nil {
				println(err.Error())
			}
		}()
	}

	if req.CoAuthorIDs != nil {
		authorIDs := make([]entity.ID, len(*req.CoAuthorIDs)+1)
		authorIDs = append(authorIDs, *authorID)
		for _, c := range *req.CoAuthorIDs {
			authorIDs = append(authorIDs, c)
		}
		log, err = s.logRepository.Update(ctx, log, req.TopicIDs, &authorIDs)
	} else {
		log, err = s.logRepository.Update(ctx, log, req.TopicIDs, nil)
	}

	if err != nil {
		return nil, err
	}

	return log, nil
}

func (s *LogService) Delete(ctx context.Context, id *entity.ID) error {
	err := s.commentRepository.DeleteByLogID(ctx, id)
	if err != nil {
		return err
	}

	return s.logRepository.Delete(ctx, id)
}

func (s *LogService) BuildDescription(content string) string {
	pcontent, err := s.plainExtractor.PlainText(content)
	if err != nil {
		pcontent = &content
	}

	pcontentRune := []rune(*pcontent)
	var description string
	if len(pcontentRune) > 100 {
		description = string(pcontentRune[:97]) + "..."
	} else {
		description = *pcontent
	}

	return description
}

func (s *LogService) PreRender(ctx context.Context, id *entity.ID) error {
	log, err := s.logRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	var rendered bytes.Buffer

	if err = goldmark.Convert([]byte(log.Content), &rendered); err != nil {
		return err
	}

	log.PreRendered = rendered.String()

	log, err = s.logRepository.Update(ctx, log, nil, nil)

	if err != nil {
		return err
	}

	return nil
}
