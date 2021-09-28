# Docs

A few technical aspects are kept documented here for future reference.
For the most part, functions will be commented.

## Redis JSONPath structure example

```bash
.
├── beats
│   ├── <HeartbeatBeat>        # {"device_name": "laptop", "timestamp": 1632748096}
│   ├── <HeartbeatBeat>        # {"device_name": "phone", "timestamp": 1632748137}
│   └── <HeartbeatBeat>        # {"device_name": "laptop", "timestamp": 1632748682}
│
├── devices
│   ├── <HeartbeatDevice>      # {"device_name": "laptop", "total_beats": 12903, "longest_missing_beat", 1201}
│   └── <HeartbeatDevice>      # {"device_name": "phone", "total_beats": 1952, "longest_missing_beat", 3219}
│
├── last_beat <HeartbeatBeat>  # {"device_name": "laptop", "timestamp": 1632748096}
└── stats <HeartbeatStats>     # {"total_visits": 45012, "total_uptime": 892340, "total_beats": 14855, "longest_missing_beat", 3219}
```
