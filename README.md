<h1 align="center">
	<br>
	<img width="220" src="https://rawgit.com/PaulRosset/previs/master/media/logo.png" alt="previs">
	<br>
	<br>
	<br>
</h1>

# Previs

> 🎯 Use Travis config to run your CI locally to keep your git history clean.

[![Snap Status](https://build.snapcraft.io/badge/PaulRosset/previs.svg)](https://build.snapcraft.io/user/PaulRosset/previs)

### Motivation

> Your very own local CI!

Travis is an amazing tool but sometimes before pushing we are producting small mistake that trigger a fail build from Travis. To correct this build we have to alterate the git history to revert the mistakes, thus our git history can become dirty and you may are lazy to re-write it correctly.

But more than a simple travis copy that run locally, **Previs** can be used as a clean environment of test. What I mean by that, is you can easily use **Previs** without using Travis, one can go without the other. 

### Installation

#### [Snap packager](https://snapcraft.io/) for Linux users

One of the strength of **Previs** is the fact that he is snapped. He is using the [snapcraft](https://docs.snapcraft.io/) which made very easy to install on linux system. The snapped software are updated automatically.  
You have to run two commands and you will be ready to use **Previs**:

- `$> sudo apt install snapd`
- `$> sudo snap install previs`

For others package manager [see](https://docs.snapcraft.io/core/install).

#### For others users

If you run an other operating system or you don't want to install it via snap, you can still install it "manually" by installing [Go](https://golang.org/doc/install) first, then by doing:

- `$> go get -v github.com/PaulRosset/previs`

### How to use Previs

Previs is simple to use, he is using the travis configuration (`.travis.yml`) to configure and provide you everythings:

Once you are at the root of your repository where the `.travis.yml` is, you can launch:

`$> previs`

However, Previs is not supporting all the stuff that Travis is supporting yet, at the moment, he is supporting these:

- Languages:
    - ~~Go~~ (dev)
    - ~~Ruby~~ (dev)
    - Nodejs
    - ~~Python~~ (dev)
    - ~~Php~~ (dev)

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

### License 

**MIT**  
Paul Rosset