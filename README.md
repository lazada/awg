# Advanced wait group #

AWG -  advanced version of wait group

Added features: 

* execution timeout
* cancel and deadline with context
* ability to return errors
* panic handling from goroutines

## Usage: ##


```
#!go
 	wg := awg.AdvancedWaitGroup{}
	
	// Add function one
	wg.Add(func() error {
		//Logic
		return nil
	})

	// Add function two
	wg.Add(func() error {
		//Another Logic
		return nil
	})
	
	wg.Start()

```


### Getting errors (or one error): ###


```
#!go
	wg := awg.AdvancedWaitGroup{}
	
	// Add function one
	wg.Add(func() error {
		//Logic
		return nil
	})

	// Add function two
	wg.Add(func() error {
		//Another Logic
		return nil
	})
	
        // Taking one error make sense if you use *.SetStopOnError(true)* option - see below
        var err error
	err = wg.SetStopOnError(true).Start().GetLastError()


	// Taking all errors
	var errs []error
	errs = wg.Start().GetAllErrors()

```



### You can add slice of functions to exec in parallel threads: ###


```
#!go

	wg := awg.AdvancedWaitGroup{}

	functions := []awg.WaitgroupFunc{
		func() error {
			//Logic
			return nil
		},
		func() error {
			//Another Logic
			return nil
		},
	}
	
	// Add function one
	wg.AddSlice(functions)

	wg.Start()
```


### You can set timeout for execution: ###


```
#!go

	wg := awg.AdvancedWaitGroup{}
	
	// Add function one
	wg.Add(func() error {
		//Logic
		return nil
	})

	wg.SetTimeout(time.Minute).Start()

```


### You can break execution after first error: ###


```
#!go

	wg := awg.AdvancedWaitGroup{}
	
	// Add function one
	wg.Add(func() error {
		//Logic
		return nil
	})

	// Add function two
	wg.Add(func() error {
		//Another Logic
		return nil
	})

	wg.SetStopOnError(true).Start()
```



### You can reset state by *.Reset()* function ###
