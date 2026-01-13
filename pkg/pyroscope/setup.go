package pyroscope

import (
	"github.com/grafana/pyroscope-go"
)

func Setup(appName, serverAddress string) (*pyroscope.Profiler, error) {
	pyroProf, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: appName,
		// the OTel collector is actually Grafana Alloy, which has a built-in Pyroscope receiver
		ServerAddress: serverAddress,
		// you can disable logging by setting this to nil
		Logger:   nil, //pyroscope.StandardLogger, -- disabled for less noisy logs
		TenantID: "tenant_123",
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		return nil, err
	}

	return pyroProf, nil
}
