package dto

import "encoding/xml"

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
