package find

import (
	"context"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/src"
)

// FindTF finds the path to the terraform executable.
func FindTF(ctx context.Context, vc version.Constraints) (string, error) {
	i := install.NewInstaller()
	return i.Ensure(ctx, []src.Source{
		&fs.Version{
			Product:     product.Terraform,
			Constraints: vc,
		},
	})
}
