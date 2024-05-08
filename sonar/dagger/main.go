// SonarQube static code analysis.
//
// Check for code quality with static anaylsis scanning.
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Sonar struct{}

// Returns a container that echoes whatever string argument is provided
func (m *Sonar) Analyze(
	ctx context.Context,

	source *Directory,

	// sonar server url
	// +optional
	// +default="https://sonarcloud.io"
	url string,

	// sonar auth token
	token *Secret,

	// sonar project key
	// +optional
	projectKey string,

	// sonar organization
	// +optional
	organization string,

	// sonar scanner options
	// +optional
	options []string,
) (string, error) {
	ctr := dag.Container().
		From("sonarsource/sonar-scanner-cli").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithSecretVariable("SONAR_TOKEN", token)

	_, err := ctr.File("/src/sonar-project.properties").Contents(ctx)
	if err != nil {
		if projectKey == "" || organization == "" {
			return "", errors.New("sonar project key and organization are required")
		}
		options = append([]string{
			fmt.Sprintf("-Dsonar.projectKey=%s", projectKey),
			fmt.Sprintf("-Dsonar.organization=%s", organization),
		}, options...)
		ctr = ctr.
			WithEnvVariable("SONAR_HOST_URL", url).
			WithEnvVariable("SONAR_SCANNER_OPTS", strings.Join(options, " "))
	}

	return ctr.Stdout(ctx)
}
