# PLC Gateway Service

## Dependencies Requires
```
go get -u github.com/swaggo/swag/cmd/swag
```

## Compile & Build
```shell
go mod vendor && make
```

## Run
```bash
# Linux
./dist/plc-gw

# OR on Windows
./dist/plc-gw.exe
```

## Service Managerment

using `./plc-gw <COMMAND>` on linux OR `./plc-gw.exe <COMMAND>` on windows

- Commands
  ```
  start           Startup the service
  stop            Stop the service
  restart         Restart the service
  install         Install as service
  uninstall       Uninstall the service
  ```
  
- Example
  ```bash
  # Install as service on linux
  ./plc-gw install

  # OR on windows
  ./plc-gw.exe install
  ```

## Online API Docs
```
http://localhost:8080/swagger/index.html
```