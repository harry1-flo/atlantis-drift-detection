package usecase

type AtlantisRequest struct {
	GithubToken   string
	GithubRepoRef string

	AtlantisURL        string
	AtlantisToken      string
	AtlantisRepository string
	AtlantisConfigFile string
}

type APIHealthResponse struct {
	Status string `json:"status"`
}

type APIPlanBodyParams struct {
	Repository string             `json:"repository"`
	Ref        string             `json:"ref"`
	Type       string             `json:"type"` // Gtihub
	Paths      []APIPlanBodyPaths `json:"paths"`
}

type APIPlanBodyPaths struct {
	Directory string `json:"directory"`
	Workspace string `json:"workspace,omitempty"` // default
}

/*
@examples
## example
version: 3
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
*/
type AtlantisConfigParams struct {
	Version  string `yaml:"version"`
	Projects []struct {
		Name     string `yaml:"name"`
		Dir      string `yaml:"dir"`
		Workflow string `yaml:"workflow"`
	} `yaml:"projects"`
}
