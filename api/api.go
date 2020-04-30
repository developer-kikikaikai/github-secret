package api

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v31/github"
	"github.com/jamesruan/sodium"
	"golang.org/x/oauth2"
)

type APISecret interface {
	Update(name, secret string) error
	Delete(name string) error
	Get(name string) (*github.Secret, error)
	List() (*github.Secrets, error)
}

type apiSecret struct {
	client *github.Client
	owner  string
	repo   string
}

func NewSecretClient(owner, repo, token string) APISecret {
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oc := oauth2.NewClient(context.Background(), sts)
	return &apiSecret{
		client: github.NewClient(oc),
		owner:  owner,
		repo:   repo,
	}
}

func (api *apiSecret) ctx() context.Context {
	return context.Background()
}

func (api *apiSecret) getActions() *github.ActionsService {
	return api.client.Actions
}

func (api *apiSecret) Update(name, secret string) error {
	actions := api.getActions()
	ctx := api.ctx()

	// Get public key for secret
	pubKey, _, err := actions.GetPublicKey(ctx, api.owner, api.repo)
	if err != nil {
		fmt.Println("error GetPublicKey:", err.Error())
		return err
	}
	pkBase64 := pubKey.GetKey()
	_pk := make([]byte, len(pkBase64))
	_, err = base64.StdEncoding.Decode(_pk, []byte(pkBase64))
	if err != nil {
		fmt.Println("base64 decode error:", pkBase64, err.Error())
		return err
	}
	pk := make([]byte, 32)
	for i, _ := range pk {
		pk[i] = _pk[i]
	}

	// encrypts by sodium
	encSec := sodium.Bytes(secret).SealedBox(sodium.BoxPublicKey{Bytes: pk})
	encSecBase64 := base64.StdEncoding.EncodeToString(encSec)

	// Update secret
	es := &github.EncryptedSecret{
		Name:           name,
		KeyID:          *pubKey.KeyID,
		EncryptedValue: encSecBase64,
	}
	_, err = actions.CreateOrUpdateSecret(ctx, api.owner, api.repo, es)
	if err != nil {
		fmt.Println("failed to update secret", err.Error())
		return err
	}
	return nil
}

func (api *apiSecret) Delete(name string) error {
	_, err := api.getActions().DeleteSecret(api.ctx(), api.owner, api.repo, name)
	return err
}

func (api *apiSecret) Get(name string) (*github.Secret, error) {
	secret, _, err := api.getActions().GetSecret(api.ctx(), api.owner, api.repo, name)
	return secret, err
}

func (api *apiSecret) List() (*github.Secrets, error) {
	secrets, _, err := api.getActions().ListSecrets(api.ctx(), api.owner, api.repo, nil)
	return secrets, err
}
