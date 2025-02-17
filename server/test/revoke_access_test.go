package test

import (
	"fmt"
	"testing"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/crypto"
	"github.com/authorizerdev/authorizer/server/db"
	"github.com/authorizerdev/authorizer/server/envstore"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/resolvers"
	"github.com/stretchr/testify/assert"
)

func revokeAccessTest(t *testing.T, s TestSetup) {
	t.Helper()
	t.Run(`should revoke access`, func(t *testing.T) {
		req, ctx := createContext(s)
		email := "revoke_access." + s.TestInfo.Email
		_, err := resolvers.MagicLinkLoginResolver(ctx, model.MagicLinkLoginInput{
			Email: email,
		})
		assert.NoError(t, err)
		verificationRequest, err := db.Provider.GetVerificationRequestByEmail(email, constants.VerificationTypeMagicLinkLogin)
		verifyRes, err := resolvers.VerifyEmailResolver(ctx, model.VerifyEmailInput{
			Token: verificationRequest.Token,
		})
		assert.NoError(t, err)
		assert.NotNil(t, verifyRes.AccessToken)

		res, err := resolvers.RevokeAccessResolver(ctx, model.UpdateAccessInput{
			UserID: verifyRes.User.ID,
		})
		assert.Error(t, err)

		h, err := crypto.EncryptPassword(envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminSecret))
		assert.Nil(t, err)
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminCookieName), h))

		res, err = resolvers.RevokeAccessResolver(ctx, model.UpdateAccessInput{
			UserID: verifyRes.User.ID,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Message)

		// it should not allow login with revoked access
		_, err = resolvers.MagicLinkLoginResolver(ctx, model.MagicLinkLoginInput{
			Email: email,
		})
		assert.Error(t, err)
		cleanData(email)
	})
}
