package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"strings"
)

func ProcessHandler(c *gin.Context) {
	rtnMe := strings.Builder{}

	processes, err := process.Processes()

	if err != nil {
		return
	}

	openMetricsMetadata := `# TYPE process_cpu_percent Gauge
# HELP process_cpu_percent CPU percentage of a process.

# TYPE process_memory_rss_bytes Gauge
# UNIT process_memory_rss_bytes bytes
# HELP process_memory_rss_bytes total bytes a process is consuming, not including swap

`
	rtnMe.WriteString(openMetricsMetadata)

	for _, p := range processes {
		cmd, _ := p.CmdlineWithContext(c.Request.Context())
		if cmd == "" {
			cmd, _ = p.ExeWithContext(c.Request.Context())
		}
		pid := p.Pid
		ppid, _ := p.PpidWithContext(c.Request.Context())
		name, _ := p.NameWithContext(c.Request.Context())
		cpuPercent, _ := p.CPUPercentWithContext(c.Request.Context())
		mem, _ := p.MemoryInfoWithContext(c.Request.Context())
		memRss := mem.RSS
		//mem, _ := p.MemoryInfoWithContext(c.Request.Context())

		rtnMe.WriteString(fmt.Sprintf(`process_cpu_percent{name="%s",cmdLine="%s",pid="%d",ppid="%d"} %f`+"\n", name, cmd, pid, ppid, cpuPercent/100))
		rtnMe.WriteString(fmt.Sprintf(`process_memory_rss_bytes{name="%s",cmdLine="%s", pid="%d",ppid="%d"} %d`+"\n", name, cmd, pid, ppid, memRss))

	}

	c.Data(http.StatusOK, "application/openmetrics-text; version=1.0.0; charset=utf-8", []byte(rtnMe.String()))

}
