package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"strings"
)

func ProcessHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-type", "text/plain")

	rtnMe := strings.Builder{}

	processes, err := process.Processes()

	if err != nil {
		return
	}

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
		memUsed := mem.RSS
		//mem, _ := p.MemoryInfoWithContext(c.Request.Context())

		rtnMe.WriteString(fmt.Sprintf(`process_cpu_percent{name="%s",cmdLine="%s",pid="%d",ppid="%d"} %f`+"\n", name, cmd, pid, ppid, cpuPercent/100))
		rtnMe.WriteString(fmt.Sprintf(`process_memory_bytes{name="%s",cmdLine="%s", pid="%d",ppid="%d"} %d`+"\n", name, cmd, pid, ppid, memUsed))

	}

	c.Data(http.StatusOK, "text/plain", []byte(rtnMe.String()))

}
