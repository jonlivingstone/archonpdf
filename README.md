## About the project

Archon PDF is a small server that merges two PDF files, one representing the odd and another the even pages (also known as recto and verso).

This utility is useful is you have a scanner with a document feeder, but which cannot process both sides of the pages at once.

Used in combination with a document management system such as [Paperless-ngx](https://github.com/paperless-ngx/paperless-ngx), this enables you to scan and ingest large two-sided documents at great speed!

## Getting started

### Prerequisites

You will need to compile this yourself, but don't worry the setup is very small.

The only dependency is **Golang**. Follow the installation instructions [here](https://go.dev/doc/install).

### Compilation

Use the command `make` to check all dependencies are configured correctly.

## Usage

### Standalone
1. On the machine that will process files, create a directory that your scanner/printer will write pages to. Typically this will be a SMB mount point. For instance `\\server\archonpdf_input`.
2. On your scanner/printer, add two network directories as destinations, one for receiving odd pages and another one to receive even pages. For instance, `\\server\archonpdf_input\odd` and `\\server\archonpdf_input\even`.  
   On some printers you may have to use the web administration at `http://<your printer ip>`
3. Create or identify a directory for the merged documents. This doesn't have to be a SMB share. If you use a document manager such as [Paperless-ngx](https://github.com/paperless-ngx/paperless-ngx), this will be your `consume` directory. Make sure that the process running `archonpdf` will have write access to this directory.
4. Run `archonpdf -input <input_dir> -merged <merge_dir>`.

### Docker

1. Compile the project into a docker image. You can use the makefile.  
   `make install`

2. Run the container. Replace `<input_dir>` and `<merged_dir>` with the directories you want archonpdf to process.  
   `docker run --rm -d -v <input_dir>:/input -v <merged_dir>:/merged --name archonpdf archonpdf:snapshot`  

### TrueNAS

The Docker image for this project is not published yet, but we can install it manually.

1. Create the SMB share to which the printer/scanner will write scanned documents.

2. Create the TAR export for the Docker image.  
   `make export`

3. Copy the generated file `archon-pdf.snapshot.tar` to your server.

4. Open a shell and load the docker image in the local repository. This will require administration privileges (sudo).  
   `docker image load -i archonpdf-snapshot.tar`

5. You can now create your TrueNAS app.

    Method A: fill the fields

    * Image: archonpdf
    * Tag: snapshot

    * Storage configuration:
        * Input folder  
            Mount path: `/input`  
            Host path: `<input_directory>`
        * Merged folder. This will typically point to the `consume` directory for `paperless`.    
            Mount path: `/merged`  
            Host path: `<merged_directory>`

    Method B: use a YAML description. You may adapt the following template.

```
services:
  archonpdf:
    environment:
      - PUID=568
      - PGID=568
    image: archonpdf:snapshot
    restart: unless-stopped
    volumes:
      - <input_directory>:/input
      - <merged_directory>:/merged
```

## Potential improvements

1. Publish the docker image

2. Publish a TrueNAS app.
