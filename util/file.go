package util

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/pathfinder-cm/pathfinder-agent/config"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

// TODO: to be abstracted
func GenerateBootstrapScriptContent(bs pfmodel.Bootstrapper) (string, int, error) {
	var tpl bytes.Buffer
	var mode int
	if bs.Type == "chef-solo" {
		const content = `
CHEF_FLAG_FILE=/tmp/chef_installed.txt
if [ ! -f "$CHEF_FLAG_FILE" ]; then
	echo "$CHEF_FLAG_FILE doesn't exist, creating file..."
	cd /tmp && curl -LO {{.BootstrapInstallerUrl}} && sudo bash ./install.sh -v {{.BootstrapVersion}} && rm install.sh && touch chef_installed.txt
fi

CHEF_REPO_DIR=/tmp/chef-repo-master
[ -d "$CHEF_REPO_DIR" ] && rm -rf $CHEF_REPO_DIR
mkdir $CHEF_REPO_DIR && wget {{.CookbooksURL}} -O - | tar -xz -C /tmp/chef-repo-master --strip-components=1

SOLO_FILE=/tmp/solo.rb
if [ ! -f "$SOLO_FILE" ]; then
	echo "$SOLO_FILE doesn't exist, creating file..."
	cat > solo.rb << EOF
cookbook_path "/tmp/chef-repo-master/cookbooks"
role_path "/tmp/chef-repo-master/roles"
EOF
fi

cat > /tmp/attributes.json << EOF
{{.BootstrapAttributes}}
EOF

chef-solo -c /tmp/solo.rb -j /tmp/attributes.json {{.BootstrapFlagOptions}}
`
		attributes, _ := json.Marshal(bs.Attributes)

		tmpl := template.Must(template.New("content").Parse(content))
		err := tmpl.Execute(&tpl, struct {
			BootstrapInstallerUrl string
			BootstrapVersion      string
			BootstrapAttributes   string
			BootstrapFlagOptions  string
			CookbooksURL          string
		}{
			BootstrapInstallerUrl: config.BootstrapInstallerUrl,
			BootstrapVersion:      config.BootstrapVersion,
			BootstrapAttributes:   string(attributes),
			BootstrapFlagOptions:  config.BootstrapFlagOptions,
			CookbooksURL:          bs.CookbooksUrl,
		})

		if err != nil {
			return "", 0, err
		}

		mode = 600
	}

	return tpl.String(), mode, nil
}
