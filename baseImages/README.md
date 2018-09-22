# Base Image

This is the lowest level of image, they are published on the official docker registry.

### Why base image ?

Many more base image can be published to extend the support for differents platforms.  
Each base image answer differents problems:

- Manage user privilieges
- Installing minimal dependencies at the lowest level
- Creating the same working directory among every base image
- Creating environment

The fact to create a base image permit to let the logic inside the dockerfile instead of creating one from scratch from the go program. Thus it permit to mitigate the go program.

### Available base image

| Platform | Versions Supported |
| -------- | ------------------ |
| node_js  | 8,9,latest         |

### Contributing by adding a base image

