package test

import (
	"fmt"
	"testing"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/crypto"
	"github.com/authorizerdev/authorizer/server/envstore"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/resolvers"
	"github.com/stretchr/testify/assert"
)

func deleteUserTest(t *testing.T, s TestSetup) {
	t.Helper()
	t.Run(`should delete users with admin secret only`, func(t *testing.T) {
		req, ctx := createContext(s)
		email := "delete_user." + s.TestInfo.Email
		resolvers.SignupResolver(ctx, model.SignUpInput{
			Email:           email,
			Password:        s.TestInfo.Password,
			ConfirmPassword: s.TestInfo.Password,
		})

		_, err := resolvers.DeleteUserResolver(ctx, model.DeleteUserInput{
			Email: email,
		})
		assert.NotNil(t, err, "unauthorized")

		h, err := crypto.EncryptPassword(envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminSecret))
		assert.Nil(t, err)
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminCookieName), h))

		_, err = resolvers.DeleteUserResolver(ctx, model.DeleteUserInput{
			Email: email,
		})
		assert.Nil(t, err)
		cleanData(email)
	})
}
