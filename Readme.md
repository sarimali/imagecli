# Imagecli

A CLI application for Bjorn to help differentiate images quickly

## Installation 

Install it when you install our command-line library:

```sh
go get -u github.com/sarimali/imagecli/...
```

# Usage


```NAME:
      imagecli - Image CLI
   
   USAGE:
      go_build_imagecli_go [global options] command [command options] [arguments...]
   
   VERSION:
      2020.02.02
   
   COMMANDS:
      compare, c  Use this to compare the files listed in the csv
      list, l     List files that will be compared
      help, h     Shows a list of commands or help for one command
   
   GLOBAL OPTIONS:
      --path value   Required path to csv file (default: "~/example.csv")
      --help, -h     show help (default: false)
      --version, -v  print the version (default: false)

```

## Examples

## Generate a report of image differences

```sh
imagecli --path ./input.csv compare

$ ls *results.csv
input.csvresults.csv
```

## Run a sample run of images that will be compared to test if the csv is ok

```sh
imagecli --path ./input.csv l

image/aa.jpg comparing with image/ba.jpg
image/ab.jpg comparing with image/bb.jpg
```