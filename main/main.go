package main

import (
	"flag"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/whaleblueio/Xray-core/main/commands/base"
	_ "github.com/whaleblueio/Xray-core/main/distro/all"
)

var LogLevel = os.Getenv("LOG_LEVEL")

func main() {
	os.Args = getArgsV4Compatible()
	initLogger()
	base.RootCommand.Long = "Xray is a platform for building proxies."
	base.RootCommand.Commands = append(
		[]*base.Command{
			cmdRun,
			cmdVersion,
		},
		base.RootCommand.Commands...,
	)
	base.Execute()
}
func initLogger() {
	logger.SetOutput(os.Stdout)
	var level = logger.InfoLevel
	if strings.Compare(LogLevel, "Warn") == 0 {
		level = logger.WarnLevel
	} else if strings.Compare(LogLevel, "Info") == 0 {
		level = logger.InfoLevel
	} else if strings.Compare(LogLevel, "Debug") == 0 {
		level = logger.DebugLevel
	} else if strings.Compare(LogLevel, "Trace") == 0 {
		level = logger.TraceLevel
	} else {
		level = logger.InfoLevel
	}
	logger.SetLevel(level)
	logger.Infof("Log level :%s", level)
}
func getArgsV4Compatible() []string {
	if len(os.Args) == 1 {
		return []string{os.Args[0], "run"}
	}
	if os.Args[1][0] != '-' {
		return os.Args
	}
	version := false
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&version, "version", false, "")
	// parse silently, no usage, no error output
	fs.Usage = func() {}
	fs.SetOutput(&null{})
	err := fs.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		//fmt.Println("DEPRECATED: -h, WILL BE REMOVED IN V5.")
		//fmt.Println("PLEASE USE: xray help")
		//fmt.Println()
		return []string{os.Args[0], "help"}
	}
	if version {
		//fmt.Println("DEPRECATED: -version, WILL BE REMOVED IN V5.")
		//fmt.Println("PLEASE USE: xray version")
		//fmt.Println()
		return []string{os.Args[0], "version"}
	}
	//fmt.Println("COMPATIBLE MODE, DEPRECATED.")
	//fmt.Println("PLEASE USE: xray run [arguments] INSTEAD.")
	//fmt.Println()
	return append([]string{os.Args[0], "run"}, os.Args[1:]...)
}

type null struct{}

func (n *null) Write(p []byte) (int, error) {
	return len(p), nil
}
