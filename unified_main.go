package main

import (
	unified "bitbucket.org/ecosse-hosting/unified/lib"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	ctx = context.TODO()
)

type SSHConnection struct {
	Username string
	Host     string
	Port     int
}

var ssh_port int
var ssh_host string
var ssh_user string
var table_output bool
var json_output bool

func main() {
	app := cli.NewApp()

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: time.Second * 300, Transport: tr}
	_, err := client.Get("https://golang.org/")
	if err != nil {
		fmt.Println(err)
	}

	d := &unified.UnifiedDBOptions{
		DbUsageEnabled: true,
		UseInMemoryDB:  true,
	}

	o := &unified.UnifiedOptions{
		DbUsage: d,
	}

	cx := unified.NewUniFiClient(client, o)

	cx.Authentication.Login(ctx, "donstewa", "Blackadder!782")

	app.Name = "Unified"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Don Stewart",
			Email: "don.stewart at icloud.com",
		},
	}
	app.Copyright = "(c) 2017 Ecosse Hosting"
	app.HelpName = "Unified - Ubiquiti UniFi Tools"
	app.Usage = "unified help"
	app.UsageText = "Unified - demonstrating the available API"


	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:  "UseDB, db, b",
			Usage: "Enable the use of a DB for storing data from the UniFi Controller.",
		},
		cli.BoolTFlag{
			Name:  "UseCache, c",
			Usage: "Enable the use of the internal DB for retreiving data. Implies -db",
		},
		cli.BoolTFlag{
			Name:  "Daemon, d",
			Usage: "Runs Unified as a daemon.",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "db",
			Aliases:     []string{"database"},
			Category:    "Database",
			Usage:       "do the doo",
			UsageText:   "doo - does the dooing",
			Description: "Manages the Unified DB if enabled.",
			ArgsUsage:   "[arrgh]",
			Flags: []cli.Flag{
				cli.BoolTFlag{Name: "db, UseDB"},
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "clean",
					Usage: "Drops the selected stored data returning the DB to an empty state.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
					Subcommands: cli.Commands{
						cli.Command{
							Name:  "all",
							Usage: "Drops all the currently stored data returning the DB to an empty state.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:  "alarms",
							Usage: "Drops all the currently stored alarm data only.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:  "events",
							Usage: "Drops all the currently stored event data only.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:  "users",
							Usage: "Drops all the currently stored user data only.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
					},
				},
			},
			SkipFlagParsing: false,
			HideHelp:        false,
			Hidden:          false,
			HelpName:        "database",
			BashComplete: func(c *cli.Context) {
				fmt.Fprintf(c.App.Writer, "--better\n")
			},
			Before: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "brace for impact\n")
				return nil
			},
			After: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
				return nil
			},
			Action: func(c *cli.Context) error {
				c.Command.FullName()
				c.Command.HasName("wop")
				c.Command.Names()
				c.Command.VisibleFlags()
				fmt.Fprintf(c.App.Writer, "dodododododoodododddooooododododooo\n")
				if c.Bool("forever") {
					c.Command.Run(c)
				}
				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "for shame\n")
				return err
			},
		},
		cli.Command{
			Name:     "devices",
			Category: "Devices",
			Usage:    "All devices command.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:    "ls",
					Aliases: []string{"list"},
					Usage:   "Displays a list of known UniFi devices (of all types).",
					Flags: []cli.Flag{
						cli.BoolTFlag{
							Name:        "table, T",
							Usage:       "Displays device short data in a table on the console.",
							Destination: &table_output,
						},
						cli.BoolFlag{
							Name:        "json, J",
							Usage:       "Displays device short data in JSON on the console.",
							Destination: &json_output,
						},
					},
					Action: func(c *cli.Context) error {
						fmt.Println("\nunified devices ls\n")
						devices, _, err := cx.Devices.ListShort(ctx, "all", nil)
						if err != nil {
							return nil
						}
						if json_output {
							table_output = false
							enc := json.NewEncoder(os.Stdout)
							enc.SetIndent("", "    ")
							enc.Encode(devices)
						}
						if table_output {
							outputDevicesToTable(devices)
						}
						return nil
					},
				},
				cli.Command{
					Name:      "inspect",
					Aliases:   []string{"info"},
					Usage:     "View detailed configuration of a running Unifi Device.",
					ArgsUsage: "[mac_address]",
					Action: func(c *cli.Context) error {
						fmt.Println("\nunified devices inspect `mac_address`\n")
						device, _, err := cx.Devices.GetByMac(ctx, c.Args().Get(0))
						if err != nil {
							return nil
						}
						enc := json.NewEncoder(os.Stdout)
						enc.SetIndent("", "    ")
						enc.Encode(device)

						return nil
					},
				},
				cli.Command{
					Name:     "ugw",
					Aliases:  []string{"gateway, usg"},
					Category: "Device Types (sub commands)",
					Usage:    "Commands relating to a UniFi Security Gateway aka USG.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
					Subcommands: cli.Commands{
						cli.Command{
							Name:    "ls",
							Aliases: []string{"list"},
							Usage:   "Displays a list of known UniFi USGs.",
							Flags: []cli.Flag{
								cli.BoolTFlag{
									Name:        "table, T",
									Usage:       "Displays device short data in a table on the console.",
									Destination: &table_output,
								},
								cli.BoolFlag{
									Name:        "json, J",
									Usage:       "Displays device short data in JSON on the console.",
									Destination: &json_output,
								},
							},
							Action: func(c *cli.Context) error {
								fmt.Println("\nunified devices ugw ls\n")
								devices, _, err := cx.Devices.ListShort(ctx, "ugw", nil)
								if err != nil {
									return nil
								}
								if json_output {
									table_output = false
									enc := json.NewEncoder(os.Stdout)
									enc.SetIndent("", "    ")
									enc.Encode(devices)
								}
								if table_output {
									outputDevicesToTable(devices)
								}
/*								table := tablewriter.NewWriter(os.Stdout)
								for _, v := range devices {
									fieldNames := structs.Names(&v)
									table.SetHeader(fieldNames)

									fieldValues := structs.Values(&v)
									valuesArray := make([]string, len(fieldValues))
									for k, w := range fieldValues {
										switch x := w.(type) {
										case string:
											valuesArray[k] = x
										case bool:
											valuesArray[k] = strconv.FormatBool(x)
										case int:
											valuesArray[k] = strconv.Itoa(x)
										}
									}
									table.Append(valuesArray)
								}
								table.Render() // Send output*/
								return nil
							},
						},
						cli.Command{
							Name:  "ps",
							Usage: "Information about a running Unifi USG.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "inspect",
							Aliases: []string{"info"},
							Usage:   "View configuration of a running Unifi USG.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
					},
				},
				cli.Command{
					Name:     "uap",
					Aliases:  []string{"accesspoint, ap"},
					Category: "Device Types (sub commands)",
					Usage:    "Commands relating to a UniFi Access Points or APs.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
					Subcommands: cli.Commands{
						cli.Command{
							Name:    "ls",
							Aliases: []string{"list"},
							Usage:   "Displays a list of known Unifi APs.",
							Flags: []cli.Flag{
								cli.BoolTFlag{
									Name:        "table, T",
									Usage:       "Displays device short data in a table on the console.",
									Destination: &table_output,
								},
								cli.BoolFlag{
									Name:        "json, J",
									Usage:       "Displays device short data in JSON on the console.",
									Destination: &json_output,
								},
							},
							Action: func(c *cli.Context) error {
								fmt.Println("\nunified devices uap ls\n")
								devices, _, err := cx.Devices.ListShort(ctx, "uap", nil)
								if err != nil {
									return nil
								}
								if json_output {
									table_output = false
									enc := json.NewEncoder(os.Stdout)
									enc.SetIndent("", "    ")
									enc.Encode(devices)
								}
								if table_output {
									outputDevicesToTable(devices)
								}
/*								table := tablewriter.NewWriter(os.Stdout)
								for _, v := range devices {
									fieldNames := structs.Names(&v)
									table.SetHeader(fieldNames)

									fieldValues := structs.Values(&v)
									valuesArray := make([]string, len(fieldValues))
									for k, w := range fieldValues {
										switch x := w.(type) {
										case string:
											valuesArray[k] = x
										case bool:
											valuesArray[k] = strconv.FormatBool(x)
										case int:
											valuesArray[k] = strconv.Itoa(x)
										}
									}
									table.Append(valuesArray)
								}
								table.Render() // Send output*/
								return nil
							},
						},
						cli.Command{
							Name:  "ps",
							Usage: "Information about a running Unifi AP.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "inspect",
							Aliases: []string{"info"},
							Usage:   "View configuration of a running AP.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
					},
				},
				cli.Command{
					Name:     "usw",
					Aliases:  []string{"switch"},
					Category: "Device Types (sub commands)",
					Usage:    "Commands relating to a UniFi Switch.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
					Subcommands: cli.Commands{
						cli.Command{
							Name:    "ls",
							Aliases: []string{"list"},
							Usage:   "Displays a list of known Unifi Switches.",
							Flags: []cli.Flag{
								cli.BoolTFlag{
									Name:        "table, T",
									Usage:       "Displays device short data in a table on the console.",
									Destination: &table_output,
								},
								cli.BoolFlag{
									Name:        "json, J",
									Usage:       "Displays device short data in JSON on the console.",
									Destination: &json_output,
								},
							},
							Action: func(c *cli.Context) error {
								fmt.Println("\nunified devices usw ls\n")
								devices, _, err := cx.Devices.ListShort(ctx, "usw", nil)
								if err != nil {
									return nil
								}
								if json_output {
									table_output = false
									enc := json.NewEncoder(os.Stdout)
									enc.SetIndent("", "    ")
									enc.Encode(devices)
								}
								if table_output {
									outputDevicesToTable(devices)
								}
/*								table := tablewriter.NewWriter(os.Stdout)
								for _, v := range devices {
									fieldNames := structs.Names(&v)
									table.SetHeader(fieldNames)

									fieldValues := structs.Values(&v)
									valuesArray := make([]string, len(fieldValues))
									for k, w := range fieldValues {
										switch x := w.(type) {
										case string:
											valuesArray[k] = x
										case bool:
											valuesArray[k] = strconv.FormatBool(x)
										case int:
											valuesArray[k] = strconv.Itoa(x)
										}
									}
									table.Append(valuesArray)
								}
								table.Render() // Send output*/
								return nil
							},
						},
						cli.Command{
							Name:  "ps",
							Usage: "Information about a running Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:      "inspect",
							Aliases:   []string{"info"},
							ArgsUsage: "[mac_address]",
							Usage:     "View configuration of a running Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "cn",
							Aliases: []string{"confignetwork"},
							Usage:   "View configuration for Config Network for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "pt",
							Aliases: []string{"porttable"},
							Usage:   "View configuration for Port Table for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "po",
							Aliases: []string{"portoverride"},
							Usage:   "View configuration for Port Overrides for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "dt",
							Aliases: []string{"down, downlink"},
							Usage:   "View configuration for Downlink Table for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "up",
							Aliases: []string{"uplink"},
							Usage:   "View configuration for Uplink Table for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "eth",
							Aliases: []string{"ethernet"},
							Usage:   "View configuration for Ethernet Table for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "ss",
							Aliases: []string{"sysstats"},
							Usage:   "View configuration for System Statistics for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "stats",
							Aliases: []string{"switchstats, portstats"},
							Usage:   "View configuration for Switch Statistics for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
						cli.Command{
							Name:    "ssh",
							Aliases: []string{"sshsessions, consoles"},
							Usage:   "View configuration for SSH Session Table for a Unifi Switch.",
							Action: func(c *cli.Context) error {
								fmt.Println("new task template: ", c.Args().First())
								return nil
							},
						},
					},
				},
			},
		},

		cli.Command{
			Name:     "site",
			Category: "Site",
			Usage:    "Commands relating to a Site.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:    "ls",
					Aliases: []string{"list"},
					Usage:   "Displays a list of known Sites.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				cli.Command{
					Name:  "ps",
					Usage: "Information about a running Siteh.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				cli.Command{
					Name:  "inspect",
					Usage: "View configuration of a running Site.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		cli.Command{
			Name:     "operator",
			Category: "Site",
			Usage:    "Commands relating to an Operator.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:    "ls",
					Aliases: []string{"list"},
					Usage:   "Displays a list of known Operators.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				cli.Command{
					Name:  "inspect",
					Usage: "View configuration of a running Site.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		cli.Command{
			Name:     "exec",
			Category: "Remote",
			Usage:    "Executes a command on a remote device via ssh.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "user, U",
					Usage:       " Connects via SSH using the `USERNAME` supplied.",
					Destination: &ssh_user,
				},
				cli.StringFlag{
					Name:        "mac, M",
					Usage:       " Connects via SSH to the specified `MAC_ADDRESS` device.",
					Destination: &ssh_host,
				},
				cli.IntFlag{
					Name:        "port, P",
					Usage:       " Connects via SSH to the specified `PORT` supplied.",
					Value:       22,
					Destination: &ssh_port,
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					fmt.Println("exec - Not enough Args passed")
					return nil
				}
				ip_address, err := cx.Devices.GetIPFromMac(ctx, c.Args().Get(0))
				if err != nil {
					return nil
				}
				_, session, err :=
					unified.ConnectToSSHHost(ssh_user, ip_address+":"+strconv.Itoa(ssh_port))
				if err != nil {
					panic(err)
				}
				fmt.Println("Established SSH Connection as " + ssh_user)
				out, err := session.CombinedOutput(os.Args[3])
				if err != nil {
					panic(err)
				}
				fmt.Println(string(out))
				return nil
			},
		},
		cli.Command{
			Name:     "unknown",
			Category: "Hotspot",
			Usage:    "Stuff to do with Hotspots.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
		},
		cli.Command{
			Name:     "net",
			Aliases:  []string{"network, lan"},
			Category: "Network",
			Usage:    "Stuff to do with Networks.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
		},
		cli.Command{
			Name:     "wireless",
			Aliases:  []string{"wlan"},
			Category: "Network",
			Usage:    "Stuff to do with Networks.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
		},
		cli.Command{
			Name:     "routing",
			Category: "Routing & Firewall",
			Usage:    "Stuff to do with the Routing.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "ls",
					Usage: "List current routing.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		cli.Command{
			Name:     "firewall",
			Category: "Routing & Firewall",
			Usage:    "Stuff to do with the Firewall.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "ls",
					Usage: "List current routing.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		cli.Command{
			Name:     "dpi",
			Category: "Metrics",
			Usage:    "Deep Packet Inspection.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "clear",
					Usage: "Clears current DPI counters.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				cli.Command{
					Name:  "enable",
					Usage: "Enables DPI counters.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				cli.Command{
					Name:  "disable",
					Usage: "Diables DPI counters.",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		cli.Command{
			Name:     "controller",
			Category: "Controller",
			Usage:    "Controller management.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "alarms",
					Usage: "Alarm based commands.",
					Action: func(c *cli.Context) error {
						color.Green("unified controller alarms help ")
						return nil
					},
					Subcommands: cli.Commands{
						cli.Command{
							Name:    "ls",
							Aliases: []string{"list"},
							Usage:   "Displays current & archived alarms.",
							Action: func(c *cli.Context) error {
								alarms, _, err := cx.Alarms.List(ctx, nil)
								if err != nil {
									return nil
								}

								fmt.Println("\n")
								for _, v := range alarms {
									table := tablewriter.NewWriter(os.Stdout)
									fieldNames := structs.Names(&v)
									table.SetHeader(fieldNames)

									fieldValues := structs.Values(&v)
									valuesArray := make([]string, len(fieldValues))
									for k, w := range fieldValues {
										switch x := w.(type) {
										case string:
											valuesArray[k] = x
										case bool:
											valuesArray[k] = strconv.FormatBool(x)
										case int:
											valuesArray[k] = strconv.Itoa(x)
										}
									}
									table.Append(valuesArray)
									table.Render() // Send output
								}

								return nil
							},
						},
					},
				},
				cli.Command{
					Name:  "events",
					Usage: "Display controller events.",
					Action: func(c *cli.Context) error {
						color.Green("unified controller events help ")
						return nil
					},
					Subcommands: cli.Commands{
						cli.Command{
							Name:    "ls",
							Aliases: []string{"list"},
							Usage:   "Displays current & archived events.",
							Action: func(c *cli.Context) error {
								events, _, err := cx.Events.List(ctx, nil)
								if err != nil {
									return nil
								}

								table := tablewriter.NewWriter(os.Stdout)
								//table.SetAutoMergeCells(true)
								//table.SetRowLine(true)
								var eventsArray = make([][]string, len(events))
								for j, v := range events {
									//table := tablewriter.NewWriter(os.Stdout)
									table.SetHeader(structs.Names(&v))
									fieldValues := structs.Values(&v)
									eventsArray[j] = make([]string, len(fieldValues))
									for k, w := range fieldValues {
										switch x := w.(type) {
										case string:
											eventsArray[j][k] = x
										case bool:
											eventsArray[j][k] = strconv.FormatBool(x)
										case int:
											eventsArray[j][k] = strconv.Itoa(x)
										}
									}
									//table.Append(eventsArray[j])
									//table.Render() // Send output
								}
								table.AppendBulk(eventsArray)
								table.Render() // Send output
								return nil
							},
						},
					},
				},
			},
		},
		cli.Command{
			Name:     "guest",
			Category: "Users",
			Usage:    "Guest management.",
			Action: func(c *cli.Context) error {
				fmt.Println("new task template: ", c.Args().First())
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("unified help - for more help information")
		return nil
	}

	app.EnableBashCompletion = true
	app.BashComplete = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "lipstick\nkiss\nme\nlipstick\nringo\n")
	}

	app.Run(os.Args)
}
func outputDevicesToTable(devices []unified.DeviceShort) {
	table := tablewriter.NewWriter(os.Stdout)
	for _, v := range devices {
		fieldNames := structs.Names(&v)
		table.SetHeader(fieldNames)

		fieldValues := structs.Values(&v)
		valuesArray := make([]string, len(fieldValues))
		for k, w := range fieldValues {
			switch x := w.(type) {
			case string:
				valuesArray[k] = x
			case bool:
				valuesArray[k] = strconv.FormatBool(x)
			case int:
				valuesArray[k] = strconv.Itoa(x)
			}
		}
		table.Append(valuesArray)
	}
	table.Render()
	// Send output
}

func wopAction(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, ":wave: over here, eh\n")
	return nil
}

func setup() {
	client := unified.NewUniFiClient(nil, nil)
	value, err := url.Parse("https://192.168.10.7:8443/api/s/")

	if err != nil {
		//log.Fatal(err)
	} else {
		client.BaseURL = value
	}
}
