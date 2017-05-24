# Unified - The CLI Tool for Ubiquiti UniFi Devices

The Unified CLI tool allows interacting with the UniFi Controller from the command line. Unified impolements a series of
Commands, Sub-Commands and Flags in a similar fashion to Git or Docker.

## Commands & SubCommands

The Commands and Sub-Commands are as follows: -

```
unified 
         --help
         --version
         controller
                 --help
                 alarms
                         ls       
                 event
                         ls
         devices
                --help
                ls
                inspect MAC_ADDRESS
                ugw
                        --help
                        ls
                        ps
                        inspect MAC_ADDRESS    
                uap
                        --help
                        ls
                        ps
                        inspect MAC_ADDRESS
                usw
                        --help
                        ls
                        ps
                        inspect MAC_ADDRESS
                        cn
                        pt
                        po
                        dt
                        up
                        eth
                        ss
                        stats
                        ssh  
         exec
                --help
         guest
         operator
                --help
                ls
                inspect
         site
                --help
                ls
                ps
                inspect
```

### Example Commands
If we wish to see a list of Alarms on the controller then we would use the command: -
 
 `unified controller alarms ls`
 
which would return something similar to: -
 
 ```
+--------------------------+----------+----------------------+-------+--------------------------+----------------------+---------------------+------------+--------------------------------+--------+--------------------------+-----------+
|           UUID           | ARCHIVED |       DATETIME       | ESSID |      HANDLEDADMINID      |     HANDLEDTIME      |         KEY         | MACADDRESS |            MESSAGE             | OCCURS |          SITEID          | SUBSYSTEM |
+--------------------------+----------+----------------------+-------+--------------------------+----------------------+---------------------+------------+--------------------------------+--------+--------------------------+-----------+
| 590487c9e4b01c675d52a786 | true     | 2017-04-29T12:32:09Z |       | 58def83ee4b0dfb95e1163bc | 2017-04-29T21:37:52Z | EVT_AP_Lost_Contact |            | AP[80:2a:a8:d9:99:00] was      |      0 | 58def75ce4b0dfb900000000 | wlan      |
|                          |          |                      |       |                          |                      |                     |            | disconnected                   |        |                          |           |
+--------------------------+----------+----------------------+-------+--------------------------+----------------------+---------------------+------------+--------------------------------+--------+--------------------------+-----------+
 ```
 
 or if we wanted to list only UniFi Switches adopted by the contoller: -
 
`unified devices usw ls`
  
which would display something similar to
  ```
  +------+--------------+--------------------------+-------------------+-------------------------+---------+-------------+--------------------------+-----------+-------------+
  | TYPE |    SERIAL    |          SITEID          |    MACADDRESS     |          NAME           |  MODEL  |     IP      |           UUID           | ISADOPTED |   VERSION   |
  +------+--------------+--------------------------+-------------------+-------------------------+---------+-------------+--------------------------+-----------+-------------+
  | usw  | F09FC0000000 | 58def75ce4b0dfb900000000 | f0:9f:c2:60:00:00 | UniFi Switch 8 POE-150W | US8P150 | 192.168.1.4 | 58def83fe4b0dfb95e1163c0 | true      | 3.7.55.6308 |
  +------+--------------+--------------------------+-------------------+-------------------------+---------+-------------+--------------------------+-----------+-------------+
  ```
  
  Where multiple elements are returned e.g. 3 UniFi AP using the command
  
