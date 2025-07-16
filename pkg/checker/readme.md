
``` go

err := checker.RunHealthCheck(ctx, monitor.Type, monitor.Target)
if err != nil {
	log.Warn("Monitor check failed", zap.String("monitor_id", monitor.ID), zap.Error(err))
} else {
	log.Info("Monitor is healthy", zap.String("monitor_id", monitor.ID))
}



```