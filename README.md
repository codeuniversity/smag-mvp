# Social Record
> Distributed scraping and analysis pipeline for a range of social media platforms

**Table of content**
- [About](#about)
- [Architectural overview](#architectural-overview)
- [Further reading](#further-reading)
  - [Detailed documentation](#detailed-documentation)
  - [Wanna contribute?](#wanna-contribute)
- [Getting started](#getting-started)
  - [Requirements](#requirements)
  - [Scraper in docker](#scraper-in-docker)
  - [Scraper locally](#scraper-locally)

## About
The goal of this project is to raise awareness about data privacy. The mean to do so is a tool to scrape, combine and analyze public data from multiple social media sources. <br>
The results will be available via an API, used for some kind of art exhibition.

## Architectural overview

![](docs/architecture.png)

You can find an more detailed overview [here](https://miro.com/app/board/o9J_kw7a-qM=/)

## Further reading

### Detailed documentation

| part        | docs                                       | contact                                          |
| :---------- | :----------------------------------------- | :----------------------------------------------- |
| Api         | [`api/README.md`](api/README.md)           | [@jo-fr](https://github.com/jo-fr)               |
| Frontend    | [`frontend/README.md`](frontend/README.md) | [@lukas-menzel](https://github.com/lukas-menzel) |
| Postgres DB | [`db/README.md`](db/README.md)             | [@alexmorten](https://github.com/alexmorten)     |

### Wanna contribute?

If you want to join us raising awareness for data privacy have a look into [`CONTRIBUTING.md`](CONTRIBUTING.md)

## Getting started

### Requirements

| depency                                                      | version                                                            |
| :----------------------------------------------------------- | :----------------------------------------------------------------- |
| [`go`](https://golang.org/doc/install)                       | `v1.13` _([go modules](https://blog.golang.org/using-go-modules))_ |
| [`docker`](https://docs.docker.com/install/)                 | `v19.x`                                                            |
| [`docker-compose`](https://docs.docker.com/compose/install/) | `v1.24.x`                                                          |

If this is your first time running this:

1. Add `127.0.0.1 my-kafka` and `127.0.0.1 minio` to your `/etc/hosts` file
2. Choose a `<user_name>` for your platform of choice `<instagram|twitter>` as a starting point and run 
   ```bash
   $ go run cli/main/main.go <instagram|twitter> <user_name>
   ```

### Scraper in docker

```bash
$ make run
```

### Scraper locally

Have a look into [`docker-compose.yml`](docker-compose.yml), set the neccessary environment variables and run it with the command from the regarding dockerfile.
