//go:build e2e

package e2e

import (

	"testing"

	"github.com/stretchr/testify/suite"


)

type AuthRouterTestSuite struct {
	E2eTestSuite
	server *TestServer
}

// SetupSuite is called before any Router tests in the suite are run.
func (suite *AuthRouterTestSuite) SetupSuite() {
	suite.E2eTestSuite.SetupSuite()
	suite.server = NewTestServer(suite.db, suite.rdClient)
}

// TearDownSuite is called after all Router tests in the suite have been run, regardless of whether they passed or failed.
func (suite *AuthRouterTestSuite) TearDownSuite() {
	// Shutdown the server after a set of tests
    err := suite.server.Shutdown()
    suite.NoError(err)

    suite.E2eTestSuite.TearDownSuite()
}


func TestAuthRouterSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}
	suite.Run(t, new(AuthRouterTestSuite))
}
