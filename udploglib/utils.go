package udploglib

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	SyslogFacilityKernel = iota
	SyslogFacilityUser
	SyslogFacilityMail
	SyslogFacilitySystemDaemons
	SyslogFacilitySecurityAuth
	SyslogFacilityInternalSyslogd
	SyslogFacilityLinePrinter
	SyslogFacilityNetworkNews
	SyslogFacilityUUCP
	SyslogFacilityClockDaemon
	SyslogFacilitySecurityAuth2
	SyslogFacilityFTP
	SyslogFacilityNTP
	SyslogFacilityLogAudit
	SyslogFacilityLogAlert
	SyslogFacilityClockDaemon2
	SyslogFacilityLocal0
	SyslogFacilityLocal1
	SyslogFacilityLocal2
	SyslogFacilityLocal3
	SyslogFacilityLocal4
	SyslogFacilityLocal5
	SyslogFacilityLocal6
	SyslogFacilityLocal7
)

const (
	SyslogSeverityEmergency = iota
	SyslogSeverityAlert
	SyslogSeverityCritical
	SyslogSeverityError
	SyslogSeverityWarning
	SyslogSeverityNotice
	SyslogSeverityInformational
	SyslogSeverityDebug
)

// SyslogFacilityString is a map containing the textual equivalence of a given facility number
var SyslogFacilityString = map[int]string{
	SyslogFacilityKernel:          "kernel",
	SyslogFacilityUser:            "user",
	SyslogFacilityMail:            "mail",
	SyslogFacilitySystemDaemons:   "system daemons",
	SyslogFacilitySecurityAuth:    "security/auth",
	SyslogFacilityInternalSyslogd: "internal syslogd",
	SyslogFacilityLinePrinter:     "line printer",
	SyslogFacilityNetworkNews:     "network news",
	SyslogFacilityUUCP:            "uucp",
	SyslogFacilityClockDaemon:     "clock daemon",
	SyslogFacilitySecurityAuth2:   "security/auth",
	SyslogFacilityFTP:             "ftp",
	SyslogFacilityNTP:             "ntp",
	SyslogFacilityLogAudit:        "log audit",
	SyslogFacilityLogAlert:        "log alert",
	SyslogFacilityClockDaemon2:    "clock daemon",
	SyslogFacilityLocal0:          "local0",
	SyslogFacilityLocal1:          "local1",
	SyslogFacilityLocal2:          "local2",
	SyslogFacilityLocal3:          "local3",
	SyslogFacilityLocal4:          "local4",
	SyslogFacilityLocal5:          "local5",
	SyslogFacilityLocal6:          "local6",
	SyslogFacilityLocal7:          "local7",
}

// SyslogSeverityString is a map containing the textual equivalence of a given severity number
var SyslogSeverityString = map[int]string{
	SyslogSeverityEmergency:     "emergency",
	SyslogSeverityAlert:         "alert",
	SyslogSeverityCritical:      "critical",
	SyslogSeverityError:         "error",
	SyslogSeverityWarning:       "warning",
	SyslogSeverityNotice:        "notice",
	SyslogSeverityInformational: "informational",
	SyslogSeverityDebug:         "debug",
}

// GetLogItem returns the log entry format, elasticsearch type, message and error (if any)
func GetLogItem(buf []byte) ([]string, error) {

	parts := strings.SplitN(string(buf), ":", 3)
	if len(parts) != 3 {
		return []string{"", "", ""}, fmt.Errorf("Invalid log item")
	}
	if parts[0] != "json" && parts[0] != "plain" {
		return []string{"", "", ""}, fmt.Errorf("Log format %s is invalid", parts[0])
	}
	if parts[1] == "" {
		return []string{"", "", ""}, fmt.Errorf("A log type must be specified")
	}
	if parts[2] == "" {
		return []string{"", "", ""}, fmt.Errorf("Log data is empty")
	}

	return []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])}, nil
}

/*
   The Priority value is calculated by first multiplying the Facility
   number by 8 and then adding the numerical value of the Severity.

	 Source: https://tools.ietf.org/html/rfc5424 [Page 10]
*/

// ExtractSyslogFacilityAndServerity returns the facility and severity of a valid syslog message
func getSyslogFacilityAndServerity(syslogMsg string) (string, int, int, error) {

	re := regexp.MustCompile(`^<([0-9]{1,3})>(.*)`)
	matches := re.FindStringSubmatch(syslogMsg)
	if len(matches) < 3 {
		return "", 0, 0, errors.New("Could not extract syslog priority from message")
	}
	priorityNum, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", 0, 0, nil
	}
	severity := int(math.Mod(float64(priorityNum), 8.0))
	facility := (priorityNum - severity) / 8
	return matches[2], facility, severity, nil

}

// GetSyslogFacilityAndServerityAsStr returns the facility and serverity as text
func GetSyslogMsgDetails(syslogMsg string) (string, string, string, error) {
	msg, facility, severity, err := getSyslogFacilityAndServerity(syslogMsg)
	if err != nil {
		return "", "", "", err
	}
	return msg, SyslogFacilityString[facility], SyslogSeverityString[severity], nil
}
