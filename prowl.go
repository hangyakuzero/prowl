// Import the necessary packages and libraries
import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"

    "encoding/csv"
    "encoding/json"
)

// Define the default values for the options
const (
    defaultTimeout = 30 * time.Second
    defaultDepth   = 5
    defaultExclude = `(?i)\.(css|js|jpg|jpeg|png|gif|pdf|zip|tar|gz|bz2|mp3|mp4)$`
)

// Define a struct to hold the options of the tool
type Options struct {
    URL       string        // URL to crawl
    UserAgent string        // User-Agent
    Headers   string        // Custom headers
    Cookies   string        // Cookies
    Depth     int           // Maximum depth
    Timeout   time.Duration // HTTP timeout
    Exclude   string        // Exclude regex
    Format    string        // Output format
    File      string        // Output file
}

// Define a function to parse the command-line arguments
func ParseArgs() (*Options, error) {
    // Create a new options instance
    options := &Options{}

    // Set the default values
    options.Timeout = defaultTimeout
    options.Depth = defaultDepth

    // Parse the command-line arguments
    flag.StringVar(&options.URL, "u", "", "URL to crawl")
    flag.StringVar(&options.UserAgent, "ua", "", "User-Agent")
    flag.StringVar(&options.Headers, "H", "", "Custom headers")
    flag.StringVar(&options.Cookies, "C", "", "Cookies")
    flag.IntVar(&options.Depth, "d", defaultDepth, "Maximum depth")
    flag.DurationVar(&options.Timeout, "t", defaultTimeout, "HTTP timeout")
    flag.StringVar(&options.Exclude, "e", defaultExclude, "Exclude regex")
    flag.StringVar(&options.Format, "f", "", "Output format (json, csv, text)")
    flag.StringVar(&options.File, "o", "", "Output file")
    flag.Parse()

    // Check if the URL is provided
    if options.URL == "" {
        return nil, fmt.Errorf("URL is not provided")
    }

    // Check if the output format is supported
    if options.Format != "" && options.Format != "json" && options.Format != "csv" && options.Format != "text" {
        return nil, fmt.Errorf("Unsupported output format")
    }

    // Return the parsed options
    return options, nil
}


    // Define a function to extract the links from the HTML page
func ExtractLinksFromHTML(content string) ([]string, error) {
    // Initialize a slice to hold the extracted links
    var links []string

    // Define a regular expression pattern to match the links in the HTML page
    pattern := `(?i)(?:href|src)="([^"]+)"`

    // Compile the regular expression pattern
    regex, err := regexp.Compile(pattern)
    if err != nil {
        return links, err
    }

    // Execute the regular expression on the content and extract the matches
    matches := regex.FindAllStringSubmatch(content, -1)
    if matches == nil {
        return links, nil
    }

    // Iterate over the matches and add the links to the slice
    for _, match := range matches {
        links = append(links, match[1])
    }

    // Return the slice of extracted links
    return links, nil
}

// Define a function to extract the links from the JavaScript files
func ExtractLinksFromJS(content string) ([]string, error) {
    // Initialize a slice to hold the extracted links
    var links []string

    // Define a regular expression pattern to match the links in the JavaScript files
    pattern := `(?i)https?://[^"']+`

    // Compile the regular expression pattern
    regex, err := regexp.Compile(pattern)
    if err != nil {
        return links, err
    }

    // Execute the regular expression on the content and extract the matches
    matches := regex.FindAllStringSubmatch(content, -1)
    if matches == nil {
        return links, nil
    }

    // Iterate over the matches and add the links to the slice
    for _, match := range matches {
        links = append(links, match[1])
    }

    // Return the slice of extracted links
    return links, nil
}

