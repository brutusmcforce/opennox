//go:build !linux

package legacy

// Other platforms retain non-recoverable termination but are not runtime-qualified
// by this Linux port. legacyAbort performs the unconditional process exit.
func platformAbort() {}
