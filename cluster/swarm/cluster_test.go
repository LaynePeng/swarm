package swarm

import (
	"testing"

	"github.com/docker/swarm/cluster"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

func createEngine(t *testing.T, ID string, containers ...dockerclient.Container) *cluster.Engine {
	engine := cluster.NewEngine(ID, 0)
	engine.Name = ID
	engine.ID = ID

	for _, container := range containers {
		engine.AddContainer(&cluster.Container{Container: container, Engine: engine})
	}

	return engine
}

func TestContainerLookup(t *testing.T) {
	c := &Cluster{
		engines: make(map[string]*cluster.Engine),
	}
	container := dockerclient.Container{
		Id:    "container-id",
		Names: []string{"/container-name1", "/container-name2"},
	}

	n := createEngine(t, "test-engine", container)
	c.engines[n.ID] = n

	// Invalid lookup
	assert.Nil(t, c.Container("invalid-id"))
	assert.Nil(t, c.Container(""))
	// Container ID lookup.
	assert.NotNil(t, c.Container("container-id"))
	// Container ID prefix lookup.
	assert.NotNil(t, c.Container("container-"))
	// Container name lookup.
	assert.NotNil(t, c.Container("container-name1"))
	assert.NotNil(t, c.Container("container-name2"))
	// Container engine/name matching.
	assert.NotNil(t, c.Container("test-engine/container-name1"))
	assert.NotNil(t, c.Container("test-engine/container-name2"))
}
