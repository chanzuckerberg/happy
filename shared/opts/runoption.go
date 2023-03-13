package opts

type RunOption struct {
	IsDestroy   *bool
	IsDryRun    *bool
	PlanMessage *string
	Targets     *[]string
}

func DestroyPlan() RunOption {
	val := true
	o := RunOption{
		IsDestroy: &val,
	}
	return o
}

func DryRun(dryRun bool) RunOption {
	val := dryRun
	o := RunOption{
		IsDryRun: &val,
	}
	return o
}

func Message(message string) RunOption {
	o := RunOption{
		PlanMessage: &message,
	}
	return o
}

func TargetAddrs(targets []string) RunOption {
	o := RunOption{
		Targets: &targets,
	}
	return o
}

func Combine(opts ...RunOption) RunOption {
	o := RunOption{}
	for _, opt := range opts {
		if opt.IsDestroy != nil {
			o.IsDestroy = opt.IsDestroy
		}
		if opt.IsDryRun != nil {
			o.IsDryRun = opt.IsDryRun
		}
		if opt.PlanMessage != nil {
			o.PlanMessage = opt.PlanMessage
		}
		if opt.Targets != nil {
			o.Targets = opt.Targets
		}
	}
	return o
}
