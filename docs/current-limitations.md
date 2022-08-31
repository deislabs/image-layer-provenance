# Current Limitations in Container Image Vulnerability Experience

## Limitations of the Container Image History Format

When `docker image history` is run, a simple image history is produced.
The output shows how each non-empty layer (each layer's filesystem change) was created.

```tsv
IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
1a06cd9c6d9d   12 minutes ago   RUN /bin/sh -c echo file2 > /file2.txt # bui…   6B        buildkit.dockerfile.v0
<missing>      12 minutes ago   RUN /bin/sh -c echo file1 > /file1.txt # bui…   6B        buildkit.dockerfile.v0
<missing>      7 days ago       /bin/sh -c #(nop) COPY dir:540bddbe02c37f419…   20.3MB
<missing>      7 days ago       /bin/sh -c #(nop)  ENV ASPNET_VERSION=6.0.8     0B
<missing>      7 days ago       /bin/sh -c ln -s /usr/share/dotnet/dotnet /u…   24B
<missing>      7 days ago       /bin/sh -c #(nop) COPY dir:a07d28c7e9124f9c9…   70.6MB
<missing>      7 days ago       /bin/sh -c #(nop)  ENV DOTNET_VERSION=6.0.8     0B
<missing>      7 days ago       /bin/sh -c #(nop)  ENV ASPNETCORE_URLS=http:…   0B
<missing>      7 days ago       /bin/sh -c apt-get update     && apt-get ins…   37MB
<missing>      2 weeks ago      /bin/sh -c #(nop)  CMD ["bash"]                 0B
<missing>      2 weeks ago      /bin/sh -c #(nop) ADD file:0eae0dca665c7044b…   80.4MB
```

The simple history from `docker image history` has the following problems:

* There is no information that indicates (1) which layers came from base images referenced in `FROM` statements or (2) which layers came from Dockerfile instructions built on top of base image layers.
* No information that maps from a layer digest to a base image ref (registry url, repo, tag, digest).
* In certain cases, the Dockerfile instruction that created a layer is mangled, as seen in the `ADD file:<hash>` and `COPY dir:<hash>` history entries above.

When a vulnerable package is detected, scanners report the following:

* the layer hash in which a vulnerable package was introduced to the filesystem,
* the Dockerfile instruction that created the layer, which in some cases (as seen above) is mangled/obfuscated or contains hashed values.

Reporting only (1) a layer hash and (2) a mangled Dockerfile instruction is not actionable.
With only these 2 pieces of limited information, maintainers cannot differentiate vulnerability alerts between:

* Actionable for vulnerabilities the user can actually fix in their layers, such as vulnerabilities introduced through Dockerfile `ADD, RUN, or COPY` instructions.
* Vulnerabilities in dependent base image layers requiring further action, such as (1) patching the base image, and (2) rebuilding to pull the patched dependencies.

## Limitations of Vulnerability Scan Reports

Vulnerability scanners typically output a long list of vulnerabilities when an image is scanned.
Scanners typically make use of `docker image history` output when generating vulnerability provenance.

A vulnerability report becomes unactionable if it cannot differentiate vulnerabilities between (1) newly-introduced layers or (2) vulnerabilities from base image layers.
It also becomes unactionable if it does not pinpoint the exact source of the vulnerability, requiring manual investigation to do so.

Here is a sample vulnerability report that is unactionable due to information gaps in `docker image history` output:

* The layer command is obfuscated (`ADD file:<hash>`), preventing image builders from identifying the exact layer and Dockerfile instruction which introduced the vulnerable package.
* There is also no indication whether the layer comes from a base image.
It also does not contain the layer's originating base image ref.

```
Vulnerability Name: Ubuntu Security Notification for Sqlite3 Vulnerabilities (USN-2698-1)

Vulnerable Package Information:
{
  "VulnerablePackages": [
    {
      "name": "libsqlite3-0",
      "installedVersion": "3.8.2-1ubuntu2",
      "requiredVersion": "3.8.2-1ubuntu2.1"
    }
  ]
}

Layer information where the vulnerable package was introduced:

{
 "packageMapping": [
    {
      "packageName": "libsqlite3-0",
      "packageVersion": "3.8.2-1ubuntu2",
      "layers": [
        {
          "layerId": 1,
          "layerHash": "fef0f9958347a4b3c846fb8ea394fbcc554ec5440c7ec72b09786230d55ccc03",
          "layerCommand": "ADD file:0a5fd3a659be172e86491f2b94fe4fcc48be603847554a6a8d3bbc87854affec in /"
        }
      ]
    }
  ]
}
```
