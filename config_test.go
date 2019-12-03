package cfg_test

import (
	"testing"

	"github.com/indebted-modules/cfg"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestDatabaseURL() {
	s.Equal(cfg.DatabaseURL(), "db-url")
}
