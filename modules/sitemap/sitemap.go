package sitemap

import (
	"strconv"

	"github.com/valyala/bytebufferpool"
	"userstyles.world/models"
)

// CreateSitemap creates a sitemap.xml which is compatbile with the search engines.
// It should return the []byte of the sitemap.xml or an error if something went wrong.
// We should take a input of our database users and append them to a hardcoded list of paths.

var (
	siteMapXMLHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>`)
	siteMapURLSet    = []byte(`<urlset
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
	http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">`)

	siteMapHCEStart = []byte(`<url><loc>`)
	siteMapHCEEnd   = []byte(`</loc></url>`)

	siteMapURLEnd = []byte(`</urlset>`)
)

var hardCodedEntries = []string{
	"https://userstyles.world/",
	"https://userstyles.world/login",
	"https://userstyles.world/register",
	"https://userstyles.world/search",
	"https://userstyles.world/explore",
	"https://userstyles.world/docs/security",
	"https://userstyles.world/docs/faq",
	"https://userstyles.world/docs/code-of-conduct",
	"https://userstyles.world/docs/content-guidelines",
	"https://userstyles.world/docs/privacy",
}

var SiteMapCache []byte

func CreateSitemap(styles []models.StyleSiteMap) ([]byte, error) {
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)

	_, _ = buffer.Write(siteMapXMLHeader)
	_, _ = buffer.Write(siteMapURLSet)

	for _, hardCodedEntry := range hardCodedEntries {
		_, _ = buffer.Write(siteMapHCEStart)
		_, _ = buffer.WriteString(hardCodedEntry)
		_, _ = buffer.Write(siteMapHCEEnd)
	}

	for _, style := range styles {
		_, _ = buffer.Write(siteMapHCEStart)
		_, _ = buffer.WriteString("https://userstyles.world/style/")
		_, _ = buffer.WriteString(strconv.Itoa(style.ID))
		_, _ = buffer.Write(siteMapHCEEnd)
	}
	_, _ = buffer.Write(siteMapURLEnd)
	// Return the buffer copy
	return buffer.Bytes()[:], nil
}

func UpdateSitemapCache() error {
	styles, err := models.GetAllSitesSiteMap()
	if err != nil {
		return err
	}
	siteMap, err := CreateSitemap(styles)
	if err != nil {
		return err
	}
	SiteMapCache = siteMap
	return nil
}
