package polygon_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)




func SetupPolygon() *polygon.Client {

	p, err := polygon.New(&resolver.ConfigMap{
		"SECRET_KEY": os.Getenv("POLYGON_SECRET_KEY"),
	})
	if err != nil {
		panic(err)
	}

	return p
}



func TestGetAllProducts(t *testing.T) {
	p := SetupPolygon()

	list, err := p.GetAllProducts(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
}