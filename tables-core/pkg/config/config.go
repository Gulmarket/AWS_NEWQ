package config

import (
    "bufio"
    "crypto/tls"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/hashicorp/go-retryablehttp"
    "github.com/hashicorp/vault/api"
    "io"
    "log"
    "net"
    "net/http"
    "os"
    "reflect"
    "strconv"
    "time"
)

const (
    envTag = "env"
)

/*
configVaultClient function create default api.Client for Vault
Important note: address is fetching from Environment using os.Getenv function
Environment key - VAULT_ADDR

Returning api.Client, error
*/
func configVaultClient() (*api.Client, error) {
    vaultEnv := os.Getenv("VAULT_ADDR")
    if vaultEnv == "" {
        return nil, errors.New("VAULT_ADDR Environment variable is required")
    }
    cfg := &api.Config{
        Address: vaultEnv,
        HttpClient: &http.Client{
            Transport: &http.Transport{
                DialContext: (&net.Dialer{
                    Timeout:   30 * time.Second,
                    KeepAlive: 30 * time.Second,
                }).DialContext,
                MaxIdleConns:          100,
                IdleConnTimeout:       90 * time.Second,
                TLSHandshakeTimeout:   10 * time.Second,
                ExpectContinueTimeout: 1 * time.Second,
                TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
            },
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                return http.ErrUseLastResponse
            },
        },
        MinRetryWait: time.Millisecond * 1000,
        MaxRetryWait: time.Millisecond * 1500,
        MaxRetries:   2,
        Timeout:      time.Second * 60,
        Backoff:      retryablehttp.LinearJitterBackoff,
        CheckRetry:   retryablehttp.DefaultRetryPolicy,
    }
    return api.NewClient(cfg)
}

/*
fetchVaultConfig receives api.Client - creating by configVaultClient function
path - to Vault secret directory (pattern for Vault v2 - /secret/data/{secret_name})
token - for accessing Vault

If fetched api.Secret is nil - return error
if read from Vault failed - return error

Returning: map[string]interface{}, error
*/
func fetchVaultConfig(cli *api.Client, path, token string) (*api.Secret, error) {
    cli.SetToken(token)
    data, dataErr := cli.Logical().Read(path)
    if dataErr != nil {
        return nil, dataErr
    }
    if data == nil {
        return nil, fmt.Errorf("path %s not found", path)
    }
    return data, nil
}

/*
FetchVaultConfig process fetching secret using incoming arguments (path, token)
configVaultClient and fetchVaultConfig functions are using inside

Returning: map[string]interface{}, error
*/
func FetchVaultConfig(path, token string) (*api.Secret, error) {
    cli, configErr := configVaultClient()
    if configErr != nil {
        return nil, configErr
    }
    return fetchVaultConfig(cli, path, token)
}

/*
FetchBytesVaultConfig process fetching secret bytes using incoming arguments (path, token)
FetchVaultConfig function are using inside

Returning: []byte, error
*/
func FetchBytesVaultConfig(path, token string) ([]byte, error) {
    secret, err := FetchVaultConfig(path, token)
    if err != nil {
        return nil, err
    }
    data, mErr := json.Marshal(secret.Data["data"])
    if mErr != nil {
        return nil, mErr
    }
    return data, nil
}

/*
FetchBytesVaultConfigEnvToken process fetching secret bytes using incoming arguments (path)
and env variable VAULT_TOKEN for token
FetchBytesVaultConfig function are using inside

Returning: []byte, error
*/
func FetchBytesVaultConfigEnvToken(path string) ([]byte, error) {
    token := os.Getenv("VAULT_TOKEN")
    if len(token) == 0 {
        return nil, errors.New("token is not set in the environment variable VAULT_TOKEN")
    }
    return FetchBytesVaultConfig(path, token)
}

/*
FetchBytesVaultConfigEnvTokenPath process fetching secret bytes using env variables VAULT_PATH, VAULT_TOKEN, VAULT_ADDR
FetchBytesVaultConfigEnvToken function are using inside

Returning: []byte, error
*/
func FetchBytesVaultConfigEnvTokenPath() ([]byte, error) {
    path := os.Getenv("VAULT_PATH")
    if len(path) == 0 {
        return nil, errors.New("path is not set in the environment variable VAULT_PATH")
    }
    return FetchBytesVaultConfigEnvToken(path)
}

