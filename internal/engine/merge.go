package engine

// GlobalMerger is responsible for global segment merging
type GlobalMerger struct{}

// NewGlobalMerger creates a new GlobalMerger instance
func NewGlobalMerger() *GlobalMerger {
	return &GlobalMerger{}
}

// Merge merges all segments globally
func (gm *GlobalMerger) Merge() error {
	// TODO: Implement the logic for global segment merging
	return nil
}
