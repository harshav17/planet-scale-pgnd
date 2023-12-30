package planetscale

type (
	ControllerProvider struct {
		Product      ProductController
		ExpenseGroup ExpenseGroupController
	}
	RepoProvider struct {
		Product      ProductRepo
		ExpenseGroup ExpenseGroupRepo
	}
)
