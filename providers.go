package planetscale

type (
	ControllerProvider struct {
		Product      ProductController
		ExpenseGroup ExpenseGroupController
		GroupMember  GroupMemberController
		Expense      ExpenseConroller
		Settlement   SettlementController
		SplitType    SplitTypeController
		User         UserController
	}
	RepoProvider struct {
		Product            ProductRepo
		ExpenseGroup       ExpenseGroupRepo
		GroupMember        GroupMemberRepo
		Expense            ExpenseRepo
		ExpenseParticipant ExpenseParticipantRepo
		Settlement         SettlementRepo
		SplitType          SplitTypeRepo
		Item               ItemRepo
		ItemSplit          ItemSplitRepo
		User               UserRepo
	}

	ServiceProvider struct {
		Balance BalanceService
		Expense ExpenseService
	}
)
