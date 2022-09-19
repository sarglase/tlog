# tlog
a little 、easy、simple go lib for golang


# usage 
## simple use
```


tlog.Info("hello) 

tlog.Error("hello) 

tlog.Debug("hello) 

```

if you want to write log  into file,set the option 
```
tlog.New(tlog.WithHook(hook.New()))
```

you can also do some sets

```
tlog.SetName("service name")
tlog.SetLevel(tlog.ErrorLevel)
```

ok.try it!