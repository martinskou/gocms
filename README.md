# Aviva Wakizashi - Possible the fastest CMS in the world.

Simplicity is a great virtue but it requires hard work to achieve it and 
education to appreciate it. And to make matters worse: complexity sells better.
    Dijkstra (1984)


## Goal

The goal is to supply a fast, simple and useable CMS. It should also be easy to extend
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

Watch for changes and re-bundle:

> go run . watch ../themes/alfa/dependencies

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

## Editors

### Page Editor

Find and place content in page positions. This is done by generating ContentLinks. And change the order of Content in a given position. And move item from one position to a new position...

### Content Editor

Create or edit a Content object. 

Text in a Quill editor.

Select image from media library or unsplash.


### Media library

List of images used in Content objects.

Upload possible.

Show items not linked to Content.



## Features

Some nice to have features

- Intelligent "404 Not found" page.
- Sitemap
- Index


## What it is and is not

It's a content management system. You can edit and manage the content and structure of your websites.

It's not a tool to design or implement webdesign. To change the design of your site, contact your webdesigner.

It's primary goal is to be simple, useable and fast.
- Fast : accomplished quickly
- Simple : readily understood or performed
- Useable : convenient and practicable for use 

Responsiveness is an element of the template and not the CMS. 

We should aim for simplicity because simplicity is a prerequisite for reliability.

Simple is often erroneously mistaken for easy. 
"Easy" means "to be at hand", "to be approachable". 
"Simple" is the opposite of "complex" which means "being intertwined", "being tied together". 
Simple != easy.

The benefits of simplicity are: ease of understanding, ease of change, ease of debugging, flexibility.


## Todo

Graph Viz integration
