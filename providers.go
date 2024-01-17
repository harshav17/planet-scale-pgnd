package planetscale

type (
	ControllerProvider struct {
		Product      ProductController
		ExpenseGroup ExpenseGroupController
		GroupMember  GroupMemberController
		Expense      ExpenseConroller
		Settlement   SettlementController
		SplitType    SplitTypeController
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
	}

	ServiceProvider struct {
		Balance BalanceService
	}
)
