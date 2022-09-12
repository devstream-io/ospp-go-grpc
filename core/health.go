package core

import (
	"time"
)

func (c *Core) healthCheck() func() {
	return func() {
		c.plugins.Range(func(key, value interface{}) bool {
			_, ok1 := key.(string)
			v, ok2 := value.(*pluginInfo)
			if !ok1 || !ok2 {
				return true // continue
			}

			t := time.Now().Unix() - int64(c.opts.healthTimeout.Seconds())
			c.opts.logger.Logf("core", LogLevelDebug, "checking health of plugin %s.%s: %d/%d", v.name, v.version, v.health, t)

			// if health timeout, shutdown plugin
			if v.health < t {
				c.opts.logger.Logf("core", LogLevelInfo, "shutdown plugin %s.%s", v.name, v.version)
				if err := c.ShutdownPlugin(v.name, v.version); err != nil {
					c.opts.logger.Logf("core", LogLevelError, "shutdown plugin %s.%s failed: %s", v.name, v.version, err.Error())
					return false
				}
			}
			return true // continue
		})
	}
}
