package sitemap

import (
	"strconv"

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

	siteURL = []byte(`https://userstyles.world/style/`)
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
	// Make a educated guess of the size of the sitemap.xml
	// Allocate 0.1MB
	buffer := make([]byte, 0, 1048576)
	buffer = append(buffer, siteMapXMLHeader...)
	buffer = append(buffer, siteMapURLSet...)

	for _, hardCodedEntry := range hardCodedEntries {
		buffer = append(buffer, siteMapHCEStart...)
		buffer = append(buffer, hardCodedEntry...)
		buffer = append(buffer, siteMapHCEEnd...)
	}

	for _, style := range styles {
		buffer = append(buffer, siteMapHCEStart...)
		buffer = append(buffer, siteURL...)
		buffer = append(buffer, strconv.Itoa(style.ID)...)
		buffer = append(buffer, siteMapHCEEnd...)
	}
	buffer = append(buffer, siteMapURLEnd...)
	// Return the filled buffer
	return buffer, nil
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
