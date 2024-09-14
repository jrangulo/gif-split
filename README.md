## Setup and Running Instructions

```
go mod vendor
go run main.go
```

## Dependecies used
[Templ](https://github.com/a-h/templ) For creating dynamic html from the server
[htmx](https://htmx.org/) For rendering/fetching on the client
All other functions are part of the Go standard library

## Approach

A simple program written using htmx/templ that receives an animated GIF, splits the image according to user defined rows and columns, and
returns pure html to the client with the split image in a grid format. Server-side rendering was chosen in order to simplify 
gif processing/sincronization. 

This approach does have one main drawback, after a certain amount of rows/columns, download size on the client side starts to get high.
Possible alternate solutions include serving the gifs as static files in order to leverage html transit compression (such as gzip), but this introduces
other considerations such as animation sincronization on the client and temporary file management on the server side
