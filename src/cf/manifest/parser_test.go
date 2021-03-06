package manifest

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testhelpers/maker"
	"testing"
)

func TestParsingApplicationName(t *testing.T) {
	m, err := Parse(strings.NewReader(maker.ManifestWithName("single app")))
	assert.NoError(t, err)
	assert.Equal(t, "manifest-app-name", m.Applications[0].Get("name").(string))
}

func TestParsingManifestServices(t *testing.T) {
	m, err := Parse(strings.NewReader(maker.ManifestWithName("global services")))
	assert.NoError(t, err)

	services := m.Applications[0].Get("services").([]string)
	assert.Equal(t, len(services), 1)
	assert.Equal(t, services[0], "work-queue")
}

func TestParsingManifestWithEmptyEnvVar(t *testing.T) {
	_, err := Parse(strings.NewReader(maker.ManifestWithName("invalid env")))
	assert.Error(t, err)
}