/*
ParseByteConfig takes array of bytes as input and cfg interface{} for unmarshal incoming data as json to structure
if array is empty - returns error, if cfg is nil - returns error

Returning: error

Important note: cfg and all inner struct field should be initialized as pointers

Example:
type Test struct {
    Inner *InnerTest `json:"inner"`
}

type InnerTest struct {
    Field string `json:"field"`
}

cfg := &Test{Inner: &InnerTest{}}

Byte array content:
{
    "inner": {
        "field":"value"
    }
}

*/
func ParseByteConfig(data []byte, cfg interface{}) error {
    if data == nil {
        return errors.New("data cannot be nil")
    }
    if cfg == nil {
        return errors.New("config struct cannot be nil")
    }
    if parseErr := json.Unmarshal(data, cfg); parseErr != nil {
        return parseErr
    }
    return nil
}

/*
ParseFileConfig takes filepath as input and cfg interface{} for unmarshal incoming data as json to structure
if filename empty - returns error, if cfg is nil - returns error
if file is not exists - returns error

Returning: error

Important note: cfg and all inner struct field should be initialized as pointers

Example:
type Test struct {
    Inner *InnerTest `json:"inner"`
}

type InnerTest struct {
    Field string `json:"field"`
}

cfg := &Test{Inner: &InnerTest{}}

File content:
{
    "inner": {
        "field":"value"
    }
}
*/

type JaegerEnv struct {
    AgentPort   string `json:"agentPort,omitempty"`
    AgentHost   string `json:"agentHost,omitempty"`
    ServiceName string ` json:"serviceName,omitempty"`
}

func SetJaegerConfig(j JaegerEnv) {
    if err := os.Setenv("JAEGER_AGENT_PORT", j.AgentPort); err != nil {
        log.Fatal(err)
    }
    if err := os.Setenv("JAEGER_AGENT_HOST", j.AgentHost); err != nil {
        log.Fatal(err)
    }
    if err := os.Setenv("JAEGER_SERVICE_NAME", j.ServiceName); err != nil {
        log.Fatal(err)
    }
    return
}
func ParseFileConfig(filepath string, cfg interface{}) error {
    if filepath == "" {
        return errors.New("filename cannot be empty")
    }
    if cfg == nil {
        return errors.New("config struct cannot be nil")
    }
    f, fErr := os.Open(filepath)
    if fErr != nil {
        return fErr
    }
    defer func() {
        if closeErr := f.Close(); closeErr != nil {
            log.Println("Failed to close resources")
        }
    }()
    reader := bufio.NewReader(f)
    data, readErr := io.ReadAll(reader)
    if readErr != nil {
        return readErr
    }
    if unmErr := json.Unmarshal(data, cfg); unmErr != nil {
        return unmErr
    }
    return nil
}

/*
ParseEnvConfig method takes config interface as argument
runs over the fields of incoming interface, if field contains tag env:"value"
after it is searching for "value" in os Environment
if find - inject os environment into field, if not - do nothing

Important note: cfg and all inner struct field should be initialized as pointers

Example:
type Test struct {
    Inner *InnerTest
}

type InnerTest struct {
    Field string `env:"field"`
}

cfg := &Test{Inner: &InnerTest{}}
*/
func ParseEnvConfig(cfg interface{}) {
    if cfg == nil {
        return
    }
    v := reflect.ValueOf(cfg)
    if v.Kind() == reflect.Ptr {
        el := v.Elem()
        for i := 0; i < el.NumField(); i++ {
            if el.Field(i).Kind() == reflect.Ptr {
                ParseEnvConfig(el.Field(i).Interface())
            } else {
                t := reflect.TypeOf(cfg).Elem()
                tagenv := t.Field(i).Tag.Get(envTag)
                env := os.Getenv(tagenv)
                if env == "" {
                    env = t.Field(i).Tag.Get("default")
                    if env == "" {
                        continue
                    }
                }
                switch el.Field(i).Kind() {
                case reflect.String:
                    el.Field(i).SetString(env)
                case reflect.Int:
                    num, _ := strconv.Atoi(env)
                    el.Field(i).SetInt(int64(num))
                case reflect.Bool:
                    b, _ := strconv.ParseBool(env)
                    el.Field(i).SetBool(b)
                }
            }
        }
    }
}

func FromEnvAndHelms(cfg interface{}) error {
    path := os.Getenv("CONFIG_PATH")
    if len(path) == 0 {
        path = "./helm/local"
    }
    ParseEnvConfig(cfg)
    if err := ParseFileConfig(fmt.Sprintf("%s/%s", path, "config.json"), cfg); err != nil {
        return err
    }
    return nil
}
