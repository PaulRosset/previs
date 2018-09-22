# Previs

> 🎯 Use Travis config to run your CI locally to keep your git history clean.

### Motivation

> Your very own local CI!

Travis is an amazing tool but sometimes before pushing we are producting small mistake that trigger a fail build from Travis. To correct this build we have to alterate the git history to revert the mistakes, thus our git history can become dirty and you may are lazy to re-write it correctly.

But more than a simple travis copy that run locally, **Previs** can be used as a clean environment of test. What I mean by that, is you can easily use **Previs** without using Travis, one can go without the others. 

### [Snap packager](https://snapcraft.io/) for Linux users

One of the strenght of **Previs** is the fact that he is snapped. He is using the [awesome software developed by Canonical](https://docs.snapcraft.io/) which made very easy to install it.  
Because without Snap you will be force to install [Go](https://golang.org/doc/install) which is sometime anoying and can be repulsive to use a tool, so here no more excuses, two command and you are ready to use **Previs**:

- `apt-get install snap`
- `snap install previs`

### How to use Previs

Previs is simple to use, he is using the travis configuration (`.travis.yml`) to configure and provide you everythings:

Once you are at the root of your repository where the `.travis.yml` is, you can launch:

`$> previs`

However, Previs is not supporting all the stuff that Travis is supporting yet, at the moment, he is supporting these:

- Languages:
    - Go
    - Ruby
    - Nodejs
    - Python
    - Php

- Commands:
    - `language`
    - `[nameoflanguage]: [version]`
    - `before_install`
    - `install`
    - `before_script`
    - `script`

As Travis, Previs understand a failed build when the program ran is returning the exit code **2**.

### Contribute

Any contributions is very welcomed, let's do something bigger and stronger together!

Points that will be improved:
- Improve the way of the docker images are wrote before build
- Find a better solutions for the low level images that reside in `baseImages/` folder. They are currently hosted on the official docker registry.
- Adding support for more Languages
- Adding support for more commands (Env Variables are **High** priority)
- Add support for indicating the loading states, when pulling, building and starting.
- The architecture and the code can be improved
- Add other functionalities.?.?...

### Dependencies

- github.com/satori/go.uuid
- gopkg.in/yaml.v2
- github.com/docker/docker/client

### License 

MIT
Paul Rosset