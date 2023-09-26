package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func ProcessHandler(c *gin.Context) {
	topNProcessesString := c.DefaultQuery("top", "all")
	var topNProcesses = -1
	if strings.ToLower(topNProcessesString) != "all" {
		i, err := strconv.ParseInt(topNProcessesString, 10, 64)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		topNProcesses = int(i)
	}

	rtnMe := strings.Builder{}

	processes, err := process.Processes()

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	openMetricsMetadata := `# TYPE process_cpu_percent Gauge
# HELP process_cpu_percent CPU percentage of a process.

# TYPE process_memory_rss_bytes Gauge
# UNIT process_memory_rss_bytes bytes
# HELP process_memory_rss_bytes total bytes a process is consuming, not including swap

`
	rtnMe.WriteString(openMetricsMetadata)

	// cpu
	sort.Slice(processes, func(i, j int) bool {
		iPercent, err2 := processes[i].CPUPercentWithContext(c.Request.Context())
		if err2 != nil {
			panic(err2)
		}
		jPercent, err2 := processes[j].CPUPercentWithContext(c.Request.Context())
		if err2 != nil {
			panic(err2)
		}
		return iPercent < jPercent
	})

	for i, p := range processes {
		cmd, _ := p.CmdlineWithContext(c.Request.Context())
		if cmd == "" {
			cmd, _ = p.ExeWithContext(c.Request.Context())
		}
		cmd = strings.ReplaceAll(cmd, "\"", "\\\"")
		pid := p.Pid
		ppid, _ := p.PpidWithContext(c.Request.Context())
		name, _ := p.NameWithContext(c.Request.Context())
		name = strings.ReplaceAll(name, "\"", "\\\"")
		cpuPercent, _ := p.CPUPercentWithContext(c.Request.Context())

		rtnMe.WriteString(fmt.Sprintf(`process_cpu_percent{name="%s",cmdLine="%s",pid="%d",ppid="%d"} %f`+"\n", name, cmd, pid, ppid, cpuPercent))

		if i == topNProcesses-1 {
			break
		}
	}
	// memory
	sort.Slice(processes, func(i, j int) bool {
		iMem, err2 := processes[i].MemoryInfoWithContext(c.Request.Context())
		if err2 != nil {
			panic(err2)
		}
		iMemVal := iMem.RSS
		jMem, err2 := processes[j].MemoryInfoWithContext(c.Request.Context())
		if err2 != nil {
			panic(err2)
		}
		jMemVal := jMem.RSS
		return iMemVal < jMemVal
	})

	for i, p := range processes {
		cmd, _ := p.CmdlineWithContext(c.Request.Context())
		if cmd == "" {
			cmd, _ = p.ExeWithContext(c.Request.Context())
		}
		pid := p.Pid
		ppid, _ := p.PpidWithContext(c.Request.Context())
		name, _ := p.NameWithContext(c.Request.Context())
		mem, _ := p.MemoryInfoWithContext(c.Request.Context())
		memRss := mem.RSS
		//mem, _ := p.MemoryInfoWithContext(c.Request.Context())
		cmd = strings.ReplaceAll(cmd, "\"", "\\\"")
		name = strings.ReplaceAll(name, "\"", "\\\"")

		rtnMe.WriteString(fmt.Sprintf(`process_memory_rss_bytes{name="%s",cmdLine="%s", pid="%d",ppid="%d"} %d`+"\n", name, cmd, pid, ppid, memRss))

		if i == topNProcesses-1 {
			break
		}
	}

	c.Data(http.StatusOK, "application/openmetrics-text; version=1.0.0; charset=utf-8", []byte(rtnMe.String()))

}
