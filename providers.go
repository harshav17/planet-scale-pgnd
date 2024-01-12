package planetscale

type (
	ControllerProvider struct {
		Product      ProductController
		ExpenseGroup ExpenseGroupController
		GroupMember  GroupMemberController
		Expense      ExpenseConroller
		Settlement   SettlementController
	}
	RepoProvider struct {
		Product      ProductRepo
		ExpenseGroup ExpenseGroupRepo
		GroupMember  GroupMemberRepo
		Expense      ExpenseRepo
		Settlement   SettlementRepo
		SplitType    SplitTypeRepo
	}
)
