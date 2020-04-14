package params

// Default simulation operation weights for messages and gov proposals
const (
	DefaultWeightMsgCreateValidator int = 0
	DefaultWeightMsgEditValidator   int = 0
	DefaultWeightMsgDelegate        int = 0
	DefaultWeightMsgUndelegate      int = 0
	DefaultWeightMsgBeginRedelegate int = 0

	DefaultWeightCommunitySpendProposal int = 5
	DefaultWeightTextProposal           int = 5
	DefaultWeightParamChangeProposal    int = 5

	DefaultWeightMsgDefineService int = 100
)
