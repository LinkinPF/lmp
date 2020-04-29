package api

import (
	"fmt"
	"lmp/config"
	"lmp/pkg/model"
	"net/http"
	"os/exec"
	"path"
	"strconv"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func init() {
	SetRouterRegister(func(router *RouterGroup) {
		engine := router.Group("/api")

		engine.GET("/ping", Ping)
		engine.POST("/data/collect", do_collect)
	})
}

func do_collect(c *Context) {
	//data := sys.Data{}
	m := fillConfigMessage(c)
	//TODO..

	fmt.Println(m)
	//data.Handle(&m)
	pid, err := strconv.Atoi(m.Pid)
	if err != nil {
		seelog.Error("pid error")
	}

	// For static bpf files
	go execute(pid)
	fmt.Println("start extracting data...")
	seelog.Info("start extracting data...")
	c.Redirect(http.StatusMovedPermanently, "http://"+config.GrafanaIp)
	return
}

func Ping(c *Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func execute(pid int) {
	collector := path.Join(config.DefaultCollectorPath, "collect.py")
	arg1 := "-P"
	arg2 := strconv.Itoa(pid)
	//fmt.Println(collector)
	cmd := exec.Command("python", collector, arg1, arg2)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		seelog.Error(err)
		return
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		seelog.Error(err)
		return
	}
	defer stderr.Close()
	//go func() {
	//	scanner := bufio.NewScanner(stdout)
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		fmt.Println(line)
	//	}
	//}()
	//
	//go func() {
	//	scanner := bufio.NewScanner(stderr)
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		fmt.Println(line)
	//	}
	//}()
	err = cmd.Start()
	if err != nil {
		seelog.Error(err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		seelog.Error(err)
		return
	}
}

func fillConfigMessage(c *Context) model.ConfigMessage {
	var m model.ConfigMessage
	if _, ok := c.GetPostForm("dispatchingdelay"); ok {
		m.DispatchingDelay = true
	} else {
		m.DispatchingDelay = false
	}
	if _, ok := c.GetPostForm("waitingqueuelength"); ok {
		m.WaitingQueueLength = true
	} else {
		m.WaitingQueueLength = false
	}
	if _, ok := c.GetPostForm("softirqtime"); ok {
		m.SoftIrqTime = true
	} else {
		m.SoftIrqTime = false
	}
	if _, ok := c.GetPostForm("hardirqtime"); ok {
		m.HardIrqTime = true
	} else {
		m.HardIrqTime = false
	}
	if _, ok := c.GetPostForm("oncputime"); ok {
		m.OnCpuTime = true
	} else {
		m.OnCpuTime = false
	}
	if pid, ok := c.GetPostForm("pid"); ok {
		m.Pid = pid
	} else {
		m.Pid = "-1"
	}
	return m
}
