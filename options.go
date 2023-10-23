package torgo

import (
	"fmt"
	"strings"
)

type Options struct {
	GeneralOptions *GeneralOptions
	// .. add more options
}

func (o *Options) Serialize() string {
	var serializedLines []string

	serializedLines = append(serializedLines, o.GeneralOptions.Serialize()...)

	return strings.Join(serializedLines, "\n")
}

type GeneralOptions struct {
	SocksPort   int
	Logging     []LogConfig
	RunAsDaemon bool
	// .. add more options
}

func (o *GeneralOptions) Serialize() []string {
	var serializedLines []string

	serializedLines = append(serializedLines, fmt.Sprintf("SocksPort %d", o.SocksPort))

	for _, lc := range o.Logging {
		serializedLines = append(serializedLines, lc.Serialize())
	}

	if o.RunAsDaemon {
		serializedLines = append(serializedLines, "RunAsDaemon 1")
	} else {
		serializedLines = append(serializedLines, "RunAsDaemon 0")
	}

	return serializedLines
}

type LogConfig struct {
	SeverityRange string   // minSeverity[-maxSeverity]
	Destinations  []string // stderr|stdout|syslog or file paths
	Domains       []string // list of logging domains, e.g. general, crypto, net, ...
}

func (lc *LogConfig) Serialize() string {
	var serializedLines []string

	domainString := ""
	if len(lc.Domains) > 0 {
		domainString = "[" + strings.Join(lc.Domains, ",") + "]"
	}

	for _, dest := range lc.Destinations {
		if dest == "stderr" || dest == "stdout" || dest == "syslog" {
			serializedLines = append(serializedLines, fmt.Sprintf("Log %s%s %s", domainString, lc.SeverityRange, dest))
		} else {
			serializedLines = append(serializedLines, fmt.Sprintf("Log %s%s file %s", domainString, lc.SeverityRange, dest))
		}
	}

	return strings.Join(serializedLines, "\n")
}

type SocksPortConfig struct {
	Address        string   // Address to bind to, can be empty
	Port           string   // Port to bind to, can be "auto"
	Flags          []string // Flags for the SocksPort
	IsolationFlags []string // Isolation flags for the SocksPort
}

func (spc *SocksPortConfig) Serialize() string {
	var addressPort string
	if spc.Address != "" {
		addressPort = spc.Address + ":" + spc.Port
	} else {
		addressPort = spc.Port
	}

	flagString := ""
	if len(spc.Flags) > 0 {
		flagString = strings.Join(spc.Flags, " ")
	}

	isolationFlagString := ""
	if len(spc.IsolationFlags) > 0 {
		isolationFlagString = strings.Join(spc.IsolationFlags, " ")
	}

	return fmt.Sprintf("SocksPort %s %s %s", addressPort, flagString, isolationFlagString)
}