// Define a function to crawl the URL
func CrawlURL(options *Options, counter *int, wg *sync.WaitGroup, queue chan string, visited map[string]bool, results chan []string) {
    // Increment the counter
    *counter++

    // Define a function to cleanup after crawling the URL
    defer func() {
        // Decrement the counter
        *counter--

        // Notify the WaitGroup that the goroutine is done
        wg.Done()
    }()

    // Get the URL from the queue
    u := <-queue

    // Check if the URL is already visited
    if _, ok := visited[u]; ok {
        return
    }

    // Set the URL as visited
    visited[u] = true

    // Parse the URL
    url, err := url.Parse(u)
    if err != nil {
        return
    }

    // Construct the request
    req, err := http.NewRequest("GET", url.String(), nil)
    if err != nil {
        return
    }

    // Set the User-Agent
    if options.UserAgent != "" {
        req.Header.Set("User-Agent", options.UserAgent)
    }

    // Set the custom headers
    if options.Headers != "" {
        headers := strings.Split(options.Headers, ";")
        for _, header := range headers {
            parts := strings.Split(header, ":")
            if len(parts) == 2 {
                req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
            }
        }
    }

    // Set the cookies
    if options.Cookies != "" {
        cookies := strings.Split(options.Cookies, ";")
        for _, cookie := range cookies {
            parts := strings.Split(cookie, "=")
            if len(parts) == 2 {
                req.AddCookie(&http.Cookie{
                    Name:  strings.TrimSpace(parts[0]),
                    Value: strings.TrimSpace(parts[1]),
                })
            }
        }
    }

    // Create a new HTTP client with the specified timeout
    client := http.Client{
        Timeout: options.Timeout,
    }

    // Send the request and get the response
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

   // Check the Content-Type of the response
contentType := resp.Header.Get("Content-Type")
if strings.Contains(contentType, "application/javascript") {
    // Extract the links from the JavaScript file
    links, err := ExtractLinksFromJS(string(body))
    if err != nil {
        return
    }

    // Add the links to the queue
    for _, link := range links {
        queue <- link
    }
}

// Check the output format and send the results to the appropriate channel
switch options.OutputFormat {
case "json":
	results <- []string{u, string(body), "json"}
case "csv":
	results <- []string{u, string(body), "csv"}
default:
	results <- []string{u, string(body)}
}
// Check if the crawling should stop
if options.MaxDepth > 0 && depth >= options.MaxDepth {
    return
}

// Check if the crawling should stop
if options.MaxPages > 0 && len(visited) >= options.MaxPages {
    return
}

// Check if the crawling should stop
if options.MaxWorkers > 0 && *counter >= options.MaxWorkers {
    return
}

// Check if the crawling should stop
if len(queue) == 0 {
    return
}

// Start a new goroutine to crawl the next URL in the queue
wg.Add(1)
go CrawlURL(options, counter, wg, queue, visited, results)
// Wait for all the goroutines to finish
wg.Wait()

// Close the channels
close(queue)
close(results)

// Check the output format and print the results to the screen
switch options.OutputFormat {
    case "json":
        printResultsJSON(results)
    case "csv":
        printResultsCSV(results)
    default:
        printResultsPlainText(results)
}
// Function to print the results in JSON format
func printResultsJSON(results chan []string) {
    // Open the output file, if specified
    var outputFile *os.File
    if options.OutputFile != "" {
        var err error
        outputFile, err = os.Create(options.OutputFile)
        if err != nil {
            log.Fatalf("Error creating output file: %s", err)
        }
        defer outputFile.Close()
    }

    // Iterate over the results and print them in JSON format
    for result := range results {
        // Create the JSON object
        obj := map[string]string{
            "url":     result[0],
            "content": result[1],
        }

        // Convert the JSON object to bytes
        jsonBytes, err := json.Marshal(obj)
        if err != nil {
            log.Printf("Error marshalling JSON object: %s", err)
            continue
        }

        // Print the JSON object to the screen, or to the output file
        if outputFile != nil {
            outputFile.Write(jsonBytes)
        } else {
            fmt.Println(string(jsonBytes))
        }
    }
}
// Function to print the results in CSV format
func printResultsCSV(results chan []string) {
    // Open the output file, if specified
    var outputFile *os.File
    if options.OutputFile != "" {
        var err error
        outputFile, err = os.Create(options.OutputFile)
        if err != nil {
            log.Fatalf("Error creating output file: %s", err)
        }
        defer outputFile.Close()
    }

    // Iterate over the results and print them in CSV format
    for result := range results {
        // Create the CSV record
        record := []string{result[0], result[1]}

        // Convert the CSV record to bytes
        csvBytes := csv.NewWriter(bytes.NewBufferString("")).Write(record)

        // Print the CSV record to the screen, or to the output file
        if outputFile != nil {
            outputFile.Write(csvBytes)
        } else {
            fmt.Println(string(csvBytes))
        }
    }
}
// Function to print the results in plain text format
func printResultsPlainText(results chan []string) {
    // Open the output file, if specified
    var outputFile *os.File
    if options.OutputFile != "" {
        var err error
        outputFile, err = os.Create(options.OutputFile)
        if err != nil {
            log.Fatalf("Error creating output file: %s", err)
        }
        defer outputFile.Close()
    }

    // Iterate over the results and print them in plain text format
    for result := range results {
        // Print the URL and the content to the screen, or to the output file
        if outputFile != nil {
            outputFile.WriteString(result[0] + "\n" + result[1] + "\n")
        } else {
            fmt.Println(result[0])
            fmt.Println(result[1])
        }
    }
}
// Parse the command-line arguments
options := parseArgs()

// Set the default options, if not specified
if options.MaxDepth == 0 {
    options.MaxDepth = DefaultMaxDepth
}
if options.MaxUrls == 0 {
    options.MaxUrls = DefaultMaxUrls
}
if options.Concurrency == 0 {
    options.Concurrency = DefaultConcurrency
}
if options.Timeout == 0 {
    options.Timeout = DefaultTimeout
}

// Validate the options
if options.MaxDepth < 0 {
    log.Fatalf("Invalid maximum depth: %d", options.MaxDepth)
}
if options.MaxUrls < 0 {
    log.Fatalf("Invalid maximum number of URLs: %d", options.MaxUrls)
}
if options.Concurrency < 0 {
    log.Fatalf("Invalid concurrency level: %d", options.Concurrency)
}
if options.Timeout < 0 {
    log.Fatalf("Invalid timeout: %d", options.Timeout)
}

// Initialize the channels
queue := make(chan string)
results := make(chan []string)

// Initialize the wait group
var wg sync.WaitGroup
// Start the crawl
wg.Add(1)
go crawl(queue, results, &wg, options)
queue <- options.SeedUrl

// Wait for the crawl to complete
wg.Wait()

// Close the channels
close(queue)
close(results)

// Print the results
switch options.OutputFormat {
case "json":
    printResultsJSON(results)
case "csv":
    printResultsCSV(results)
default:
    printResultsPlainText(results)
}
