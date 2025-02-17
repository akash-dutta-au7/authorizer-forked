package test

import (
	"testing"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/envstore"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/resolvers"
	"github.com/stretchr/testify/assert"
)

func adminSignupTests(t *testing.T, s TestSetup) {
	t.Helper()
	t.Run(`should complete admin signup`, func(t *testing.T) {
		_, ctx := createContext(s)
		_, err := resolvers.AdminSignupResolver(ctx, model.AdminSignupInput{
			AdminSecret: "admin",
		})

		assert.NotNil(t, err)
		// reset env for test to pass
		envstore.EnvStoreObj.UpdateEnvVariable(constants.StringStoreIdentifier, constants.EnvKeyAdminSecret, "")

		_, err = resolvers.AdminSignupResolver(ctx, model.AdminSignupInput{
			AdminSecret: "admin123",
		})

		assert.Nil(t, err)
	})
}
