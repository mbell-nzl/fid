package config

// InvalidConfigError identifies a supplied config file that cannot be used.
// Path names the input for CLI invalid-argument details; Err keeps the cause.
type InvalidConfigError struct {
	Path string
	Err  error
}

func (err *InvalidConfigError) Error() string { return err.Err.Error() }

func (err *InvalidConfigError) Unwrap() error { return err.Err }
