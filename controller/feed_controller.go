package controller

import (
	"analog-be/service"
	"context"
	"net/http"
	"strings"

	"github.com/NARUBROWN/spine/pkg/httpx"
)

type FeedController struct {
	service *service.FeedService
}

func NewFeedController(service *service.FeedService) *FeedController {
	return &FeedController{service: service}
}

func (c *FeedController) GetFeed(ctx context.Context) httpx.Response[string] {
	return httpx.Response[string]{
		Body: c.service.GetRSSFeed(),
		Options: httpx.ResponseOptions{
			Headers: map[string]string{"Content-Type": "application/rss+xml"},
		},
	}
}

func (c *FeedController) GetSitemap(ctx context.Context, file string) httpx.Response[string] {
	if !strings.HasPrefix(file, "sitemap-") || !strings.HasSuffix(file, ".xml") {
		return httpx.Response[string]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest,
			},
		}
	}

	st := c.service.GetSitemap(file)
	if st == "" {
		return httpx.Response[string]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound,
			},
		}
	}

	return httpx.Response[string]{
		Body: st,
		Options: httpx.ResponseOptions{
			Headers: map[string]string{"Content-Type": "application/xml"},
		},
	}
}
