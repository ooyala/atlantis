Atlantis
========

Atlantis is an Open Source PaaS for HTTP applications built on Docker and written in Go.  It makes it easy to build and deploy applications in a safe, repeatable fashion, and flexibly route requests to the appropriate containers.

We're using Atlantis heavily at Ooyala for new applications; while it still has some rough edges around getting it up and running, the experience of using it for deploying applications is fairly smooth.

For an introduction or more information on how to use Atlantis, see [A User's Guide to Getting Started With Atlantis](http://ooyala.github.io/atlantis/user-getting-started.html).  For more technical details on the project, continue reading here.


This Repository
---------------

This repository is shared across all components; it contains shared datatypes and utility functions.  It is included as a submodule in the other components to make these available without code duplication.  (At the time we started Atlantis, Go package managers were generally fairly immature; it might make sense to switch to one of those, but for the moment, submodules work.)


Architecture
------------

Atlantis consists of multiple components, each with their own repository:

\[insert Cad's image here\]

Each component is briefly discussed below, and detailed in more information in its project repository.

- The [manager](https://github.com/ooyala/atlantis-manager), as its name suggests, manages the whole system.  It is the primary interaction point with Atlantis: it coallates information from the various sources and provides APIs to interact with them.  In the near future, some of this functionality will be moved into the components, but the manager will still be the center of a full Atlantis cluster.

- The [supervisors](https://github.com/ooyala/atlantis-supervisor) run the deployed applications in Docker containers.  They also handle network security; we use iptables to provide a level of isolation similar to EC2 security groups.  If a container doesn't request permissions to a known service, that service is blocked from within that container.

- The [routers](https://github.com/ooyala/atlantis-router) handle routing of HTTP requests based on flexible rules read out of Zookeeper.

- The [builder](https://github.com/ooyala/atlantis-builder) builds containers from known templates based on a simple configuration file, and pushes the images to the registry for deployment to the supervisors.

- The [registry](https://github.com/ooyala/go-docker-registry/) is a Go reimplentation of (most of) the [original Docker registry](https://github.com/docker/docker-registry).  We had some issues with the Python implementation that resulted in us building our own.  This will likely be deprecated in favor of the [Registry 2.0](https://github.com/docker/distribution) once it's fully stable.

In addition, we depend on the following pre-existing services:

- [Zookeeper](https//zookeeper.apache.org/) is used to store configuration data and ensure that updates propagate immediately.

- [Jenkins](https://jenkins-ci.org/) is used to manage building Docker containers for the applications.  We also have a minimal standalone build server for testing, but it doesn't provide niceties like logging or auditing.

- [LDAP](http://en.wikipedia.org/wiki/Lightweight_Directory_Access_Protocol) is optionally used to manage ownership of applications.  Teams of multiple users can be assigned applications, and only users in that team will be able to affect those applications.  Teams can be managed through the Atlantis manager or with any other LDAP interface.


Aquarium
--------
Since running all of these pieces in development can be a pain, we've developed [aquarium](https://github.com/ooyala/atlantis-aquarium), a script to make it easy to run Atlantis within a Vagrant instance.  Following the directions in its Readme, you should be able to get a full Atlantis cluster up for development in minutes (though downloading the various images and dependences can take some time depending on your connection speed.)


Regions and zones
-----------------

One of Ooyala's core engineering principles is avoiding a single point of failure by using [multiple regions in EC2](http://engineering.ooyala.com/blog/staying-when-amazon-web-services-isn%E2%80%99t).  Atlantis supports this model by running multiple independent clusters in different region.  Each has its copy of every component except for the builder; the builder is only in one region, to ensure consistency, and the actual images are pushed to S3, so are available from all regions.

For finer-grained redundancy, Atlantis also supports [availability zones](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html).  The manager knows what zones are available to it, and each supervisor is configured with its zone.  When an application is deployed, it is deployed among all availability zones; this ensures that an issue limited to a single availability zone issue will not take out an application.

Routers should similarly be split among availability zones to ensure uptime.
