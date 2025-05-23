package api

import (
	statsService "github.com/whaleblueio/Xray-core/app/stats/command"
	"github.com/whaleblueio/Xray-core/main/commands/base"
)

var cmdGetStats = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api stats [--server=127.0.0.1:8080] [-name '']",
	Short:       "Get statistics",
	Long: `
Get statistics from Xray.
Arguments:
	-s, -server 
		The API server address. Default 127.0.0.1:8080
	-t, -timeout
		Timeout seconds to call API. Default 3
	-name
		Name of the stat counter.
	-reset
		Reset the counter to fetching its value.
Example:
	{{.Exec}} {{.LongName}} --server=127.0.0.1:8080 -name "inbound>>>statin>>>traffic>>>downlink"
`,
	Run: executeGetStats,
}

func executeGetStats(cmd *base.Command, args []string) {
	setSharedFlags(cmd)
	statName := cmd.Flag.String("name", "", "")
	reset := cmd.Flag.Bool("reset", false, "")
	cmd.Flag.Parse(args)

	conn, ctx, close := dialAPIServer()
	defer close()

	client := statsService.NewStatsServiceClient(conn)
	r := &statsService.GetStatsRequest{
		Name:   *statName,
		Reset_: *reset,
	}
	resp, err := client.GetStats(ctx, r)
	if err != nil {
		base.Fatalf("failed to get stats: %s", err)
	}
	showResponese(resp)
}
