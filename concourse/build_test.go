package concourse

import (
	"os"
	"reflect"
	"testing"
)

func TestNewBuildMetadata(t *testing.T) {
	env := map[string]string{
		"ATC_EXTERNAL_URL":    "https://ci.example.com",
		"BUILD_TEAM_NAME":     "main",
		"BUILD_PIPELINE_NAME": "demo",
		"BUILD_JOB_NAME":      "my test",
		"BUILD_NAME":          "1",
	}

	cases := map[string]struct {
		host         string
		instanceVars string
		want         BuildMetadata
	}{
		"environment only": {
			want: BuildMetadata{
				Host:         "https://ci.example.com",
				TeamName:     "main",
				PipelineName: "demo",
				InstanceVars: "",
				JobName:      "my test",
				BuildName:    "1",
				URL:          "https://ci.example.com/teams/main/pipelines/demo/jobs/my%20test/builds/1",
			},
		},
		"url override": {
			host: "https://example.com",
			want: BuildMetadata{
				Host:         "https://example.com",
				TeamName:     "main",
				PipelineName: "demo",
				InstanceVars: "",
				JobName:      "my test",
				BuildName:    "1",
				URL:          "https://example.com/teams/main/pipelines/demo/jobs/my%20test/builds/1",
			},
		},
		"url with instance vars": {
			instanceVars: `{"image_name":"my-image","pr_number":"1234"}`,
			want: BuildMetadata{
				Host:         "https://ci.example.com",
				TeamName:     "main",
				PipelineName: "demo",
				InstanceVars: `{"image_name":"my-image","pr_number":"1234"}`,
				JobName:      "my test",
				BuildName:    "1",
				URL:          `https://ci.example.com/teams/main/pipelines/demo/jobs/my%20test/builds/1?vars.image_name=%22my-image%22&vars.pr_number=%221234%22`,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			for k, v := range env {
				os.Setenv(k, v)
			}
			if c.instanceVars != "" {
				os.Setenv("BUILD_PIPELINE_INSTANCE_VARS", c.instanceVars)
			} else {
				os.Unsetenv("BUILD_PIPELINE_INSTANCE_VARS")
			}

			metadata := NewBuildMetadata(c.host)
			if !reflect.DeepEqual(metadata, c.want) {
				t.Fatalf("unexpected BuildMetadata value from GetBuildMetadata:\n\t(GOT): %#v\n\t(WNT): %#v", metadata, c.want)
			}
		})
	}
}
