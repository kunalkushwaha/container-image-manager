Overall expectations from this image-manager tool
Taken from comment from #stgraber 

The lxc command line tool has a concept of "remotes" which can be:

A LXD host
An image server using the lxc-images protocol (such as https://images.linuxcontainers.org)
An image server using the system-image protocol (such as https://system-image.ubuntu.com)
A registry pointing to a set of remotes using one of the aforementioned protocols (such as https://registry.linuxcontainers.org)

The default remote configuration will be:
 - A remote called "local" marked as default and set to talk using the LXD protocol over a unix socket to the local LXD daemon.
 - This remote is only automatically added if a local LXD daemon is running.
 - A remote called "images" using the registry protocol and the https://registry.linuxcontainers.org server.
 - Those URLs shouldn't ever be hardcoded deep in the code, they should only be set in whatever generates the default remote configuration as anyone's free to implement their own registry server if they so wish (and can then add it with "lxc remote add my-registry https://my-registry.local").

Now let's take a couple of concrete examples and what I'd expect them to do behind the scenes:
lxc start images:lxc-images/ubuntu/trusty/amd64/default blah

This will do the following:
 - Look into the list of remotes for "images" and extract its url (https+registry://registry.linuxcontainers.org), trusted_keys (if present) and trusted_certs (if present)
 - Based on that URL, we know that the protocol is "registry"
 - Fetch https://registry.linuxcontainers.org/1.0/index.json and https://registry.linuxcontainers.org/1.0/index.json.asc
 - Validate the GPG signature against the keys in trusted_keys
 - Parse index.json and look for "lxc-images" then just as with the local remote lookup earlier, extract its url (https+lxc-images://images.linuxcontainers.org), trusted_keys (if present) and trusted_certs (if present).
 - Based on that URL, we know that the protocol is "lxc-images"
 - Fetch https://images.linuxcontainers.org/meta/1.0/index-user and https://images.linuxcontainers.org/meta/1.0/index-user.asc
 - Validate the GPG signature against the keys we got in trusted_keys for that remote in the registry
 - Parse index-user and look for "ubuntu;trusty;amd64;default", extract the URL (last column)
 - Download $URL/meta.tar.xz, $URL/meta.tar.xz.asc, $URL/rootfs.tar.xz and $URL/rootfs.tar.xz.asc
 - Validate all of those using GPG (same as index-user)
 - Create the container (since we have what's needed for that now)
 - Parse the config inside meta.tar.xz and load it into the DB
 - Unpack rootfs.tar.xz as the container filessytem

That's a relatively high level overview, especially for the last few steps. We'll also have to add some caching to that and look at images expiry for lxc-images, but that's not a priority right now.

Oh, one more thing that might be confusing right now. LXD is currently calling create() directly from the go-lxc binding.
That's just a temporary solution to unblock us, but isn't something that'll stay. Once we have the image store support, that code will go away and LXD will do the container creation itself.
