package azure

import (
	"context"
	"fmt"
	"github.com/gerdooshell/tax-logger-client-go/internal/environment"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"

	"github.com/gerdooshell/tax-logger-client-go/internal/cache/lrucache"
	secretService "github.com/gerdooshell/tax-logger-client-go/internal/secret-service"
)

var secretServiceCache = lrucache.NewLRUCache[string](20)
var secretServiceMu sync.Mutex

func NewSecretService(uri string, environment environment.Environment) secretService.SecretService {
	secretServiceMu.Lock()
	defer secretServiceMu.Unlock()
	service, err := secretServiceCache.Read(uri)
	if err == nil {
		return service.(*azureSecretService)
	}
	newService := &azureSecretService{
		uri:         uri,
		cache:       lrucache.NewLRUCache[string](100),
		environment: environment,
	}
	secretServiceCache.Add(uri, newService)
	return newService
}

type azureSecretService struct {
	uri         string
	environment environment.Environment
	client      *azsecrets.Client
	cache       lrucache.LRUCache[string]
}

func (az *azureSecretService) connectToVault() (err error) {
	var cred azcore.TokenCredential
	if az.environment == environment.Dev {
		cred, err = azidentity.NewDefaultAzureCredential(nil)
	} else {
		cred, err = azidentity.NewManagedIdentityCredential(nil)
	}
	if err != nil {
		return err
	}
	options := azsecrets.ClientOptions{
		DisableChallengeResourceVerification: true,
	}
	client, err := azsecrets.NewClient(az.uri, cred, &options)
	if err != nil {
		return err
	}
	az.client = client
	return nil
}

func (az *azureSecretService) GetSecretValue(ctx context.Context, secretKey string) (<-chan string, <-chan error) {
	out := make(chan string, 1)
	errChan := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errChan)
		cachedValue, err := az.cache.Read(secretKey)
		if err == nil {
			fmt.Println("cached value:", cachedValue.(string))
			out <- cachedValue.(string)
			return
		}
		if err = az.connectToVault(); err != nil {
			fmt.Println("step 1:", err, secretKey)
			errChan <- err
			return
		}
		version := ""
		resp, err := az.client.GetSecret(ctx, secretKey, version, nil)
		if err != nil {
			fmt.Println("step 2:", err, secretKey)
			errChan <- err
			return
		}
		value := resp.Value
		if value == nil {
			fmt.Println("step 3:", fmt.Errorf("secret key not found"), secretKey)
			errChan <- fmt.Errorf("secret key not found")
			return
		}
		az.cache.Add(secretKey, *value)
		fmt.Println("fetched value:", *value)
		out <- *value
	}()

	return out, errChan
}
