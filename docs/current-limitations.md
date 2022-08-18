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

* There is no information that indicates (1) which layers came from base images or (2) which layers came from Dockerfile instructions built on top of base image layers.
* In certain cases, the Dockerfile instruction that created a layer is mangled, as seen in the `ADD file:<hash>` and `COPY dir:<hash>` history entries above.

When a vulnerable package is detected, scanners report the following:

* the layer in which a vulnerable package was introduced to the filesystem (ex. the layer created by `RUN pip install vuln-pkg`),
* the Dockerfile instruction that created the layer, which in some cases (as seen above) is mangled/obfuscated or contains hashed values.

Scan results contain noise and information gaps (as seen above), making scan results unactionable.

## Limitations of Vulnerability Scan Reports Due to Being Unactionable

Vulnerability scanners typically output a long list of vulnerabilities when an image is scanned.
A vulnerability report entry becomes unactionable if:

* The vulnerability was introduced from an imported base image's layer.
In other words, the vulnerability does not come from a Dockerfile instruction created on top of base image layers.
  * The vulnerability report entry will point out the vulnerability's layer.
  * However, there is no information indicating whether a layer comes from an imported base image or not.
  * Teams need to manually investigate the vulnerable package and layer's origin.
  * If it is from a base image layer, this type of vulnerability is unactionable because image builders need to wait for base image maintainers to deliver a fixed base image.
* The vulnerability report entry contains the layer digest in which the vulnerability got introduced to the filesystem.
However, it can contain a mangled/obfuscated Dockerfile instruction in its history (such as `COPY dir:<hash>` or `ADD file:<hash>`).
This is because `docker image history` output (which scanners rely on) contains information gaps.
  * Teams cannot use the obfuscated image history and mangled Dockerfile instruction to find the source of the vulnerability in the source Dockerfile.
  * Teams need to manually investigate the package and layer's origin without making use of the obfuscated image and layer history.

Here is a sample vulnerability report that is unactionable due to information gaps in `docker image history` output:

* The layer command is obfuscated (`ADD file:<hash>`), preventing image builders from identifying the exact layer and Dockerfile instruction which introduced the vulnerable package.
* There is also no indication whether the layer comes from an imported base image's layers or not.

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
