package resources

// IdentifierException mirrors net.minecraft.IdentifierException in package
// resources (the Java class lives in net.minecraft but is only referenced from
// resource-location parsing, so it is colocated here).
type IdentifierException struct {
	Message string
	Cause   error
}

// Error mirrors IdentifierException.getMessage().
func (e *IdentifierException) Error() string {
	return e.Message
}

// Unwrap mirrors Throwable.getCause().
func (e *IdentifierException) Unwrap() error {
	return e.Cause
}
