package metric

import ()

type Counter struct {
	counts map[string]float64
	total  int
}

func (c *Counter) Initialize(keys ...string) error {
	var err error
	c.counts = make(map[string]float64)
	for _, k := range keys {
		c.counts[k] = float64(0)
	}
	return err
}

func (c *Counter) Count(value string) {
	cnt, ok := c.counts[value]
	if !ok {
		c.counts[value] = 1
	} else {
		c.counts[value] = cnt + 1
	}

	c.total += 1
}

func (c *Counter) Finalize() error {
	return nil
}

func (c *Counter) Result() MetricResult {
	return MetricResult{c.counts, c.total}
}
