package main

import (
	unified "bitbucket.org/ecosse-hosting/unified/lib"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"github.com/jawher/mow.cli"
)

var (
	ctx = context.TODO()
)

type SSHConnection struct {
	Username string
	Host     string
	Port     int
}

var ssh_portOption bool
var ssh_hostOption bool
var ssh_userOption bool
var table_output bool
var json_output bool
var username bool
var password bool
var useDBOption bool
var daemon bool

func main() {
	app := cli.App("unified", "Unified CLI for Ubiquiti UniFi")
	app.Version("v version", "unified 0.0.1")
	//app.Spec = "-u -p [COMMAND [ARG...]]"

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

	var (
		/*
		useDB = app.Bool(
			cli.BoolOpt{
				Name:      "db useDB",
				Value:     true,
				Desc:      "Enable the use of a DB for storing data from the UniFi Controller.",
				EnvVar:    "UNIFIED_USEDB",
				SetByUser: &useDBOption,
			},
		)

		useCache = app.Bool(
			cli.BoolOpt{
				Name:      "c useCache",
				Value:     true,
				Desc:      "Enable the use of the internal DB for retreiving data. Implies -db.",
				EnvVar:    "UNIFIED_USE_CACHE",
				SetByUser: &useDBOption,
			},
		)

		daemon = app.Bool(
			cli.BoolOpt{
				Name:      "d daemon",
				Value:     false,
				Desc:      "Runs Unified as a daemon.",
				EnvVar:    "UNIFIED_DAEMON",
				SetByUser: &daemon,
			},
		)
		*/

		user = app.String(
			cli.StringOpt{
				Name:      "u username",
				Value:     "",
				Desc:      "Runs Unified as a daemon.",
				EnvVar:    "UNIFIED_USER",
				SetByUser: &username,
			},
		)

		pass = app.String(cli.StringOpt{
			Name:      "p password",
			Desc:      "Runs Unified as a daemon.",
			EnvVar:    "UNIFIED_PASSWORD",
			SetByUser: &password,
		},
		)
	)

	app.Before = func() {
		cx.Authentication.Login(ctx, *user, *pass)
	}


	app.Command("db", "Manages the Unified DB if enabled.", func(cmd *cli.Cmd) {
		cmd.Command(
			"clean",
			"Drops the selected stored data returning the DB to an empty state.",
			func(cmd2 *cli.Cmd) {
				cmd2.Command(
					"all",
					"Drops all the currently stored data returning the DB to an empty state.",
					func(*cli.Cmd) {
						cmd2.Action = func() {
							fmt.Println("new task template: ")
						}
					})
				cmd2.Command(
					"alarms",
					"Drops all the currently stored alarm data only.",
					func(*cli.Cmd) {
						cmd2.Action = func() {
							fmt.Println("new task template: ")
						}
					})
				cmd2.Command(
					"events",
					"Drops all the currently stored event data only.",
					func(*cli.Cmd) {
						cmd2.Action = func() {
							fmt.Println("new task template: ")
						}
					})
				cmd2.Command(
					"users",
					"Drops all the currently stored user data only.",
					func(*cli.Cmd) {
						cmd2.Action = func() {
							fmt.Println("new task template: ")
						}
					})

			})
	})

	app.Command("devices", "All devices command.", func(cmd *cli.Cmd) {
		cmd.Command(
			"ls",
			"Displays a list of known UniFi devices (of all types).",
			func(cmd2 *cli.Cmd) {
				cmd2.Spec = "[-t | -j]"
				tableo := cmd2.Bool(cli.BoolOpt{
					Name:      "t table",
					Value: 	   true,
					Desc:      "Displays device short data in a table on the console.",
					SetByUser: &table_output,
				})
				jsono := cmd2.Bool(cli.BoolOpt{
					Name:      "j json",
					Desc:      "Displays device short data in JSON on the console.",
					SetByUser: &json_output,
				})
				cmd2.Action = func() {
					fmt.Println("\nunified devices ls\n")
					devices, _, err := cx.Devices.ListShort(ctx, "all", nil)
					if err != nil {

					}
					if *jsono {
						devices = DeviceArrayToJSON(devices)
					} else {
						if *tableo {
							outputDevicesToTable(devices)
						}
					}

				}
			})
		cmd.Command(
			"inspect",
			"View detailed configuration of a running Unifi Device.",
			func(cmd2 *cli.Cmd) {
				macAddress := cmd2.StringArg("MAC_ADDRESS", "", "The MAC address of the device to inspect.")
				cmd2.Action = func() {
					fmt.Println("\nunified devices inspect\n")
					device, _, err := cx.Devices.GetByMac(ctx, *macAddress)
					if err != nil {

					}
					DeviceToJSON(device)
				}
			})
		cmd.Command(
			"ugw",
			"Commands relating to a UniFi Security Gateway (UGW) aka USG.",
			func(cmd2 *cli.Cmd) {
				cmd2.Command(
					"ls",
					"Displays a list of known UniFi USGs.",
					func(cmd3 *cli.Cmd) {
						cmd3.Spec = "[-t | -j]"
						tableo := cmd3.Bool(cli.BoolOpt{
							Name:      "t table",
							Value: 	   true,
							Desc:      "Displays device short data in a table on the console.",
							SetByUser: &table_output,
						})
						jsono := cmd3.Bool(cli.BoolOpt{
							Name:      "j json",
							Desc:      "Displays device short data in JSON on the console.",
							SetByUser: &json_output,
						})
						cmd3.Action = func() {
							fmt.Println("\nunified devices ugw ls\n")
							devices, _, err := cx.Devices.ListShort(ctx, "ugw", nil)
							if err != nil {

							}
							if *jsono {
								devices = DeviceArrayToJSON(devices)
							} else {
								if *tableo {
									outputDevicesToTable(devices)
								}
							}

						}
					})
				cmd2.Command(
					"inspect",
					"View configuration of a running Unifi USG.",
					func(cmd3 *cli.Cmd) {
						macAddress := cmd3.StringArg("MAC_ADDRESS", "",
							"The MAC address of the device to inspect.")
						cmd3.Action = func() {
							fmt.Println("\nunified devices ugw inspect MAC_ADDRESS\n")
							device, _, err := cx.Devices.GetByMac(ctx, *macAddress)
							if err != nil {

							}
							DeviceToJSON(device)
						}
					})
			})
		cmd.Command(
			"uap",
			"Commands relating to a UniFi Access Point (UAP).",
			func(cmd2 *cli.Cmd) {
				cmd2.Command(
					"ls",
					"Displays a list of known UniFi UAPs.",
					func(cmd3 *cli.Cmd) {
						cmd3.Spec = "[-t | -j]"
						tableo := cmd3.Bool(cli.BoolOpt{
							Name:      "t table",
							Value: 	   true,
							Desc:      "Displays device short data in a table on the console.",
							SetByUser: &table_output,
						})
						jsono := cmd3.Bool(cli.BoolOpt{
							Name:      "j json",
							Desc:      "Displays device short data in JSON on the console.",
							SetByUser: &json_output,
						})
						cmd3.Action = func() {
							fmt.Println("\nunified devices uap ls\n")
							devices, _, err := cx.Devices.ListShort(ctx, "uap", nil)
							if err != nil {

							}
							if *jsono {
								devices = DeviceArrayToJSON(devices)
							} else {
								if *tableo {
									outputDevicesToTable(devices)
								}
							}

						}
					})
				cmd2.Command(
					"inspect",
					"View configuration of a running Unifi UAP.",
					func(cmd3 *cli.Cmd) {
						macAddress := cmd3.StringArg("MAC_ADDRESS", "",
							"The MAC address of the device to inspect.")
						cmd3.Action = func() {
							fmt.Println("\nunified devices ugw inspect MAC_ADDRESS\n")
							device, _, err := cx.Devices.GetByMac(ctx, *macAddress)
							if err != nil {

							}
							DeviceToJSON(device)
						}
					})
			})
		cmd.Command(
			"usw",
			"Commands relating to a UniFi Switch (USW).",
			func(cmd2 *cli.Cmd) {
				cmd2.Command(
					"ls",
					"Displays a list of known UniFi USWs.",
					func(cmd3 *cli.Cmd) {
						cmd3.Spec = "[-t | -j]"
						tableo := cmd3.Bool(cli.BoolOpt{
							Name:      "t table",
							Value: 	   true,
							Desc:      "Displays device short data in a table on the console.",
							SetByUser: &table_output,
						})
						jsono := cmd3.Bool(cli.BoolOpt{
							Name:      "j json",
							Desc:      "Displays device short data in JSON on the console.",
							SetByUser: &json_output,
						})
						cmd3.Action = func() {
							fmt.Println("\nunified devices usw ls\n")
							devices, _, err := cx.Devices.ListShort(ctx, "usw", nil)
							if err != nil {

							}
							if *jsono {
								devices = DeviceArrayToJSON(devices)
							} else {
								if *tableo {
									outputDevicesToTable(devices)
								}
							}

						}
					})
				cmd2.Command(
					"inspect",
					"View configuration of a running Unifi USW.",
					func(cmd3 *cli.Cmd) {
						macAddress := cmd3.StringArg("MAC_ADDRESS", "",
							"The MAC address of the device to inspect.")
						cmd3.Action = func() {
							fmt.Println("\nunified devices usw inspect MAC_ADDRESS\n")
							device, _, err := cx.Devices.GetByMac(ctx, *macAddress)
							if err != nil {

							}
							DeviceToJSON(device)
						}
					})
			})
	})

	app.Command("exec", "Open a remote SSH Shell.", func(cmd *cli.Cmd) {
		cmd.Spec = "MAC_ADDRESS (-U [-P])"
		macAddress := cmd.StringArg("MAC_ADDRESS", "",
			"The MAC address of the device to ssh to.")
		ssh_user := cmd.String(cli.StringOpt{
			Name: "U sshuser",
			Desc: "Connects via SSH using the -U SSH_USERNAME supplied.",
			EnvVar:    "UNIFIED_SSH_USERNAME",
			SetByUser: &ssh_userOption,
		})
		/*
		cmd.String(cli.StringOpt{
			Name: "P sshpass",
			Desc: "Connects via SSH using the -P SSH_PASSWORD supplied.",
			EnvVar:    "UNIFIED_SSH_PASSWORD",
			SetByUser: &sshpass,
		})
		*/
		ssh_port := cmd.Int(cli.IntOpt{
			Name: "P ssh_port",
			Desc: "Connects via SSH using the supplied port or defaults to 22.",
			Value: 22,
			EnvVar:    "UNIFIED_SSH_PORT",
			SetByUser: &ssh_portOption,
		})
		cmd.Action = func() {
			ip_address, err := cx.Devices.GetIPFromMac(ctx, *macAddress)
			if err != nil {

			}
			_, session, err :=
				unified.ConnectToSSHHost(*ssh_user, ip_address+":"+strconv.Itoa(*ssh_port))
			if err != nil {
				panic(err)
			}
			fmt.Println("Established SSH Connection as " + *ssh_user)
			out, err := session.CombinedOutput(os.Args[3])
			if err != nil {
				panic(err)
			}
			fmt.Println(string(out))
		}
	})

	app.Run(os.Args)
}

