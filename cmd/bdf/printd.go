package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/server"
	"github.com/shibukawa/bdf/print/spool"
)

// printdCmd runs the print server in the foreground until SIGINT or
// SIGTERM (Ctrl-C and Ctrl-Break on Windows): a service manager (systemd,
// launchd, NSSM, WinSW) keeps it running.
func printdCmd(args []string) {
	fs := flag.NewFlagSet("printd", flag.ExitOnError)
	config := fs.String("config", "", "the configuration file (YAML)")
	checkOnly := fs.Bool("check", false, "check the configuration, ask each printer what it supports, and exit; the spool is not opened")
	compact := fs.Bool("compact", false, "rewrite the spool's file without the space of removed jobs, and exit; the server must be stopped")
	listen := fs.String("listen", "", "the address to listen on, over the configuration's")
	profile := fs.String("airprint-profile", "", "write a configuration profile (.mobileconfig) that adds the IPP printers to iPhones and iPads at a `host` name or address, with the server's certificate, and exit")
	zone := fs.String("dns-sd", "", "write the DNS-SD records of the IPP printers for a DNS `domain` (wide-area DNS-SD: clients of other subnets find them), and exit")
	logFormat := fs.String("log", "json", "the format of the log on the standard error: json or text")
	debug := fs.Bool("debug", false, "log what is done step by step")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: bdf printd -config printd.yaml [flags]
  runs the print server in the foreground: a spool of jobs, its queues and
  printers, and the HTTP API (/api/v1) that bdf print -server and bdf jobs use;
  SIGINT or SIGTERM stops it after the documents being sent`)
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if *config == "" || fs.NArg() != 0 {
		fs.Usage()
		os.Exit(2)
	}
	level := slog.LevelInfo
	if *debug {
		level = slog.LevelDebug
	}
	var h slog.Handler
	switch *logFormat {
	case "json":
		h = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	case "text":
		h = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	default:
		check(fmt.Errorf("-log %q: json or text", *logFormat))
	}
	log := slog.New(h)

	cfg, err := server.LoadConfig(*config)
	check(err)
	if *listen != "" {
		cfg.Listen = *listen
	}
	switch {
	case *profile != "":
		if cfg.IPP == nil {
			check(fmt.Errorf("%s has no ipp: section", *config))
		}
		_, port, err := net.SplitHostPort(cfg.IPP.Listen)
		check(err)
		n, err := strconv.Atoi(port)
		check(err)
		b, err := server.AirPrintProfile(cfg, *profile, n)
		check(err)
		os.Stdout.Write(b)
	case *zone != "":
		if cfg.IPP == nil {
			check(fmt.Errorf("%s has no ipp: section", *config))
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		caps := map[string]*ipp.Printer{}
		for _, r := range server.CheckPrinters(ctx, cfg) {
			if r.Err != nil {
				fmt.Fprintf(os.Stderr, "bdf printd: printer %s: %v (its capabilities are left out)\n", r.Name, r.Err)
				continue
			}
			caps[r.Name] = r.Printer
		}
		_, port, err := net.SplitHostPort(cfg.IPP.Listen)
		check(err)
		n, err := strconv.Atoi(port)
		check(err)
		out, err := server.ZoneRecords(cfg, *zone, n, caps)
		check(err)
		fmt.Print(out)
	case *compact:
		before, after, err := spool.Compact(cfg.Spool)
		check(err)
		fmt.Printf("%s: %d → %d bytes\n", cfg.Spool, before, after)
	case *checkOnly:
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		failed := false
		tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "PRINTER\tURL\tSTATE\tMODEL\tFORMATS")
		for _, r := range server.CheckPrinters(ctx, cfg) {
			if r.Err != nil {
				failed = true
				fmt.Fprintf(tw, "%s\t%s\terror: %v\t\t\n", r.Name, r.URL, r.Err)
				continue
			}
			p := r.Printer
			state := p.State.String()
			if !p.AcceptingJobs {
				state += ", not accepting"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", r.Name, r.URL, state, p.MakeModel, strings.Join(p.Formats, ","))
		}
		tw.Flush()
		if failed {
			os.Exit(1)
		}
	default:
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		conv := printConverter{fontDirs: cfg.Fonts.Dirs, noSystem: cfg.Fonts.NoSystem}
		srv, err := server.New(ctx, cfg, server.Options{Log: log, Converter: conv})
		check(err)
		if err := srv.Run(ctx); err != nil {
			log.Error("printd", "err", err)
			os.Exit(1)
		}
	}
}
