package planetscale

type (
	ControllerProvider struct {
		Product      ProductController
		ExpenseGroup ExpenseGroupController
		GroupMember  GroupMemberController
	}
	RepoProvider struct {
		Product      ProductRepo
		ExpenseGroup ExpenseGroupRepo
		GroupMember  GroupMemberRepo
	}
)
