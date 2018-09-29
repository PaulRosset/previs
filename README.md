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

Travis is an amazing tool but sometimes before pushing we are producting small mistake that trigger a fail build from Travis. To correct this build we have to alterate the git history to revert the mistakes, thus our git history can become dirty or you may are lazy to re-write it correctly.

But more than a simple travis copy that run locally, **Previs** can be used as a clean environment of testing but also run application in it to test their inner working. What I mean by that, is you can easily use **Previs** without using Travis, one can go without the other.

**Previs** is still in active development.

### Installation

Previs is using docker at his heart, so of course you will need the docker deamon in order to use it.  
One of the simplest way to install it, it's by using the bash script:
```
$> curl -fsSL get.docker.com -o get-docker.sh
$> sh get-docker.sh
```

If you are on Mac, it's [here](https://docs.docker.com/docker-for-mac/install/#install-and-run-docker-for-mac)

You will may have some troubleshooting when installing docker especially with the right management of it, if so, do the following:  
`sudo usermod -a -G docker $USER`  
If it still persist after, reboot your system.

#### [Snap packager](https://snapcraft.io/) for Linux users

One of the strength of **Previs** is the fact that he is snapped. He is using the [snapcraft](https://docs.snapcraft.io/) which made very easy to install on linux system. The snapped software are updated automatically.  
You have to run two commands and you will be ready to use **Previs**:

- `$> sudo apt install snapd`
- `$> sudo snap install previs --beta`

For others package manager [see](https://docs.snapcraft.io/core/install).

> Snap binary are located in `/snap/bin/`, make sure it is integrated in your $PATH env.

#### For others users

If you run an other operating system or you don't want to install it via snap, you can still install it "manually" by installing [Go](https://golang.org/doc/install) first, then by doing:

- `$> go get -v github.com/PaulRosset/previs`

Then make sure that your env variable $PATH contain the path where the go binary live.

### How to use Previs

Previs is simple to use, he is using the travis configuration (`.travis.yml`) to configure and provide you everythings:

Once you are at the root of your repository where the `.travis.yml` is, you can launch:

`$> previs`

However, Previs is not supporting all the stuff that Travis is supporting yet, at the moment, he is supporting these:

- Languages:
    - Go
    - ~~Ruby~~ (dev)
    - Nodejs
    - Python
    - ~~Php~~ (dev)

Go [here](https://github.com/PaulRosset/previs/tree/master/baseImages) to check it out the supported language and version

- Commands:
    - `language`
    - `[nameoflanguage]: [version]`
    - `before_install`
    - `install`
    - `before_script`
    - `script`
    - `after_script`

Previs understand a failed build when the program ran is returning other than the **0** exit code.

### Contribute

Any contributions is very welcomed, let's do something bigger and stronger together!

Points that will be improved:
- Improve the way of the docker images are wrote before build
- Find a better solutions for the low level images that reside in `baseImages/` folder. They are currently hosted on the official docker registry. The workflow to add a version or a language support is not very convenient, because we have to create the images locally then send it to the official docker registry, can we a better way to do this ?
- Adding support for more Languages
- Adding support for more commands (Env Variables are **High** priority)
- Clean when aborting via SIGNALs
- Add other functionalities...

### License 

**MIT**  
Paul Rosset
