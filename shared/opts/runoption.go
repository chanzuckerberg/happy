package opts

type RunOption struct {
	destroy *bool
	dryRun  *bool
	message *string
	targets *[]string
}

func (o RunOption) IsDestroy() *bool {
	return o.destroy
}

func (o RunOption) IsDryRun() *bool {
	return o.dryRun
}

func (o RunOption) Message() *string {
	return o.message
}

func (o RunOption) Targets() *[]string {
	return o.targets
}

func DestroyPlan() RunOption {
	val := true
	o := RunOption{
		destroy: &val,
	}
	return o
}

func DryRun(dryRun bool) RunOption {
	val := dryRun
	o := RunOption{
		dryRun: &val,
	}
	return o
}

func Message(message string) RunOption {
	o := RunOption{
		message: &message,
	}
	return o
}

func TargetAddrs(targets []string) RunOption {
	o := RunOption{
		targets: &targets,
	}
	return o
}

func Combine(opts ...RunOption) RunOption {
	o := RunOption{}
	for _, opt := range opts {
		if opt.destroy != nil {
			o.destroy = opt.destroy
		}
		if opt.dryRun != nil {
			o.dryRun = opt.dryRun
		}
		if opt.message != nil {
			o.message = opt.message
		}
		if opt.targets != nil {
			o.targets = opt.targets
		}
	}
	return o
}
