# Aviva Wakizashi - Possible the fastest CMS in the world.

Simplicity is a great virtue but it requires hard work to achieve it and 
education to appreciate it. And to make matters worse: complexity sells better.
    Dijkstra (1984)


## Goal

The goal is to supply a simple, useable and fast CMS. It should also be easy to extend
in a controlable manner, and easy to install and update.


## Go tools

Use:

> cd src

Run server using go run:

> go run . startserver ../config/config.json

Build server:

> go build .
> cp gocms ../bin/
> cd ..
> bin/gocms config/config.json

Build some dependency (folder must include build.json file):

> go run . bundle ../dependencies



## CMS

The CMS relies on files from the dependencies folder beeing prepared by running 

go run gocms.go bundle

Which will merge the required CSS and JS files and place these in th folder assets.

The folder assets should only be used by the CMS, not by your own themes.


## Themes

A Theme consists of a folder with the following content:

theme-name
- templates
- assets
-- css
-- img
-- js

A theme can be activated by setting theme-folder in the config file.

{theme: theme-alfa}

On server restart the themes assets are merged

In the theme template, you can access theme using:

<img src="/themes/alfa/assets/img/img.png">

Like the CMS, the themes must only include one bundled CSS file, and one bundled JS file.






