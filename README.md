# This project serves a web page where you can list the movies of an actor sorted by rating.

## Try it out live at [https://movies-of.bxabi.com](https://movies-of.bxabi.com)

It is written using Revel, a high-productivity web framework for the [Go language](http://www.golang.org/).

### To use it, you have to get an API key from [The Movie Database](https://themoviedb.org), and save it in a file called 'apiKey' in this folder.

### Install Revel: 
    
    go get github.com/revel/revel
    go get github.com/revel/cmd/revel

### Start the web server:

   revel run movies-of \[prod\]
   (~/go/bin/revel if you don't want to add revel to your path)

### Go to http://localhost:9000/

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## Help

* The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/examples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).

## Note:

the routes/routes.go file is autogenerated when the "revel run myapp" is executed.
