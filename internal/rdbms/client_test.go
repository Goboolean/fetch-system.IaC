package rdbms_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/rdbms"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)



func SetupPostgreSQL() *rdbms.Client {
	c, err := rdbms.NewDB(&resolver.ConfigMap{
		"HOST": os.Getenv("POSTGRES_HOST"),
		"PORT": os.Getenv("POSTGRES_PORT"),
		"USER": os.Getenv("POSTGRES_USER"),
		"PASSWORD": os.Getenv("POSTGRES_PASSWORD"),
		"DATABASE": os.Getenv("POSTGRES_DATABASE"),
	})

	if err != nil {
		panic(err)
	}

	return c
}

func TeardownPostgreSQL(c *rdbms.Client) {
	c.Close()
}



func TestClient(t *testing.T) {

	c := SetupPostgreSQL()
	defer TeardownPostgreSQL(c)

	t.Run("Ping()", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := c.Ping(ctx)
		assert.NoError(t, err)
	})
}
