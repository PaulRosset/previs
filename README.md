<h1 align="center">
	<br>
	<img width="220" src="https://rawgit.com/PaulRosset/previs/master/media/logo.png" alt="previs">
	<br>
	<br>
	<br>
</h1>

# Previs

> 🎯 Use Travis config to run your CI locally to keep your git history clean.

[![Travis CI Build Status](https://api.travis-ci.com/PaulRosset/previs.svg?branch=master)](https://travis-ci.org/PaulRosset/previs)
[![Snap Status](https://build.snapcraft.io/badge/PaulRosset/previs.svg)](https://build.snapcraft.io/user/PaulRosset/previs)

### Motivation

> Your very own local CI!

Previs is using the travis configuration mechanism to provide your own local service of continuous integration.
No more accidental error that trigger a fail build, no more plumbing on your git history.
But more than a local trevis, it provides you a way of testing in a sandboxed environment.

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

Or you can still download the binary of release.

### How to use Previs

Previs is simple to use, he is using the travis mechanism to run 'things' in a clean environment.

However, since we are running things locally, we can't reproduce everythings that Travis is providing, so when the `.travis.yml` get too complex, you have the possibility to switch on a custom file named `.previs.yml`, then previs will take the configuration of the `.previs.yml`.

Once you are at the root of your repository where the `.travis.yml` or `.previs.yml` is, you can launch:

`$> previs [-p]`

The `-p` command indicate to previs to take the configuration of the `.previs.yml` instead of the `.travis.yml`.

Previs is for the moment supporting these commands:

- Languages:
    - Normally, all the language are already supported except the only languages name that differ between the travis config and the name of the image registered on the docker registry, in that case we have to add the entry in the dictionary. As an example the nodejs language demonstrate it. In the travis configuration we have to provide the name `node_js` but in the docker registry the official nodejs image is registered as `node`.

- Commands:
    - `language`
    - `[nameoflanguage]: [version]`
    - `before_install`
    - `install`
    - `before_script`
    - `script`

Previs understand a failed build when the program ran is returning other than the **0** exit code.

### Testing

Concerning the workflow of testing, rather than create unit tests on multiple call systems that already been tested especially docker calls, we instead run the program in real world use case to verify nothing broke and prevent regressions.  
To understand it, you can check out the `.travis.yml` file that serve th is purpose.

### Contribute

Any contributions is very welcomed, let's do something bigger and stronger together!

Points that will be improved:
- Improve the way of the docker images are wrote before build
- Adding support for more commands
- Clean when aborting via SIGNALs

### License 

**MIT**  
Paul Rosset
