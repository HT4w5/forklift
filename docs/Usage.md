# Usage

## Config
```json
{
    "log": {
        "level": ""
    },
    "exec": {
        "path": "/usr/bin/sing-box",
        "log_fwd": true
    },
    "profile": {
        "url": "",
        "update": "",
        "ua": "",
        "patches": [
            "patch1"
        ]
    },
    "patches": [
        {
            "tag": "patch1",
            "content": {
                "+foo": [
                    "bar1",
                    "bar2"
                ]
            }
        }
    ]
}
```

#### log
Log settings.
|Field|Description|Type|
|:---:|:---|:---|
|level|Log level. Possible values: "debug", "info", "warn", "error", "none". Defaults to "error".|string|

#### exec
Sing-box execution settings.
|Field|Description|Type|
|:---:|:---|:---|
|path|Path to sing-box binary. Defaults to "sing-box".|string|
|log_fwd|Forward sing-box stdout and stderr to logs. Defaults to false.|bool|

#### profile
Sing-box profile settings.
|Field|Description|Type|
|:---:|:---|:---|
|url|URL to profile.|string|
|update|Crontab for updating profile.|string|
|ua|User agent used when fetching profile.|string|
|patches|Tags of patches to be applied to the raw profile. Must be defined in patches.|[]string|

#### patches
Patches definition. See [patching](Patching.md).