`unified devices uap ls`
    
  it would look similar to the following: -
  
  ```
  +------+--------------+--------------------------+-------------------+----------------+-------+--------------+--------------------------+-----------+-------------+
  | TYPE |    SERIAL    |          SITEID          |    MACADDRESS     |      NAME      | MODEL |      IP      |           UUID           | ISADOPTED |   VERSION   |
  +------+--------------+--------------------------+-------------------+----------------+-------+--------------+--------------------------+-----------+-------------+
  | uap  | 802AA8C66367 | 58def75ce4b0dfb95e1163ab | 80:2a:a8:c6:63:67 | Manse Landing  | U7LT  | 192.168.1.10 | 58df0767e4b0dfb95e116437 | true      | 3.7.55.6308 |
  | uap  | 802AA8D9C521 | 58def75ce4b0dfb95e1163ab | 80:2a:a8:d9:c5:21 | Manse Hallway  | U7PG2 | 192.168.1.20 | 58defcb8e4b0dfb95e1163f5 | true      | 3.7.55.6308 |
  | uap  | F09FC22010D1 | 58def75ce4b0dfb95e1163ab | f0:9f:c2:20:10:d1 | Manse Orangery | U7LT  | 192.168.1.11 | 58df07c9e4b0dfb95e11643b | true      | 3.7.55.6308 |
  +------+--------------+--------------------------+-------------------+----------------+-------+--------------+--------------------------+-----------+-------------+
  
  ```
  
#### The `inspect` Command
  The exception to this tabular output format is where an inspect is done. Due to the large amounts of data returned and
  given that an inspect only can be issued against a single device the output is in JSON format. 
  
  e.g. for a usw i.e. Unified UniFi Switch, the command would be: -
  
`unified devices info MAC_ADDRESS`

then the output would be of the form: -