/*
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
						DeviceToJSON(device)

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
									devices = DeviceArrayToJSON(devices)
								}
								if table_output {
									outputDevicesToTable(devices)
								}
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
								fmt.Println("\nunified devices ugw inspect `mac_address`\n")
								device, _, err := cx.Devices.GetByMac(ctx, c.Args().Get(0))
								if err != nil {
									return nil
								}
								DeviceToJSON(device)
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
									devices = DeviceArrayToJSON(devices)
								}
								if table_output {
									outputDevicesToTable(devices)
								}
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
								fmt.Println("\nunified devices uap inspect `mac_address`\n")
								device, _, err := cx.Devices.GetByMac(ctx, c.Args().Get(0))
								if err != nil {
									return nil
								}
								DeviceToJSON(device)
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
									devices = DeviceArrayToJSON(devices)
								}
								if table_output {
									outputDevicesToTable(devices)
								}
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
								fmt.Println("\nunified devices usw inspect `mac_address`\n")
								device, _, err := cx.Devices.GetByMac(ctx, c.Args().Get(0))
								if err != nil {
									return nil
								}
								DeviceToJSON(device)
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
*/
func DeviceToJSON(device *unified.Device) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(device)
}

func DeviceArrayToJSON(devices []unified.DeviceShort) []unified.DeviceShort {
	table_output = false
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(devices)
	return devices
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

	func setup() {
	client := unified.NewUniFiClient(nil, nil)
	value, err := url.Parse("https://192.168.10.7:8443/api/s/")

	if err != nil {
		//log.Fatal(err)
	} else {
		client.BaseURL = value
	}
}
