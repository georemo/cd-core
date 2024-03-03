// cd custom error

package base

type CdError struct{}

func (m *CdError) Error() string {
	return "Unexpected Error!"
}
