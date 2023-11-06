package planetscale

type (
	ControllerProvider struct {
		Product ProductController
	}
	RepoProvider struct {
		Product ProductRepo
	}
)
