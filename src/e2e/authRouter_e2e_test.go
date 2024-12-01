//go:build e2e

package e2e

import (

	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/app/database"
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

func (suite *AuthRouterTestSuite) SetupTest() {
	// Create predetermined roles at the beginning of each test
	err := database.CreateDefaultRoles(suite.db)
	suite.NoError(err)
}

func (suite *AuthRouterTestSuite) TearDownTest() {

	// truncarTablesMaintainingTheStructureAndRestrictions
	tables := []string{
		"comment_votes",
		"posts_likes",
		"comments",
		"posts",
		"forums",
		"categories",
		"tokens",
		"users",
		"role_permissions",
		"roles",
	}

	for _, table := range tables {
		err := suite.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)).Error
		suite.NoError(err)
	}

	// Clean Redis at the end of each test
	err := suite.rdClient.FlushAll(suite.testCtx).Err()
	suite.NoError(err)
}





func TestAuthRouterSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}
	suite.Run(t, new(AuthRouterTestSuite))
}
