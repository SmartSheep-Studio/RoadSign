package sign

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var App *RoadApp

func ReadInConfig(root string) error {
	instance := &RoadApp{
		Sites: []*SiteConfig{},
	}

	if err := filepath.Walk(root, func(fp string, info os.FileInfo, err error) error {
		var site SiteConfig
		if info.IsDir() {
			return nil
		} else if file, err := os.OpenFile(fp, os.O_RDONLY, 0755); err != nil {
			return err
		} else if data, err := io.ReadAll(file); err != nil {
			return err
		} else if err := yaml.Unmarshal(data, &site); err != nil {
			return err
		} else {
			defer file.Close()

			// Extract file name as site id
			site.ID = strings.SplitN(filepath.Base(fp), ".", 2)[0]
			instance.Sites = append(instance.Sites, &site)
		}

		return nil
	}); err != nil {
		return err
	}

	App = instance

	return nil
}
