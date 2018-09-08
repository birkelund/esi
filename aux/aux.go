package aux

type lazyResource struct {
	loaded bool
}

func (r *lazyResource) load(loader func()) {
	loader()
}
