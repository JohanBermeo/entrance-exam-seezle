package calculation

// Result holds the computed value for an operation.
type Result struct {
	ID    string  `json:"id"`
	Value float64 `json:"value"`
	Error string  `json:"error,omitempty"`
}

// Results is a collection of results indexed by operation ID.
type Results map[string]Result

// Get returns the result for the given operation ID.
func (r Results) Get(id string) (Result, bool) {
	res, ok := r[id]
	return res, ok
}

// Set stores a result for the given operation ID.
func (r Results) Set(id string, value float64) {
	r[id] = Result{ID: id, Value: value}
}

// SetError stores an error result for the given operation ID.
func (r Results) SetError(id string, err string) {
	r[id] = Result{ID: id, Error: err}
}

// HasError returns true if the result has an error.
func (r Result) HasError() bool {
	return r.Error != ""
}

// OutputValues returns the values for the requested output IDs in order.
func (r Results) OutputValues(outputIDs []string) []float64 {
	values := make([]float64, len(outputIDs))
	for i, id := range outputIDs {
		if res, ok := r[id]; ok && !res.HasError() {
			values[i] = res.Value
		}
	}
	return values
}

// OutputErrors returns errors for the requested output IDs in order.
func (r Results) OutputErrors(outputIDs []string) []string {
	errors := make([]string, len(outputIDs))
	for i, id := range outputIDs {
		if res, ok := r[id]; ok && res.HasError() {
			errors[i] = res.Error
		}
	}
	return errors
}
