// Filename: internal/validator/validator.go

package validator

// We create a type that wraps our validation errors map
type Validator struct {
	Errors map[string]string
}

// New() creates a new Validator instance
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid() checks the Errors map for entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// In() checks if an element can be found in a provided list of elements
func In(element string, list ...string) bool {
	for i := range list {
		if element == list[i] {
			return true
		}
	}
	return false
}

// AddError() adds an error entry to the Errors map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check() performs the validation checks and calls the AddError()
// method in turn if an error entry needs to be added
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}