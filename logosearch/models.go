package logosearch

import (
	"fmt"
	"time"

	"github.com/AwesomeLogos/bimi-explorer/generated"
)

type Image struct {
	Img  string "json:\"img\""
	Name string "json:\"name\""
	Src  string "json:\"src\""
}

type Source struct {
	Handle       string  "json:\"handle\""
	Images       []Image "json:\"images\""
	LastModified string  "json:\"lastmodified\""
	Logo         string  "json:\"logo\""
	Name         string  "json:\"name\""
	Provider     string  "json:\"provider\""
	ProviderIcon string  "json:\"provider_icon\""
	Url          string  "json:\"url\""
	Website      string  "json:\"website\""
}

func GenerateIndex(domains []generated.Domain) Source {

	var source Source
	source.Logo = "https://www.vectorlogo.zone/logos/bimigroup/bimigroup-icon.svg"
	source.Handle = "bimi"
	source.Images = []Image{}
	source.Name = "BIMI"
	source.Provider = "remote"
	source.ProviderIcon = "https://logosear.ch/images/remote.svg"
	source.Url = "https://bimi-explorer.svg.zone/"
	source.Website = "https://bimigroup.org/"

	lastModified := time.Time{}

	for _, domain := range domains {
		image := Image{
			Img:  domain.Imgurl.String,
			Name: domain.Domain, //LATER: strip public suffix
			Src:  fmt.Sprintf("https://bimi-explorer.svg.zone/bimi/%s/", domain.Domain),
		}
		if domain.Updated.Time.After(lastModified) {
			lastModified = domain.Updated.Time
		}
		source.Images = append(source.Images, image)
	}
	source.LastModified = lastModified.UTC().Format(time.RFC3339)

	return source
}
