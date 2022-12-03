# Prowl

Prowl is a web crawler that can be used to crawl websites and extract their links. This is the enhanced version of the original Hakrawler tool, which adds support for multiple output formats, custom HTTP headers, cookies, and proxy settings.

## Installation

To install the enhanced Prowl tool, you need to have the Go programming language installed on your system. You can follow the instructions on the [Go website](https://golang.org/doc/install) to install Go on your system.

Once you have Go installed, you can use the `go get` command to install the enhanced Prowl tool.

$ go get -u github.com/hangyakuzero/prowl
  
  
This will download the source code for the tool and install it on your system. You can then run the `prowl` command to use the tool.

## Usage

To use the enhanced Prowl tool, you need to specify the URL of the website that you want to crawl, as well as any optional flags that you want to use. The basic syntax for the `prowl` command is as follows:

$ prowl [FLAGS] URL


For example, if you want to crawl the website `https://example.com`, you can use the following command:

$ prowl -u https://example.com -d 3 -n 1000 -c 10 -t 10s -f json -o results.json

This command will crawl the website starting at the seed URL `https://example.com`, with a maximum depth of 3, a maximum number of URLs of 1000, a concurrency level of 10, a timeout of 10 seconds, and an output format of json. The results will be written to the file `results.json`.


### Flags

The enhanced Prowl tool supports the following flags:

- `-u`: the seed URL to crawl (required)
- `-d`: the maximum depth of the crawl (default: 3)
- `-n`: the maximum number of URLs to crawl (default: 1000)
- `-c`: the concurrency level (default: 10)
- `-t`: the timeout for HTTP requests (default: 10 seconds)
- `-f`: the output format (json, csv, or plain text) (default: plain text)
- `-o`: the output file (optional)

### Examples

Here are some examples of how to use the enhanced Prowl tool with different flags:

- To crawl the website `https://example.com` and save the results in JSON format to the file `results.json`, you can use the following command:

$ prowl -format=json -save=results.json https://example.com

# Features

The enhanced Prowl tool includes the following features:

- Support for multiple output formats, including JSON, CSV, and plain text.
- Ability to specify custom HTTP headers, cookies, and proxy settings for the request.
- Robust error handling and logging of the results.

## Limitations

The enhanced Prowl tool has the following limitations:

- It can only crawl websites that allow robots and do not have any restrictions on the number of requests or the rate of requests.
- It can only crawl websites that use the HTTP or HTTPS protocol.
- It does not support crawling of websites that require authentication or require a specific user agent or language.
- It does not support crawling of websites that use JavaScript or other client-side technologies to generate their content.
- It does not support crawling of websites that use CAPTCHAs or other anti-bot measures.

## Contributing

If you want to contribute to the enhanced Prowl tool, you can fork the repository on GitHub, make your changes, and then create a pull request. Your changes will be reviewed and merged if they are deemed to be useful and in line with the goals of the project.

## License

The enhanced Prowl tool is released under the [MIT License](LICENSE). This means that you are free to use, modify, and distribute the tool, as long as you include the copyright notice and the license terms in your distribution.

## Credits

The enhanced Prowl tool is based on the original [Hakrawler](https://github.com/hakluke/hakrawler) tool, which was created by [Luke Stephens](https://github.com/hakluke). The enhancements and improvements to the tool were made by [Maharshi](https://github.com/hangyakuzero).