```json
{
    "_id": "58def83fe4b0dfb95e1163c0",
    "adopted": true,
    "board_rev": 12,
    "cfgversion": "4318f00c56e5c563",
    "config_network": {
        "dns1": "192.168.1.3",
        "dns2": "194.72.9.38",
        "dnssuffix": "portree.ecosse-hosting.net",
        "gateway": "192.168.1.1",
        "ip": "192.168.1.4",
        "netmask": "255.255.240.0",
        "type": "static"
    },
    "connect_request_ip": "192.168.1.4",
    "connect_request_port": "45808",
    "device_id": "58def83fe4b0dfb95e1163c0",
    "downlink_table": [
        {
            "full_duplex": true,
            "mac": "80:2a:a8:d9:c5:21",
            "port_idx": 2,
            "speed": 1000
        },
        {
            "full_duplex": true,
            "mac": "80:2a:a8:c6:63:67",
            "port_idx": 3,
            "speed": 1000
        },
        {
            "full_duplex": true,
            "mac": "f0:9f:c2:20:10:d1",
            "port_idx": 4,
            "speed": 1000
        }
    ],
    "ethernet_table": [
        {
            "mac": "f0:9f:c2:60:28:7d",
            "name": "eth0",
            "num_port": 10
        },
        {
            "mac": "f0:9f:c2:60:28:7e",
            "name": "srv0"
        }
    ],
    "fw_caps": 5,
    "general_temperature": 68,
    "guest-num_sta": 1,
    "inform_authkey": "104882425765cca451b980b6dd1612bf",
    "inform_ip": "192.168.10.7",
    "inform_url": "http://192.168.10.7:8080/inform",
    "ip": "192.168.1.4",
    "mac": "f0:9f:c2:60:28:7d",
    "name": "UniFi Switch 8 POE-150W",
    "model": "US8P150",
    "known_cfgversion": "4318f00c56e5c563",
    "led_override": "default",
    "port_overrides": [
        {
            "autoneg": true,
            "name": "Port 1",
            "op_mode": "switch",
            "poe_mode": "auto",
            "port_idx": 1,
            "portconf_id": "58def767e4b0dfb95e1163b5"
        },
        {
            "autoneg": true,
            "name": "Port 3",
            "op_mode": "switch",
            "poe_mode": "pasv24",
            "port_idx": 3,
            "portconf_id": "58def767e4b0dfb95e1163b5"
        },
        {
            "autoneg": true,
            "name": "Port 4",
            "op_mode": "switch",
            "poe_mode": "pasv24",
            "port_idx": 4,
            "portconf_id": "58def767e4b0dfb95e1163b5"
        }
    ],
    "port_table": [
        {
            "autoneg": true,
            "bytes-r": 5453,
            "dot1x_mode": "auto",
            "dot1x_status": "authorized",
            "enable": true,
            "full_duplex": true,
            "is_uplink": true,
            "media": "GE",
            "name": "Port 1",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Unknown",
            "poe_current": "0.00",
            "poe_mode": "auto",
            "poe_power": "0.00",
            "poe_voltage": "0.00",
            "port_idx": 1,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "rx_broadcast": 24035,
            "rx_bytes": 112046410012,
            "rx_bytes-r": 2969,
            "rx_multicast": 34846,
            "rx_packets": 84941015,
            "speed": 1000,
            "stp_pathcost": 20000,
            "stp_state": "forwarding",
            "tx_broadcast": 1919076,
            "tx_bytes": 7337154370,
            "tx_bytes-r": 2484,
            "tx_multicast": 1023977,
            "tx_packets": 50532849,
            "up": true
        },
        {
            "autoneg": true,
            "bytes-r": 809,
            "dot1x_mode": "auto",
            "dot1x_status": "authorized",
            "enable": true,
            "full_duplex": true,
            "media": "GE",
            "name": "Port 2",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Class 0",
            "poe_current": "89.96",
            "poe_enable": true,
            "poe_good": true,
            "poe_mode": "auto",
            "poe_power": "4.80",
            "poe_voltage": "53.40",
            "port_idx": 2,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "rx_broadcast": 364009,
            "rx_bytes": 2859262738,
            "rx_bytes-r": 403,
            "rx_dropped": 35511,
            "rx_multicast": 281572,
            "rx_packets": 13203066,
            "speed": 1000,
            "stp_pathcost": 20000,
            "stp_state": "forwarding",
            "tx_broadcast": 1577444,
            "tx_bytes": 24143287135,
            "tx_bytes-r": 405,
            "tx_multicast": 1825076,
            "tx_packets": 22862462,
            "up": true
        },
        {
            "autoneg": true,
            "bytes-r": 4599,
            "dot1x_mode": "auto",
            "dot1x_status": "authorized",
            "enable": true,
            "full_duplex": true,
            "media": "GE",
            "name": "Port 3",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Unknown",
            "poe_current": "0.00",
            "poe_enable": true,
            "poe_mode": "pasv24",
            "poe_power": "0.00",
            "poe_voltage": "24.00",
            "port_idx": 3,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "rx_broadcast": 244444,
            "rx_bytes": 660528897,
            "rx_bytes-r": 1306,
            "rx_dropped": 4680,
            "rx_multicast": 135710,
            "rx_packets": 4017536,
            "speed": 1000,
            "stp_pathcost": 20000,
            "stp_state": "forwarding",
            "tx_broadcast": 1698623,
            "tx_bytes": 5400855809,
            "tx_bytes-r": 3293,
            "tx_multicast": 948729,
            "tx_packets": 7881729,
            "up": true
        },
        {
            "autoneg": true,
            "bytes-r": 653,
            "dot1x_mode": "auto",
            "dot1x_status": "authorized",
            "enable": true,
            "full_duplex": true,
            "media": "GE",
            "name": "Port 4",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Unknown",
            "poe_current": "0.00",
            "poe_enable": true,
            "poe_mode": "pasv24",
            "poe_power": "0.00",
            "poe_voltage": "24.00",
            "port_idx": 4,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "rx_broadcast": 226714,
            "rx_bytes": 396900977,
            "rx_bytes-r": 240,
            "rx_dropped": 7871,
            "rx_multicast": 52632,
            "rx_packets": 1656848,
            "speed": 1000,
            "stp_pathcost": 20000,
            "stp_state": "forwarding",
            "tx_broadcast": 1716330,
            "tx_bytes": 2899970291,
            "tx_bytes-r": 413,
            "tx_multicast": 1819854,
            "tx_packets": 6042140,
            "up": true
        },
        {
            "autoneg": true,
            "dot1x_mode": "n/a",
            "dot1x_status": "n/a",
            "enable": true,
            "media": "GE",
            "name": "Port 5",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Unknown",
            "poe_current": "0.00",
            "poe_mode": "auto",
            "poe_power": "0.00",
            "poe_voltage": "0.00",
            "port_idx": 5,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "stp_state": "disabled"
        },
        {
            "autoneg": true,
            "dot1x_mode": "n/a",
            "dot1x_status": "n/a",
            "enable": true,
            "media": "GE",
            "name": "Port 6",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Unknown",
            "poe_current": "0.00",
            "poe_mode": "auto",
            "poe_power": "0.00",
            "poe_voltage": "0.00",
            "port_idx": 6,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "stp_state": "disabled"
        },
        {
            "autoneg": true,
            "dot1x_mode": "n/a",
            "dot1x_status": "n/a",
            "enable": true,
            "media": "GE",
            "name": "Port 7",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Unknown",
            "poe_current": "0.00",
            "poe_mode": "auto",
            "poe_power": "0.00",
            "poe_voltage": "0.00",
            "port_idx": 7,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "stp_state": "disabled"
        },
        {
            "autoneg": true,
            "bytes-r": 2804,
            "dot1x_mode": "auto",
            "dot1x_status": "authorized",
            "enable": true,
            "full_duplex": true,
            "media": "GE",
            "name": "Port 8",
            "op_mode": "switch",
            "poe_caps": 7,
            "poe_class": "Class 0",
            "poe_current": "49.19",
            "poe_enable": true,
            "poe_good": true,
            "poe_mode": "auto",
            "poe_power": "2.63",
            "poe_voltage": "53.40",
            "port_idx": 8,
            "port_poe": true,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "rx_broadcast": 14,
            "rx_bytes": 1062062793,
            "rx_bytes-r": 843,
            "rx_dropped": 4453,
            "rx_multicast": 13185,
            "rx_packets": 3042485,
            "speed": 1000,
            "stp_pathcost": 20000,
            "stp_state": "forwarding",
            "tx_broadcast": 1943093,
            "tx_bytes": 2115091623,
            "tx_bytes-r": 1961,
            "tx_multicast": 1024955,
            "tx_packets": 5906675,
            "up": true
        },
        {
            "autoneg": true,
            "bytes-r": 6625,
            "dot1x_mode": "auto",
            "dot1x_status": "authorized",
            "enable": true,
            "full_duplex": true,
            "media": "SFP",
            "name": "SFP 1",
            "op_mode": "switch",
            "port_idx": 9,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "rx_broadcast": 1049992,
            "rx_bytes": 6546168915,
            "rx_bytes-r": 3833,
            "rx_dropped": 142379,
            "rx_multicast": 3365662,
            "rx_packets": 39204375,
            "speed": 1000,
            "stp_pathcost": 20000,
            "stp_state": "forwarding",
            "tx_broadcast": 893084,
            "tx_bytes": 82715547276,
            "tx_bytes-r": 2792,
            "tx_multicast": 995394,
            "tx_packets": 64407627,
            "up": true
        },
        {
            "autoneg": true,
            "dot1x_mode": "n/a",
            "dot1x_status": "n/a",
            "enable": true,
            "media": "SFP",
            "name": "SFP 2",
            "op_mode": "switch",
            "port_idx": 10,
            "portconf_id": "58def767e4b0dfb95e1163b5",
            "stp_state": "disabled"
        }
    ],
    "serial": "F09FC260287D",
    "site_id": "58def75ce4b0dfb95e1163ab",
    "state": 1,
    "stp_priority": "32768",
    "stp_version": "rstp",
    "type": "usw",
    "uplink_depth": 1,
    "version": "3.7.55.6308"
}
```