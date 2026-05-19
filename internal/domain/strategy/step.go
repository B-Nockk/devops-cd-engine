package strategy

type ActionType string

const (
	ActionPullArtifact   ActionType = "pull_artifact"
	ActionStartContainer ActionType = "start_container"
	ActionCheckHealth    ActionType = "check_health"
	ActionSwitchTraffic  ActionType = "switch_traffic"
	ActionStopContainer  ActionType = "stop_container"
)

type Step struct {
	Name         string
	ActionType   ActionType
	Target       string
	RollbackStep *Step
}
