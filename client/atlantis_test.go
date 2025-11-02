package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zkfmapf123/at-plan/usecase"
	"github.com/zkfmapf123/at-plan/utils"
)

func Test_ValidConfigFile(t *testing.T) {

	pwd, _ := utils.GetPwd()
	filePath := filepath.Join(pwd, "atlantis_config_file.yaml")

	yamlContent := `version: 3
projects:
  - name: ec2
    dir: examples/ec2
    workflow: terraform
  - name: sg
    dir: examples/sg
    workflow: terraform
workflows:
  terraform:
    plan:
      steps:
        - init:
            extra_args: ["-upgrade", "-reconfigure"]
        - env:
            name: TERRAGRUNT_TFPATH
            command: 'echo "terraform${ATLANTIS_TERRAFORM_VERSION}"'
        - env:
            name: TF_IN_AUTOMATION
            value: "true"
        - run:
            command: terraform plan -input=false -out=$PLANFILE
            output: strip_refreshing
    apply:
      steps:
        - env:
            name: TERRAGRUNT_TFPATH
            command: 'echo "terraform${ATLANTIS_TERRAFORM_VERSION}"'
        - env:
            name: TF_IN_AUTOMATION
            value: "true"
        - run: terraform apply $PLANFILE
`

	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	defer os.Remove(filePath)

	at := getMock()
	at.Request.AtlantisConfigFile = filePath

	err = at.ValidConfigFile()
	assert.NoError(t, err)

	_, err = at.SetConfigParmas()
	assert.NoError(t, err)

	assert.Equal(t, at.atlantisConfigParmas.Version, "3")
	assert.Equal(t, len(at.atlantisConfigParmas.Projects), 2)
	assert.Equal(t, at.atlantisConfigParmas.Projects[0].Name, "ec2")
	assert.Equal(t, at.atlantisConfigParmas.Projects[0].Dir, "examples/ec2")
	assert.Equal(t, at.atlantisConfigParmas.Projects[0].Workflow, "terraform")
	assert.Equal(t, at.atlantisConfigParmas.Projects[1].Name, "sg")
	assert.Equal(t, at.atlantisConfigParmas.Projects[1].Dir, "examples/sg")
	assert.Equal(t, at.atlantisConfigParmas.Projects[1].Workflow, "terraform")
}

func Test_ValidRepository(t *testing.T) {

	at := getMock()
	err := at.ValidRepository()

	assert.NoError(t, err)
}

func getMock() AtlantisParams {
	return AtlantisParams{
		Request: &usecase.AtlantisRequest{
			AtlantisURL:        "https://atlantis.example.com",
			AtlantisToken:      "atlantis_token",
			AtlantisRepository: "zkfmapf123/atlantis-fargate",
			AtlantisConfigFile: "atlantis_config_file.yaml",
		},
		httpClient: utils.NewATHTTP(),
	}
}
