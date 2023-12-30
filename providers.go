package planetscale

type (
	ControllerProvider struct {
		Product      ProductController
		ExpenseGroup ExpenseGroupController
		GroupMember  GroupMemberController
		Expense      ExpenseConroller
	}
	RepoProvider struct {
		Product      ProductRepo
		ExpenseGroup ExpenseGroupRepo
		GroupMember  GroupMemberRepo
		Expense      ExpenseRepo
	}
)
