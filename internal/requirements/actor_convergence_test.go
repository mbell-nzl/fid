package requirements

import (
	"fmt"
	"sync/atomic"
)

var retainedTestIdempotencySequence uint64

func testActorForRole(role string) string { return "agent:" + role }

func testIdempotencyKey() string {
	return fmt.Sprintf("retained-test-%d", atomic.AddUint64(&retainedTestIdempotencySequence, 1))
}
