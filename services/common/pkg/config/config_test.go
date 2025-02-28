package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMustLoadConfig(t *testing.T) {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("SERVER_PORT", "9090")
	_ = os.Setenv("SERVER_ADDR", "127.0.0.1")
	_ = os.Setenv("SERVER_READ_TIMEOUT", "5s")
	_ = os.Setenv("SERVER_WRITE_TIMEOUT", "10s")
	_ = os.Setenv("SERVER_IDLE_TIMEOUT", "30s")
	_ = os.Setenv("DB_USER", "testuser")
	_ = os.Setenv("DB_PASSWORD", "testpass")
	_ = os.Setenv("DB_HOST", "testhost")
	_ = os.Setenv("DB_PORT", "1234")
	_ = os.Setenv("DB_NAME", "testdb")
	_ = os.Setenv("PRIVATE_KEY", "test_private_key")
	_ = os.Setenv("PUBLIC_KEY", "test_public_key")
	_ = os.Setenv("GOOGLE_CLIENT_ID", "test_google_client_id")
	_ = os.Setenv("GOOGLE_CLIENT_SECRET", "test_google_client_secret")
	_ = os.Setenv("TOKEN_TTL", "20m")

	config := MustLoadConfig()

	assert.Equal(t, "test", config.Env, "should match the test ENV value")
	assert.Equal(t, 9090, config.Server.Port, "should match the test SERVER_PORT value")
	assert.Equal(t, "127.0.0.1", config.Server.Addr, "should match the test SERVER_ADDR value")
	assert.Equal(t, 5*time.Second, config.Server.ReadTimeout, "should match the test SERVER_READ_TIMEOUT value")
	assert.Equal(t, 10*time.Second, config.Server.WriteTimeout, "should match the test SERVER_WRITE_TIMEOUT value")
	assert.Equal(t, 30*time.Second, config.Server.IdleTimeout, "should match the test SERVER_IDLE_TIMEOUT value")

	assert.Equal(t, "testuser", config.Database.User, "should match the test DB_USER value")
	assert.Equal(t, "testpass", config.Database.Password, "should match the test DB_PASSWORD value")
	assert.Equal(t, "testhost", config.Database.Host, "should match the test DB_HOST value")
	assert.Equal(t, "1234", config.Database.Port, "should match the test DB_PORT value")
	assert.Equal(t, "testdb", config.Database.Name, "should match the test DB_NAME value")

	assert.Equal(t, "test_private_key", config.PrivateKey, "should match the test PRIVATE_KEY value")
	assert.Equal(t, "test_public_key", config.PublicKey, "should match the test PUBLIC_KEY value")
	assert.Equal(t, "test_google_client_id", config.GoogleClientID, "should match the test GOOGLE_CLIENT_ID value")
	assert.Equal(t, "test_google_client_secret", config.GoogleClientSecret, "should match the test GOOGLE_CLIENT_SECRET value")
	assert.Equal(t, 20*time.Minute, config.TokenTTL, "should match the test TOKEN_TTL value")

	_ = os.Unsetenv("ENV")
	_ = os.Unsetenv("SERVER_PORT")
	_ = os.Unsetenv("SERVER_ADDR")
	_ = os.Unsetenv("SERVER_READ_TIMEOUT")
	_ = os.Unsetenv("SERVER_WRITE_TIMEOUT")
	_ = os.Unsetenv("SERVER_IDLE_TIMEOUT")
	_ = os.Unsetenv("DB_USER")
	_ = os.Unsetenv("DB_PASSWORD")
	_ = os.Unsetenv("DB_HOST")
	_ = os.Unsetenv("DB_PORT")
	_ = os.Unsetenv("DB_NAME")
	_ = os.Unsetenv("PRIVATE_KEY")
	_ = os.Unsetenv("PUBLIC_KEY")
	_ = os.Unsetenv("GOOGLE_CLIENT_ID")
	_ = os.Unsetenv("GOOGLE_CLIENT_SECRET")
	_ = os.Unsetenv("TOKEN_TTL")
}

