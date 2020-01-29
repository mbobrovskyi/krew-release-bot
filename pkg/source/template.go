package source

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/rajatjindal/krew-release-bot/pkg/krew"
	"github.com/sirupsen/logrus"
)

//ProcessTemplate process the .krew.yaml template for the release request
func ProcessTemplate(templateFile string, values interface{}) (string, []byte, error) {
	name := path.Base(templateFile)
	t := template.New(name).Funcs(map[string]interface{}{
		"addURIAndSha": func(url, tag string) string {
			t := struct {
				TagName string
			}{
				TagName: tag,
			}
			buf := new(bytes.Buffer)
			temp, err := template.New("url").Parse(url)
			if err != nil {
				panic(err)
			}

			err = temp.Execute(buf, t)
			if err != nil {
				panic(err)
			}

			logrus.Infof("getting sha256 for %s", buf.String())
			sha256, err := getSha256ForAsset(buf.String())
			if err != nil {
				panic(err)
			}

			return fmt.Sprintf(`uri: %s
    sha256: %s`, buf.String(), sha256)
		},
	})

	templateObject, err := t.ParseFiles(templateFile)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	err = templateObject.Execute(buf, values)
	if err != nil {
		return "", nil, err
	}

	pluginName, err := krew.GetPluginName(buf.Bytes())
	if err != nil {
		return "", nil, err
	}

	return pluginName, buf.Bytes(), nil
}
