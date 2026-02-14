package service

import (
	"analog-be/entity"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

const SITE_PREFIX = "https://log.ana.st/sitemaps/"
const ARITCLE_URL_FORMAT = "https://log.ana.st/%s/logs/%s"

const RSS_FEED_PREFIX = "<?xml version=\"1.0\" encoding=\"UTF-8\" ?><rss version=\"2.0\"><channel><title>Analog</title><link>https://log.ana.st/</link><description>Latest articles from Analog</description><copyright>2026 Application and Architecture Club, Sunrin Internet High School</copyright><ttl>60</ttl>"
const RSS_FEED_SUFFIX = "</channel></rss>"

type FeedService struct {
	logService      *LogService
	rssFeed         string
	isUpdating      bool
	needToBeUpdated bool
}

type SitemapIndex struct {
	XMLName xml.Name         `xml:"sitemapindex"`
	Sitemap []SitemapElement `xml:"sitemap"`
}

type SitemapElement struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc     string   `xml:"loc"`
	Lastmod string   `xml:"lastmod"`
}

type SitemapFile struct {
	XMLName xml.Name     `xml:"urlset"`
	Urls    []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func NewFeedService(logService *LogService) *FeedService {
	fs := &FeedService{logService: logService, rssFeed: "", isUpdating: false, needToBeUpdated: false}
	go func() {
		err := fs.UpdateFeed(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()
	return fs
}

func (f *FeedService) UpdateFeed(ctx context.Context) error {
	if f.isUpdating {
		f.needToBeUpdated = true
		return nil
	}

	f.needToBeUpdated = false
	f.isUpdating = true

	err := f.UpdateRSSFeed(ctx)
	if err != nil {
		return err
	}
	err = f.UpdateSitemap(nil)
	if err != nil {
		return err
	}

	f.isUpdating = false
	if f.needToBeUpdated {
		err = f.UpdateFeed(ctx)
	}

	return nil
}

func (f *FeedService) UpdateRSSFeed(ctx context.Context) error {
	rssFeed, err := f.GenerateRSSFeed(ctx)
	if err != nil {
		return err
	}

	f.rssFeed = rssFeed
	return nil
}

func (f *FeedService) GetRSSFeed() string {
	return f.rssFeed
}

func (f *FeedService) GenerateRSSFeed(ctx context.Context) (string, error) {
	list, err := f.logService.GetList(ctx, 20, 0)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(RSS_FEED_PREFIX)
	sb.WriteString(fmt.Sprintf("<lastBuildDate>%s</lastBuildDate>", list.Items[0].CreatedAt.Format(time.RFC1123Z)))

	for _, log := range list.Items {
		sb.WriteString(fmt.Sprintf("<item><title>%s</title><description>%s</description><guid isPermaLink=\"false\">%X</guid><link>%s</link><pubDate>%s</pubDate></item>", log.Title, log.Description, log.ID, BuildLogURL(log), log.CreatedAt.Format(time.RFC1123Z)))
	}

	sb.WriteString(RSS_FEED_SUFFIX)

	return sb.String(), nil
}

func (f *FeedService) UpdateSitemap(log *entity.Log) error {
	info, err := os.Stat("./sitemap")
	if (err != nil && errors.Is(err, os.ErrNotExist)) || !info.IsDir() {
		err = os.Mkdir("./sitemap", 0644)

		if err != nil {
			panic(err)
		}

		err := UpdateIndexMap(0)
		if err != nil {
			return err
		}

		return nil
	}

	// 아티클 생성으로 인한 호출이 아니면 건너뜀
	if log == nil {
		return nil
	}

	idx, err := GetLatestSitemap()
	if err != nil {
		return err
	}

	if ok, err := IsSitemapWritable(idx); err == nil || !ok {
		idx++
		err := UpdateIndexMap(0)
		if err != nil {
			return err
		}
	}

	// 사이트맵 저장
	err = WriteSitemap(idx, []SitemapURL{{Loc: BuildLogURL(log)}})
	if err != nil {
		return err
	}

	return nil
}

func (f *FeedService) GetSitemap(name string) string {
	file, err := os.ReadFile("./sitemap/" + name) // TODO: 해당 로직은 매우 위험함. 향후 Spine에서 공식적으로 static resource를 지원하게 되면 해당 로직을 대체할 것
	if err != nil {
		return ""
	}

	return string(file)
}

func UpdateIndexMap(idx int) error {
	file, err := os.Create("./sitemap/sitemap-index.xml")
	if err != nil {
		return err
	}

	enc := xml.NewEncoder(file)

	for i := 0; i < idx; i++ {
		enc.Encode(SitemapElement{Loc: fmt.Sprintf(SITE_PREFIX+"sitemap-%d.xml", i), Lastmod: time.Now().Format(time.RFC1123Z)})
	}
	enc.Flush()

	file.Close()
	return nil
}

func GetLatestSitemap() (int, error) {
	xmlFile, err := os.Open("./sitemap/sitemap-index.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return -1, err
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	var inElement string

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch ele := t.(type) {
		case xml.StartElement:
			inElement = ele.Name.Local
			if inElement == "sitemap" {
				var se SitemapElement
				decoder.DecodeElement(&se, &ele)

				var idx int
				_, err := fmt.Sscanf(se.Loc, SITE_PREFIX+"sitemap-%d.xml", &idx)
				if err != nil {
					fmt.Println("Error parsing sitemap index:", err)
					return 0, err
				}

				return idx, nil
			}
		default:
		}
	}

	return 0, nil
}

func IsSitemapWritable(idx int) (bool, error) {
	// 사이즈 체크
	if stat, serr := os.Stat(fmt.Sprintf("./sitemap/sitemap-%d.xml", idx)); serr == nil && stat.Size() > 49_000_000 {
		return false, nil
	}

	// URL 개수 체크
	xmlFile, err := os.Open(fmt.Sprintf("./sitemap/sitemap-%d.xml", idx))

	if err != nil {
		fmt.Println("Error opening file:", err)
		return false, err
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	var count int

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch ele := t.(type) {
		case xml.StartElement:
			if ele.Name.Local == "url" {
				count++
			}
		default:
		}
	}

	if count >= 50_000 {
		return false, nil
	}
	return true, nil
}

func WriteSitemap(idx int, urls []SitemapURL) error {
	f, err := os.Open(fmt.Sprintf("./sitemap/sitemap-%d.xml", idx))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		f = os.NewFile(3, fmt.Sprintf("./sitemap/sitemap-%d.xml", idx))
		f.WriteString(xml.Header)
		f.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n</urlset>")
	}

	defer f.Close()

	offset := int64(len("</urlset>"))
	f.Seek(-offset, io.SeekEnd)

	enc := xml.NewEncoder(f)
	for _, url := range urls {
		err := enc.Encode(url)
		if err != nil {
			return err
		}
	}
	enc.Flush()

	f.WriteString("</urlset>")
	return nil
}

func BuildLogURL(log *entity.Log) string {
	return fmt.Sprintf(
		ARITCLE_URL_FORMAT,
		log.LoggedBy[0].Handle,
		fmt.Sprintf("%s-%X", url.PathEscape(strings.ReplaceAll(log.Title, " ", "-")), log.ID), // Help Me (ID: 1) -> Help-Me-1; (제목)-(아이디 HEX)
	)
}
