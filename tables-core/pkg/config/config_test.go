package config

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "os"
    "testing"
)

const (
    validURL        = "msa/data/test"
    validToken      = "s.l5Wk4daFYgqQd0f2AGYvQBVf"
    invalidToken    = "invalid-token"
    vaultAddr       = "http://10.4.145.12:8200"
    testEnvValue    = "test"
    validFileName   = "test.json"
    invalidFileName = "invalid.test.json"
)

type TestVaultConfig struct {
    VaultAddr string `json:"vault_addr" env:"VAULT_ADDR"`
}

type TestConfig struct {
    Vault *TestVaultConfig `json:"vault"`
    Test  string           `json:"test" env:"TEST_ENV"`
}

func TestFetchBytesVaultConfigEnvTokenPath(t *testing.T) {
    envErr := os.Setenv("VAULT_ADDR", vaultAddr)
    envErr = os.Setenv("VAULT_TOKEN", validToken)
    envErr = os.Setenv("VAULT_PATH", validURL)

    defer func() {
        os.Unsetenv("VAULT_ADDR")
        os.Unsetenv("VAULT_TOKEN")
        os.Unsetenv("VAULT_PATH")
    }()

    data, err := FetchBytesVaultConfigEnvTokenPath()

    assert.Nil(t, envErr)
    assert.Nil(t, err)
    assert.NotNil(t, data)

    fmt.Printf("%s\n", string(data))
}

func NewTestConfig() *TestConfig {
    return &TestConfig{Vault: &TestVaultConfig{}}
}

func TestFetchVaultConfig(t *testing.T) {
    os.Setenv("VAULT_ADDR", vaultAddr)

    data, err := FetchVaultConfig(validURL, validToken)

    assert.NotNil(t, data)
    assert.Nil(t, err)
    assert.True(t, len(data.Data) > 0)
}

func TestFetchVaultConfig_Failed(t *testing.T) {
    data, err := FetchVaultConfig("test", validToken)

    assert.Nil(t, data)
    assert.NotNil(t, err)
}

func TestFetchVaultConfig_OsEnvNotFound(t *testing.T) {
    os.Setenv("VAULT_ADDR", "")
    defer os.Setenv("VAULT_ADDR", vaultAddr)

    data, err := FetchVaultConfig(validURL, validToken)

    assert.Nil(t, data)
    assert.NotNil(t, err)
}

func TestFetchVaultConfig_TokenIsIncorrect(t *testing.T) {
    data, err := FetchVaultConfig(validURL, invalidToken)

    assert.Nil(t, data)
    assert.NotNil(t, err)
}

func TestParseEnvConfig(t *testing.T) {
    os.Setenv("VAULT_ADDR", vaultAddr)
    os.Setenv("TEST_ENV", testEnvValue)

    cfg := NewTestConfig()

    ParseEnvConfig(cfg)

    assert.Equal(t, cfg.Test, testEnvValue)
    assert.Equal(t, cfg.Vault.VaultAddr, vaultAddr)
}

func TestParseFileConfig(t *testing.T) {
    cfg := NewTestConfig()

    parseErr := ParseFileConfig(validFileName, cfg)

    assert.Nil(t, parseErr)
    assert.Equal(t, cfg.Test, testEnvValue)
    assert.Equal(t, cfg.Vault.VaultAddr, vaultAddr)
}

func TestParseFileConfig_NoFile(t *testing.T) {
    cfg := NewTestConfig()

    parseErr := ParseFileConfig(invalidFileName, cfg)

    assert.NotNil(t, parseErr)
    assert.NotEqual(t, cfg.Vault.VaultAddr, vaultAddr)
    assert.NotEqual(t, cfg.Test, testEnvValue)
}

func TestParseFileConfig_EmptyFilename(t *testing.T) {
    cfg := NewTestConfig()

    parseErr := ParseFileConfig("", cfg)

    assert.NotNil(t, parseErr)
    assert.NotEqual(t, cfg.Vault.VaultAddr, vaultAddr)
    assert.NotEqual(t, cfg.Test, testEnvValue)
}

func TestParseFileConfig_NilConfig(t *testing.T) {
    cfg := NewTestConfig()

    parseErr := ParseFileConfig(validFileName, nil)

    assert.NotNil(t, parseErr)
    assert.NotEqual(t, cfg.Vault.VaultAddr, vaultAddr)
    assert.NotEqual(t, cfg.Test, testEnvValue)
}

func TestParseByteConfig(t *testing.T) {
    cfg := NewTestConfig()

    data, readErr := ioutil.ReadFile(validFileName)

    parseErr := ParseByteConfig(data, cfg)

    assert.Nil(t, readErr)
    assert.Nil(t, parseErr)
    assert.Equal(t, cfg.Test, testEnvValue)
    assert.Equal(t, cfg.Vault.VaultAddr, vaultAddr)
}

func TestParseByteConfig_NoFile(t *testing.T) {
    cfg := NewTestConfig()

    data, readErr := ioutil.ReadFile(invalidFileName)

    parseErr := ParseByteConfig(data, cfg)

    assert.NotNil(t, readErr)
    assert.NotNil(t, parseErr)
    assert.NotEqual(t, cfg.Test, testEnvValue)
    assert.NotEqual(t, cfg.Vault.VaultAddr, vaultAddr)
}

func TestParseByteConfig_NilData(t *testing.T) {
    cfg := NewTestConfig()

    parseErr := ParseByteConfig(nil, cfg)

    assert.NotNil(t, parseErr)
    assert.NotEqual(t, cfg.Test, testEnvValue)
    assert.NotEqual(t, cfg.Vault.VaultAddr, vaultAddr)
}

func TestParseByteConfig_NilConfig(t *testing.T) {
    cfg := NewTestConfig()

    data, readErr := ioutil.ReadFile(validFileName)

    parseErr := ParseByteConfig(data, nil)

    assert.Nil(t, readErr)
    assert.NotNil(t, parseErr)
    assert.NotEqual(t, cfg.Test, testEnvValue)
    assert.NotEqual(t, cfg.Vault.VaultAddr, vaultAddr)
}
