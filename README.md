# log
The log system for GoLang.

## Support:

#### 1. Log level color output. 
#### 2. Log file today split.
#### 3. Error traceback.

## Usage:

#### 1. On project add dir for config.
#### 2. Add logs.toml config file to config dir.
#### 3. The logs.toml file example.

````
[log]
    relative    = true          # true: Default use relative path, false: use absolute path.
    dir         =  "logs"       # Log storage directory.
    name        = "logs.log"    # Log filename.
    prefix      =  ""           # Prefix for cutting log filename.
    level       =  "INFO"       # The log level, if the value is OFF, then log off.
````