package diff

// Pipeline chains together the full diff processing sequence:
// Compare → Filter → Sort → ApplyContext → ApplyMask.
//
// This provides a single entry point for callers (e.g. cmd/vaultdiff/main.go)
// that want to run the standard transformation sequence without wiring each
// step manually.

// PipelineOptions aggregates every option struct used across the pipeline.
type PipelineOptions struct {
	Filter  FilterOptions
	Sort    SortOptions
	Context ContextOptions
	Mask    MaskOptions
}

// DefaultPipelineOptions returns a PipelineOptions with each sub-option set
// to its own documented default.
func DefaultPipelineOptions() PipelineOptions {
	return PipelineOptions{
		Filter:  DefaultFilterOptions(),
		Sort:    DefaultSortOptions(),
		Context: DefaultContextOptions(),
		Mask:    DefaultMaskOptions(),
	}
}

// Run executes the full diff pipeline against two secret maps and returns the
// processed slice of Change values ready for rendering.
//
// Steps:
//  1. Compare  – produce raw Change slice from secretA and secretB.
//  2. Filter   – drop change types the caller did not request.
//  3. Sort     – order the remaining changes deterministically.
//  4. Context  – pad unchanged lines around real changes (unified-diff style).
//  5. Mask     – redact sensitive key values before they reach output.
func Run(
	secretA, secretB map[string]interface{},
	opts PipelineOptions,
) []Change {
	changes := Compare(secretA, secretB)
	changes = Filter(changes, opts.Filter)
	changes = Sort(changes, opts.Sort)
	changes = ApplyContext(changes, opts.Context)
	changes = ApplyMask(changes, opts.Mask)
	return changes
}
