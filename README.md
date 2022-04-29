# Fileserver

## Run server
```
go run main.go server
```

## Run client
```
go run main.go client
```

### Commands available in client

#### Set a username
```
/username USERNAME
```

#### Get all channels
```
/channels
```

#### Subscribe to channel X (If you subscribe to a channel, you will not able to use other commands. To exit, finish the console process)
```
/subscribe CHANNEL
```

#### Send to channel X
```
/send CHANNEL PATH or FILENAME
```

#### Exit
```
/exit
```