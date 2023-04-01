# RGB

Response Generator Bot.

- Manage, organize and share AI interactions
- engineer prompts with prompt generation rules
- integrate input / output targets like:
  - reddit
  - facebook
  - twitter / mastadon
  - slack / discord

## Table of Contents

- [Running](#running)
- [Design](#design)
- [Resources](#resources)

## Running the app
The included Makefile provides 3 commands:

1) Compile the HTMX server into an executable `./server`
```sh
make build
```

2) Compile and run the HTMX server:
```sh 
make run
```

3) **To ease template development**, *make watch* uses [Reflex](https://github.com/cespare/reflex) 
to watch for changes to `pkg/templates/*.gohtml and reload the server when a template changes.
To use it, install reflex:

```sh
go install https://github.com/cespare/reflex@latest
```
then run:
```sh 
make watch
```


## Application Design

## Resources

- [HTMX](https://htmx.org/docs/)
- [Locality of Behavior (LoB)](https://htmx.org/essays/locality-of-behaviour/)
- [HATEOAS](https://htmx.org/essays/hateoas/)

## Contributing
Please contribute using [Github Flow](https://guides.github.com/introduction/flow/). Create a branch, add commits, and [open a pull request](https://github.com/fraction/readme-boilerplate/compare/).
