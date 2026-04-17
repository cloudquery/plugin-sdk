package scheduler

import "fmt"

// QueueType identifies a Storage backend.
type QueueType string

const (
	QueueTypeInMemory QueueType = "in-memory"
	QueueTypeBadger   QueueType = "badger"
)

// AllQueueTypes is used for error messages listing valid options.
var AllQueueTypes = []QueueType{QueueTypeInMemory, QueueTypeBadger}

// QueueConfig is the user-facing spec.queue configuration. Populated from
// source plugin spec. Validated during plugin spec.Validate().
type QueueConfig struct {
	// Type of backend. Defaults to QueueTypeInMemory when unset.
	Type QueueType `json:"type,omitempty"`
	// Path is the directory for the Badger backend. Required when Type=badger.
	Path string `json:"path,omitempty"`
}

// Validate checks backend-specific required fields.
func (c *QueueConfig) Validate() error {
	if c == nil {
		return nil
	}
	switch c.Type {
	case "", QueueTypeInMemory:
		return nil
	case QueueTypeBadger:
		if c.Path == "" {
			return fmt.Errorf("queue: type=%q requires path", QueueTypeBadger)
		}
		return nil
	default:
		return fmt.Errorf("queue: unknown type %q; supported: %v", c.Type, AllQueueTypes)
	}
}

// ValidateWithStrategy enforces that non-in-memory backends may only be used
// with the shuffle-queue scheduler strategy.
func (c *QueueConfig) ValidateWithStrategy(s Strategy) error {
	if c == nil {
		return nil
	}
	if err := c.Validate(); err != nil {
		return err
	}
	if c.Type == "" || c.Type == QueueTypeInMemory {
		return nil
	}
	if s != StrategyShuffleQueue {
		return fmt.Errorf("queue: type=%q requires scheduler=shuffle-queue (got %s)", c.Type, s.String())
	}
	return nil
}
